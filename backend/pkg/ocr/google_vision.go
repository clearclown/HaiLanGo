package ocr

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	vision "cloud.google.com/go/vision/v2/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/clearclown/HaiLanGo/backend/pkg/retry"
	"google.golang.org/api/option"
)

// GoogleVisionClient はGoogle Vision APIクライアント
type GoogleVisionClient struct {
	apiKey           string
	httpClient       *http.Client
	endpoint         string
	credentialsFile  string
	useServiceAcct   bool
}

// NewGoogleVisionClient は新しいGoogle Vision APIクライアントを作成する（APIキー認証）
func NewGoogleVisionClient(apiKey string) *GoogleVisionClient {
	return &GoogleVisionClient{
		apiKey:         apiKey,
		httpClient:     &http.Client{},
		endpoint:       "https://vision.googleapis.com/v1/images:annotate",
		useServiceAcct: false,
	}
}

// NewGoogleVisionClientWithCredentials はサービスアカウント認証でクライアントを作成する
func NewGoogleVisionClientWithCredentials(credentialsFile string) *GoogleVisionClient {
	return &GoogleVisionClient{
		credentialsFile: credentialsFile,
		useServiceAcct:  true,
	}
}

// ProcessImage は画像データをOCR処理する
func (g *GoogleVisionClient) ProcessImage(ctx context.Context, imageData []byte, languages []string) (*OCRResult, error) {
	if g.useServiceAcct {
		return g.processImageWithServiceAccount(ctx, imageData, languages)
	}
	return g.processImageWithAPIKey(ctx, imageData, languages)
}

// processImageWithServiceAccount はサービスアカウントを使用してOCR処理する
func (g *GoogleVisionClient) processImageWithServiceAccount(ctx context.Context, imageData []byte, languages []string) (*OCRResult, error) {
	var opts []option.ClientOption

	// GOOGLE_APPLICATION_CREDENTIALS環境変数が設定されている場合はそれを使用
	// そうでなければ明示的に指定されたファイルを使用
	if g.credentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(g.credentialsFile))
	}

	client, err := vision.NewImageAnnotatorClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create vision client: %w", err)
	}
	defer client.Close()

	image := &visionpb.Image{
		Content: imageData,
	}

	// 言語ヒントを設定
	var imageContext *visionpb.ImageContext
	if len(languages) > 0 {
		imageContext = &visionpb.ImageContext{
			LanguageHints: languages,
		}
	}

	// テキスト検出リクエスト（BatchAnnotateImagesを使用）
	req := &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{
			{
				Image: image,
				Features: []*visionpb.Feature{
					{
						Type: visionpb.Feature_TEXT_DETECTION,
					},
				},
				ImageContext: imageContext,
			},
		},
	}

	resp, err := client.BatchAnnotateImages(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to annotate image: %w", err)
	}

	// レスポンスチェック
	if len(resp.Responses) == 0 {
		return &OCRResult{
			Text:             "",
			DetectedLanguage: "",
			Confidence:       0.0,
			Pages:            []PageOCRResult{},
		}, nil
	}

	annotateResp := resp.Responses[0]
	if annotateResp.Error != nil {
		return nil, fmt.Errorf("vision API error: %s", annotateResp.Error.Message)
	}

	// レスポンスから結果を抽出
	if len(annotateResp.TextAnnotations) == 0 {
		return &OCRResult{
			Text:             "",
			DetectedLanguage: "",
			Confidence:       0.0,
			Pages:            []PageOCRResult{},
		}, nil
	}

	// 最初のアノテーションが全体のテキスト
	fullText := annotateResp.TextAnnotations[0].Description
	locale := annotateResp.TextAnnotations[0].Locale

	return &OCRResult{
		Text:             fullText,
		DetectedLanguage: locale,
		Confidence:       0.95, // Google Vision APIは信頼度を提供しないので固定値
		Pages: []PageOCRResult{
			{
				PageNumber: 1,
				Text:       fullText,
				Confidence: 0.95,
			},
		},
	}, nil
}

// processImageWithAPIKey はAPIキーを使用してOCR処理する（従来の方法）
func (g *GoogleVisionClient) processImageWithAPIKey(ctx context.Context, imageData []byte, languages []string) (*OCRResult, error) {
	var result *OCRResult
	var lastErr error

	// リトライロジックを使用してAPI呼び出し
	retryConfig := retry.Config{
		MaxRetries:     3,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     10 * time.Second,
		Multiplier:     2.0,
	}

	shouldRetry := func(err error) bool {
		// ネットワークエラーやレート制限エラーの場合はリトライ
		return strings.Contains(err.Error(), "temporary") ||
			strings.Contains(err.Error(), "timeout") ||
			strings.Contains(err.Error(), "rate limit")
	}

	err := retry.Do(ctx, retryConfig, func(ctx context.Context) error {
		r, err := g.callAPI(ctx, imageData, languages)
		if err != nil {
			lastErr = err
			return err
		}
		result = r
		return nil
	}, shouldRetry)

	if err != nil {
		return nil, fmt.Errorf("Google Vision API call failed after retries: %w", lastErr)
	}

	return result, nil
}

// callAPI はGoogle Vision APIを呼び出す（APIキー認証）
func (g *GoogleVisionClient) callAPI(ctx context.Context, imageData []byte, languages []string) (*OCRResult, error) {
	// 画像をBase64エンコード
	encodedImage := base64.StdEncoding.EncodeToString(imageData)

	// リクエストボディを構築
	requestBody := map[string]interface{}{
		"requests": []map[string]interface{}{
			{
				"image": map[string]string{
					"content": encodedImage,
				},
				"features": []map[string]interface{}{
					{
						"type": "TEXT_DETECTION",
					},
				},
			},
		},
	}

	// 言語ヒントを追加
	if len(languages) > 0 {
		requestBody["requests"].([]map[string]interface{})[0]["imageContext"] = map[string]interface{}{
			"languageHints": languages,
		}
	}

	// JSONエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// HTTPリクエストを作成
	url := fmt.Sprintf("%s?key=%s", g.endpoint, g.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// APIを呼び出し
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// レスポンスを読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// HTTPステータスコードをチェック
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error status %d: %s", resp.StatusCode, string(body))
	}

	// レスポンスをパース
	var apiResponse struct {
		Responses []struct {
			TextAnnotations []struct {
				Description string `json:"description"`
				Locale      string `json:"locale"`
			} `json:"textAnnotations"`
		} `json:"responses"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// エラーチェック
	if apiResponse.Error != nil {
		return nil, fmt.Errorf("API error: %s (code: %d)", apiResponse.Error.Message, apiResponse.Error.Code)
	}

	// レスポンスから結果を抽出
	if len(apiResponse.Responses) == 0 || len(apiResponse.Responses[0].TextAnnotations) == 0 {
		return &OCRResult{
			Text:             "",
			DetectedLanguage: "",
			Confidence:       0.0,
			Pages:            []PageOCRResult{},
		}, nil
	}

	// 最初のアノテーションが全体のテキスト
	fullText := apiResponse.Responses[0].TextAnnotations[0].Description
	locale := apiResponse.Responses[0].TextAnnotations[0].Locale

	return &OCRResult{
		Text:             fullText,
		DetectedLanguage: locale,
		Confidence:       0.95, // Google Vision APIは信頼度を提供しないので固定値
		Pages: []PageOCRResult{
			{
				PageNumber: 1,
				Text:       fullText,
				Confidence: 0.95,
			},
		},
	}, nil
}

// GetCredentialsFile はGOOGLE_APPLICATION_CREDENTIALS環境変数またはデフォルトパスからファイルを取得
func GetCredentialsFile() string {
	// 環境変数が設定されている場合はそれを使用
	if creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); creds != "" {
		return creds
	}
	// デフォルトはプロジェクトルートのgcp_hailingo.json
	return ""
}
