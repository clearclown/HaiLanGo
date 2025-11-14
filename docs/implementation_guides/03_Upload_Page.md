# 実装指示書: 書籍アップロードページ

## 概要
PDF/画像ファイルをアップロードし、OCR処理を開始するための機能を実装する。ドラッグ&ドロップ、進捗表示、エラーハンドリングを含む。

## 担当範囲
- **フロントエンド**: `frontend/web/app/upload/page.tsx`
- **コンポーネント**: `frontend/web/components/upload/*`
- **バックエンドAPI**: すでに実装済み（`/api/v1/upload/*`）

## 前提条件
- Node.js 18+、pnpm がインストール済み
- バックエンドAPI が http://localhost:8080 で起動中

## 実装ステップ

### Step 1: 型定義の追加

**ファイル**: `frontend/web/types/upload.ts`

```typescript
export interface UploadFile {
  file: File;
  id: string;
  status: 'pending' | 'uploading' | 'completed' | 'failed';
  progress: number;
  error?: string;
}

export interface UploadMetadata {
  book_id: string;
  title: string;
  target_language: string;
  native_language: string;
  reference_language?: string;
}
```

### Step 2: API クライアントの拡張

**ファイル**: `frontend/web/lib/api/client.ts`

**追加する内容**:

```typescript
import type { UploadMetadata } from '@/types/upload';

upload = {
  createBook: async (metadata: Omit<UploadMetadata, 'book_id'>): Promise<{ book_id: string }> => {
    return this.fetch<{ book_id: string }>('/api/v1/upload/create', {
      method: 'POST',
      body: JSON.stringify(metadata),
    });
  },

  uploadFile: async (
    bookId: string,
    file: File,
    onProgress?: (progress: number) => void
  ): Promise<{ success: boolean; file_id: string }> => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('book_id', bookId);

    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest();

      xhr.upload.addEventListener('progress', (e) => {
        if (e.lengthComputable && onProgress) {
          const progress = (e.loaded / e.total) * 100;
          onProgress(progress);
        }
      });

      xhr.addEventListener('load', () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          resolve(JSON.parse(xhr.responseText));
        } else {
          reject(new Error(`Upload failed: ${xhr.statusText}`));
        }
      });

      xhr.addEventListener('error', () => reject(new Error('Upload failed')));

      xhr.open('POST', `${API_BASE_URL}/api/v1/upload/file`);
      xhr.send(formData);
    });
  },

  complete: async (bookId: string): Promise<{ success: boolean }> => {
    return this.fetch<{ success: boolean }>('/api/v1/upload/complete', {
      method: 'POST',
      body: JSON.stringify({ book_id: bookId }),
    });
  },
};
```

### Step 3: FileDropzone コンポーネントの作成

**ファイル**: `frontend/web/components/upload/FileDropzone.tsx`

```typescript
'use client';

import { useState, useRef, DragEvent } from 'react';

interface FileDropzoneProps {
  onFilesSelected: (files: File[]) => void;
  accept?: string;
  maxFiles?: number;
}

export function FileDropzone({
  onFilesSelected,
  accept = '.pdf,.png,.jpg,.jpeg,.heic',
  maxFiles = 100,
}: FileDropzoneProps) {
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);

    const files = Array.from(e.dataTransfer.files);
    if (files.length > maxFiles) {
      alert(`最大${maxFiles}ファイルまでアップロードできます`);
      return;
    }

    onFilesSelected(files);
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files) {
      onFilesSelected(Array.from(files));
    }
  };

  const handleClick = () => {
    fileInputRef.current?.click();
  };

  return (
    <div
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      onClick={handleClick}
      className={`border-2 border-dashed rounded-lg p-12 text-center cursor-pointer transition-colors ${
        isDragging
          ? 'border-blue-500 bg-blue-50'
          : 'border-gray-300 hover:border-gray-400'
      }`}
    >
      <input
        ref={fileInputRef}
        type="file"
        multiple
        accept={accept}
        onChange={handleFileSelect}
        className="hidden"
      />

      <div className="flex flex-col items-center gap-4">
        <div className="text-6xl">📁</div>
        <div>
          <p className="text-lg font-medium mb-2">
            ファイルを選択またはドラッグ&ドロップ
          </p>
          <p className="text-sm text-gray-600">
            PDF / PNG / JPG / HEIC
          </p>
          <p className="text-xs text-gray-500 mt-1">
            最大{maxFiles}ファイルまで
          </p>
        </div>
      </div>
    </div>
  );
}
```

### Step 4: UploadProgress コンポーネントの作成

**ファイル**: `frontend/web/components/upload/UploadProgress.tsx`

```typescript
'use client';

import type { UploadFile } from '@/types/upload';

interface UploadProgressProps {
  files: UploadFile[];
  onRemove?: (fileId: string) => void;
}

export function UploadProgress({ files, onRemove }: UploadProgressProps) {
  const getStatusIcon = (status: UploadFile['status']) => {
    switch (status) {
      case 'pending':
        return '⏳';
      case 'uploading':
        return '⬆️';
      case 'completed':
        return '✅';
      case 'failed':
        return '❌';
    }
  };

  const getStatusText = (status: UploadFile['status']) => {
    switch (status) {
      case 'pending':
        return '待機中';
      case 'uploading':
        return 'アップロード中';
      case 'completed':
        return '完了';
      case 'failed':
        return '失敗';
    }
  };

  const totalProgress = files.length > 0
    ? files.reduce((sum, file) => sum + file.progress, 0) / files.length
    : 0;

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="font-semibold">アップロード進捗</h3>
        <span className="text-sm text-gray-600">
          {files.filter(f => f.status === 'completed').length} / {files.length} ファイル完了
        </span>
      </div>

      {/* Total Progress Bar */}
      <div>
        <div className="flex justify-between text-sm text-gray-600 mb-1">
          <span>全体の進捗</span>
          <span>{Math.round(totalProgress)}%</span>
        </div>
        <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
          <div
            className="h-full bg-blue-500 transition-all duration-300"
            style={{ width: `${totalProgress}%` }}
          />
        </div>
      </div>

      {/* Individual File Progress */}
      <div className="space-y-2 max-h-96 overflow-y-auto">
        {files.map(file => (
          <div key={file.id} className="bg-white border rounded-lg p-4">
            <div className="flex items-center gap-3">
              <span className="text-2xl">{getStatusIcon(file.status)}</span>
              <div className="flex-1 min-w-0">
                <p className="font-medium truncate">{file.file.name}</p>
                <p className="text-sm text-gray-600">
                  {getStatusText(file.status)}
                  {file.error && ` - ${file.error}`}
                </p>
              </div>
              {file.status !== 'completed' && onRemove && (
                <button
                  type="button"
                  onClick={() => onRemove(file.id)}
                  className="text-gray-400 hover:text-gray-600"
                >
                  ✕
                </button>
              )}
            </div>

            {file.status === 'uploading' && (
              <div className="mt-2">
                <div className="h-1 bg-gray-200 rounded-full overflow-hidden">
                  <div
                    className="h-full bg-blue-500 transition-all duration-300"
                    style={{ width: `${file.progress}%` }}
                  />
                </div>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
```

### Step 5: Upload ページの実装

**ファイル**: `frontend/web/app/upload/page.tsx`

```typescript
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api/client';
import type { UploadFile } from '@/types/upload';
import { FileDropzone } from '@/components/upload/FileDropzone';
import { UploadProgress } from '@/components/upload/UploadProgress';
import { v4 as uuidv4 } from 'uuid';

export default function UploadPage() {
  const router = useRouter();
  const [step, setStep] = useState<'metadata' | 'files' | 'uploading' | 'completed'>('metadata');

  // Metadata state
  const [title, setTitle] = useState('');
  const [targetLanguage, setTargetLanguage] = useState('');
  const [nativeLanguage, setNativeLanguage] = useState('ja');
  const [referenceLanguage, setReferenceLanguage] = useState('');

  // Upload state
  const [bookId, setBookId] = useState('');
  const [uploadFiles, setUploadFiles] = useState<UploadFile[]>([]);
  const [isUploading, setIsUploading] = useState(false);

  const languages = [
    { code: 'ja', name: '日本語' },
    { code: 'en', name: '英語' },
    { code: 'zh', name: '中国語' },
    { code: 'ru', name: 'ロシア語' },
    { code: 'fa', name: 'ペルシャ語' },
    { code: 'he', name: 'ヘブライ語' },
    { code: 'es', name: 'スペイン語' },
    { code: 'fr', name: 'フランス語' },
    { code: 'pt', name: 'ポルトガル語' },
    { code: 'de', name: 'ドイツ語' },
    { code: 'it', name: 'イタリア語' },
    { code: 'tr', name: 'トルコ語' },
  ];

  const handleMetadataSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!title || !targetLanguage || !nativeLanguage) {
      alert('必須項目を入力してください');
      return;
    }

    try {
      const response = await apiClient.upload.createBook({
        title,
        target_language: targetLanguage,
        native_language: nativeLanguage,
        reference_language: referenceLanguage || undefined,
      });

      setBookId(response.book_id);
      setStep('files');
    } catch (error) {
      console.error('Failed to create book:', error);
      alert('本の作成に失敗しました');
    }
  };

  const handleFilesSelected = (files: File[]) => {
    const uploadFiles: UploadFile[] = files.map(file => ({
      file,
      id: uuidv4(),
      status: 'pending',
      progress: 0,
    }));

    setUploadFiles(uploadFiles);
    setStep('uploading');
    startUpload(uploadFiles);
  };

  const startUpload = async (files: UploadFile[]) => {
    setIsUploading(true);

    for (const file of files) {
      try {
        // Update status to uploading
        setUploadFiles(prev =>
          prev.map(f => f.id === file.id ? { ...f, status: 'uploading' } : f)
        );

        // Upload file
        await apiClient.upload.uploadFile(
          bookId,
          file.file,
          (progress) => {
            setUploadFiles(prev =>
              prev.map(f => f.id === file.id ? { ...f, progress } : f)
            );
          }
        );

        // Mark as completed
        setUploadFiles(prev =>
          prev.map(f => f.id === file.id ? { ...f, status: 'completed', progress: 100 } : f)
        );
      } catch (error) {
        console.error('Upload failed:', file.file.name, error);
        setUploadFiles(prev =>
          prev.map(f =>
            f.id === file.id
              ? { ...f, status: 'failed', error: 'アップロードに失敗しました' }
              : f
          )
        );
      }
    }

    // Complete upload
    try {
      await apiClient.upload.complete(bookId);
      setIsUploading(false);
      setStep('completed');
    } catch (error) {
      console.error('Failed to complete upload:', error);
      alert('アップロードの完了処理に失敗しました');
    }
  };

  const handleRemoveFile = (fileId: string) => {
    setUploadFiles(prev => prev.filter(f => f.id !== fileId));
  };

  const handleGoToBooks = () => {
    router.push('/books');
  };

  return (
    <div className="min-h-screen bg-background-secondary">
      <div className="max-w-4xl mx-auto px-4 py-8">
        {/* Progress Indicator */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
            {['メタデータ', 'ファイル選択', 'アップロード', '完了'].map((label, index) => {
              const stepNumber = index + 1;
              const currentStepIndex = ['metadata', 'files', 'uploading', 'completed'].indexOf(step) + 1;
              const isActive = stepNumber <= currentStepIndex;

              return (
                <div key={label} className="flex-1 flex items-center">
                  <div className={`flex items-center ${index > 0 ? 'w-full' : ''}`}>
                    {index > 0 && (
                      <div className={`flex-1 h-1 ${isActive ? 'bg-blue-500' : 'bg-gray-300'}`} />
                    )}
                    <div
                      className={`flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium ${
                        isActive ? 'bg-blue-500 text-white' : 'bg-gray-300 text-gray-600'
                      }`}
                    >
                      {stepNumber}
                    </div>
                  </div>
                  <span className={`ml-2 text-sm ${isActive ? 'text-blue-600 font-medium' : 'text-gray-500'}`}>
                    {label}
                  </span>
                </div>
              );
            })}
          </div>
        </div>

        {/* Step: Metadata */}
        {step === 'metadata' && (
          <div className="bg-white rounded-lg shadow-sm p-6">
            <h1 className="text-2xl font-bold mb-6">本の情報を入力</h1>

            <form onSubmit={handleMetadataSubmit} className="space-y-6">
              <div>
                <label htmlFor="title" className="block text-sm font-medium mb-2">
                  本のタイトル *
                </label>
                <input
                  id="title"
                  type="text"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="例: ロシア語入門"
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>

              <div>
                <label htmlFor="targetLanguage" className="block text-sm font-medium mb-2">
                  学習先言語 *
                </label>
                <select
                  id="targetLanguage"
                  value={targetLanguage}
                  onChange={(e) => setTargetLanguage(e.target.value)}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                >
                  <option value="">選択してください</option>
                  {languages.map(lang => (
                    <option key={lang.code} value={lang.code}>{lang.name}</option>
                  ))}
                </select>
              </div>

              <div>
                <label htmlFor="nativeLanguage" className="block text-sm font-medium mb-2">
                  母国語 *
                </label>
                <select
                  id="nativeLanguage"
                  value={nativeLanguage}
                  onChange={(e) => setNativeLanguage(e.target.value)}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                >
                  {languages.map(lang => (
                    <option key={lang.code} value={lang.code}>{lang.name}</option>
                  ))}
                </select>
              </div>

              <div>
                <label htmlFor="referenceLanguage" className="block text-sm font-medium mb-2">
                  参照言語（本に使用されている言語）
                </label>
                <select
                  id="referenceLanguage"
                  value={referenceLanguage}
                  onChange={(e) => setReferenceLanguage(e.target.value)}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="">選択してください</option>
                  {languages.map(lang => (
                    <option key={lang.code} value={lang.code}>{lang.name}</option>
                  ))}
                </select>
                <p className="text-sm text-gray-500 mt-1">
                  学習先言語と異なる言語で書かれている場合のみ選択
                </p>
              </div>

              <div className="flex gap-4">
                <button
                  type="button"
                  onClick={() => router.push('/books')}
                  className="px-6 py-3 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300"
                >
                  キャンセル
                </button>
                <button
                  type="submit"
                  className="flex-1 px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
                >
                  次へ
                </button>
              </div>
            </form>
          </div>
        )}

        {/* Step: File Selection */}
        {step === 'files' && (
          <div className="bg-white rounded-lg shadow-sm p-6">
            <h1 className="text-2xl font-bold mb-6">ファイルを選択</h1>
            <FileDropzone onFilesSelected={handleFilesSelected} />
          </div>
        )}

        {/* Step: Uploading */}
        {step === 'uploading' && (
          <div className="bg-white rounded-lg shadow-sm p-6">
            <h1 className="text-2xl font-bold mb-6">アップロード中</h1>
            <UploadProgress files={uploadFiles} onRemove={isUploading ? undefined : handleRemoveFile} />
          </div>
        )}

        {/* Step: Completed */}
        {step === 'completed' && (
          <div className="bg-white rounded-lg shadow-sm p-6 text-center">
            <div className="text-6xl mb-4">🎉</div>
            <h1 className="text-2xl font-bold mb-4">アップロード完了！</h1>
            <p className="text-gray-600 mb-8">
              OCR処理が開始されました。処理が完了したら学習を開始できます。
            </p>
            <button
              onClick={handleGoToBooks}
              className="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
            >
              マイ本へ
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
```

## テスト方法

1. ブラウザで http://localhost:3000/upload にアクセス

2. **確認項目**:
   - [ ] メタデータフォームが表示される
   - [ ] ファイルドロップゾーンが表示される
   - [ ] ドラッグ&ドロップでファイル選択できる
   - [ ] アップロード進捗が表示される
   - [ ] 完了画面が表示される
   - [ ] マイ本ページに遷移できる

## 完了条件

- [ ] 型定義が作成されている
- [ ] API クライアントが拡張されている
- [ ] FileDropzone コンポーネントが動作する
- [ ] UploadProgress コンポーネントが動作する
- [ ] Upload ページが正しくレンダリングされる
- [ ] ファイルアップロードが完了できる
- [ ] エラーハンドリングが適切に実装されている

## 参考資料

- [書籍アップロードRD](../../docs/featureRDs/2_書籍アップロード.md)
- [UI/UX設計書](../../docs/ui_ux_design_document.md)
