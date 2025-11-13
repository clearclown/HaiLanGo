# モック構築戦略 - APIキーなしでもテスト可能な仕組み

## 概要

HaiLanGoプロジェクトでは、外部APIキーがなくても開発・テストが可能なように、自動モック構築システムを実装しています。これにより、開発者はAPIキーの取得に時間をかけることなく、すぐに開発を開始できます。

## モックシステムの仕組み

### 基本概念

1. **環境変数による制御**: `USE_MOCK_APIS=true` を設定すると、すべての外部API呼び出しが自動的にモックに置き換わります
2. **テスト時の自動モック**: テスト実行時は `TEST_USE_MOCKS=true` により自動的にモックが使用されます
3. **モックデータの管理**: モックレスポンスは `mocks/data/` ディレクトリに保存され、再利用可能です

### 対応するAPI

以下の外部APIがモック対応しています：

- **OCR API**
  - Google Cloud Vision API
  - Azure Computer Vision API

- **TTS API**
  - Google Cloud TTS
  - Amazon Polly
  - ElevenLabs

- **STT API**
  - Google Cloud STT
  - OpenAI Whisper API

- **辞書API**
  - Oxford Dictionary API
  - Free Dictionary API

- **決済API**
  - Stripe API

## 実装方法

### バックエンド（Go）

#### 1. インターフェース定義

```go
// backend/pkg/ocr/ocr.go
package ocr

// OCRClient はOCR APIのインターフェース
type OCRClient interface {
    ProcessImage(ctx context.Context, imageData []byte, languages []string) (*OCRResult, error)
}

// GoogleVisionClient は実際のGoogle Vision APIクライアント
type GoogleVisionClient struct {
    apiKey string
}

// MockOCRClient はモックOCRクライアント
type MockOCRClient struct {
    responses map[string]*OCRResult
}
```

#### 2. ファクトリーパターンによる切り替え

```go
// backend/pkg/ocr/factory.go
package ocr

import "os"

// NewOCRClient は環境変数に基づいて適切なクライアントを返す
func NewOCRClient() OCRClient {
    useMocks := os.Getenv("USE_MOCK_APIS") == "true" ||
               os.Getenv("TEST_USE_MOCKS") == "true"

    if useMocks {
        return NewMockOCRClient()
    }

    apiKey := os.Getenv("GOOGLE_CLOUD_VISION_API_KEY")
    if apiKey == "" {
        // APIキーがない場合は自動的にモックを使用
        return NewMockOCRClient()
    }

    return NewGoogleVisionClient(apiKey)
}
```

#### 3. モック実装

```go
// backend/pkg/ocr/mock.go
package ocr

import (
    "context"
    "encoding/json"
    "os"
    "path/filepath"
)

type MockOCRClient struct {
    dataDir string
}

func NewMockOCRClient() *MockOCRClient {
    dataDir := os.Getenv("MOCK_DATA_DIR")
    if dataDir == "" {
        dataDir = "./mocks/data"
    }

    return &MockOCRClient{
        dataDir: dataDir,
    }
}

func (m *MockOCRClient) ProcessImage(ctx context.Context, imageData []byte, languages []string) (*OCRResult, error) {
    // モックデータファイルから読み込み
    mockFile := filepath.Join(m.dataDir, "ocr", "sample_response.json")

    data, err := os.ReadFile(mockFile)
    if err != nil {
        // ファイルがない場合はデフォルトのモックレスポンスを返す
        return m.generateDefaultResponse(imageData, languages), nil
    }

    var result OCRResult
    if err := json.Unmarshal(data, &result); err != nil {
        return m.generateDefaultResponse(imageData, languages), nil
    }

    return &result, nil
}

func (m *MockOCRClient) generateDefaultResponse(imageData []byte, languages []string) *OCRResult {
    return &OCRResult{
        Text: "Здравствуйте! Это пример текста из OCR.",
        DetectedLanguage: "ru",
        Confidence: 0.95,
    }
}
```

### フロントエンド（TypeScript/Next.js）

#### 1. APIクライアントの抽象化

```typescript
// frontend/web/lib/api/client.ts
interface APIClient {
  ocr: OCRClient;
  tts: TTSClient;
  stt: STTClient;
}

// 実際のAPIクライアント
class RealAPIClient implements APIClient {
  // ...
}

// モックAPIクライアント
class MockAPIClient implements APIClient {
  // ...
}
```

#### 2. 環境変数による切り替え

```typescript
// frontend/web/lib/api/factory.ts
export function createAPIClient(): APIClient {
  const useMocks = process.env.NEXT_PUBLIC_USE_MOCK_APIS === 'true' ||
                   process.env.NODE_ENV === 'test';

  if (useMocks) {
    return new MockAPIClient();
  }

  return new RealAPIClient();
}
```

#### 3. Vitestでのモック使用

```typescript
// frontend/web/lib/api/__mocks__/client.ts
import { vi } from 'vitest';

export const mockOCRClient = {
  processImage: vi.fn().mockResolvedValue({
    text: 'Здравствуйте!',
    detectedLanguage: 'ru',
    confidence: 0.95,
  }),
};

export const mockAPIClient = {
  ocr: mockOCRClient,
  // ...
};
```

## モックデータの管理

### ディレクトリ構造

```
mocks/
├── data/
│   ├── ocr/
│   │   ├── sample_response.json
│   │   └── multi_language_response.json
│   ├── tts/
│   │   ├── russian_hello.mp3
│   │   └── english_hello.mp3
│   ├── stt/
│   │   └── pronunciation_evaluation.json
│   └── dictionary/
│       └── word_definition.json
└── scripts/
    └── generate_mock_data.go  # モックデータ生成スクリプト
```

### モックデータの生成

```go
// mocks/scripts/generate_mock_data.go
package main

import (
    "encoding/json"
    "os"
    "path/filepath"
)

func generateOCRMockData() {
    response := map[string]interface{}{
        "text": "Здравствуйте! Это пример текста.",
        "detectedLanguage": "ru",
        "confidence": 0.95,
        "pages": []map[string]interface{}{
            {
                "pageNumber": 1,
                "text": "Здравствуйте!",
            },
        },
    }

    data, _ := json.MarshalIndent(response, "", "  ")
    os.WriteFile("mocks/data/ocr/sample_response.json", data, 0644)
}
```

## テストでの使用

### Goテスト

```go
// backend/internal/service/ocr/service_test.go
package ocr

import (
    "os"
    "testing"
)

func TestMain(m *testing.M) {
    // テスト実行時は自動的にモックを使用
    os.Setenv("TEST_USE_MOCKS", "true")
    code := m.Run()
    os.Exit(code)
}

func TestProcessPage(t *testing.T) {
    // モックが自動的に使用される
    service := NewOCRService()
    result, err := service.ProcessPage(context.Background(), imageData, "ru", "ja")
    // テスト実行...
}
```

### Vitestテスト

```typescript
// frontend/web/components/learning/PageLearning.test.tsx
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { PageLearning } from './PageLearning';

// モックを自動的に使用
vi.mock('@/lib/api/client', () => ({
  createAPIClient: () => ({
    ocr: {
      processImage: vi.fn().mockResolvedValue({
        text: 'Здравствуйте!',
      }),
    },
  }),
}));

describe('PageLearning', () => {
  it('should render page content', () => {
    render(<PageLearning bookId="test-book" pageNumber={1} />);
    expect(screen.getByText('Здравствуйте!')).toBeInTheDocument();
  });
});
```

### Playwright E2Eテスト

```typescript
// frontend/web/e2e/page-learning.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Page Learning', () => {
  test.beforeEach(async ({ page }) => {
    // モックAPIサーバーを起動
    await page.goto('http://localhost:3000?useMocks=true');
  });

  test('should display OCR result', async ({ page }) => {
    await page.goto('/books/test-book/pages/1');
    await expect(page.locator('text=Здравствуйте!')).toBeVisible();
  });
});
```

## CI/CDでの使用

### GitHub Actions

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Run tests with mocks
        env:
          TEST_USE_MOCKS: true
          USE_MOCK_APIS: true
        run: |
          go test ./...
          pnpm test
```

## ベストプラクティス

### 1. モックデータのバージョン管理

モックデータはGitにコミットし、チーム全体で共有します。

```bash
# .gitignore に追加しない
mocks/data/
```

### 2. モックデータの更新

実際のAPIレスポンスを取得してモックデータを更新するスクリプトを用意します。

```bash
# モックデータを更新
go run mocks/scripts/update_mock_data.go --api=ocr --save-to=mocks/data/ocr/
```

### 3. モックと実APIの切り替えテスト

定期的に実APIを使用したテストを実行し、モックが実APIと一致していることを確認します。

```bash
# 実APIを使用したテスト（APIキー必要）
USE_MOCK_APIS=false go test ./...
```

## トラブルシューティング

### モックが使用されない

1. 環境変数を確認: `USE_MOCK_APIS=true` または `TEST_USE_MOCKS=true`
2. モックデータファイルが存在するか確認: `mocks/data/`
3. ログでモック使用を確認: `LOG_LEVEL=debug`

### モックデータが古い

1. モックデータを更新: `go run mocks/scripts/update_mock_data.go`
2. Gitから最新のモックデータを取得: `git pull`

## まとめ

このモックシステムにより、開発者は：

- ✅ APIキーなしで開発を開始できる
- ✅ テストを高速に実行できる（外部API呼び出しなし）
- ✅ オフラインでも開発可能
- ✅ CI/CDでコストを削減できる

APIキーを取得した後は、`USE_MOCK_APIS=false` に設定することで、実際のAPIを使用できます。
