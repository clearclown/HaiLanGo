'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { AppLayout } from '@/components/layout';
import { Button, Card, CardContent, Input } from '@/components/ui';
import { apiClient } from '@/lib/api/client';
import type { UploadFile } from '@/types/upload';
import { FileDropzone } from '@/components/upload/FileDropzone';
import { UploadProgress } from '@/components/upload/UploadProgress';

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
      id: crypto.randomUUID(),
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
    <AppLayout title="本を追加" showBack backHref="/books">
      <div className="container-app py-6 lg:py-8">
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
                      <div className={`flex-1 h-1 transition-colors ${isActive ? 'bg-primary' : 'bg-gray-200'}`} />
                    )}
                    <div
                      className={`flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium transition-colors ${
                        isActive ? 'bg-primary text-white' : 'bg-gray-200 text-gray-500'
                      }`}
                    >
                      {stepNumber}
                    </div>
                  </div>
                  <span className={`ml-2 text-xs sm:text-sm hidden sm:block ${isActive ? 'text-primary font-medium' : 'text-gray-500'}`}>
                    {label}
                  </span>
                </div>
              );
            })}
          </div>
        </div>

        {/* Step: Metadata */}
        {step === 'metadata' && (
          <Card>
            <CardContent className="pt-6">
              <h1 className="text-xl lg:text-2xl font-bold mb-6">本の情報を入力</h1>

              <form onSubmit={handleMetadataSubmit} className="space-y-6">
                <Input
                  id="title"
                  label="本のタイトル *"
                  type="text"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="例: ロシア語入門"
                  required
                />

                <div>
                  <label htmlFor="targetLanguage" className="label">
                    学習先言語 *
                  </label>
                  <select
                    id="targetLanguage"
                    value={targetLanguage}
                    onChange={(e) => setTargetLanguage(e.target.value)}
                    className="input"
                    required
                  >
                    <option value="">選択してください</option>
                    {languages.map(lang => (
                      <option key={lang.code} value={lang.code}>{lang.name}</option>
                    ))}
                  </select>
                </div>

                <div>
                  <label htmlFor="nativeLanguage" className="label">
                    母国語 *
                  </label>
                  <select
                    id="nativeLanguage"
                    value={nativeLanguage}
                    onChange={(e) => setNativeLanguage(e.target.value)}
                    className="input"
                    required
                  >
                    {languages.map(lang => (
                      <option key={lang.code} value={lang.code}>{lang.name}</option>
                    ))}
                  </select>
                </div>

                <div>
                  <label htmlFor="referenceLanguage" className="label">
                    参照言語（本に使用されている言語）
                  </label>
                  <select
                    id="referenceLanguage"
                    value={referenceLanguage}
                    onChange={(e) => setReferenceLanguage(e.target.value)}
                    className="input"
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

                <div className="flex gap-3 pt-4">
                  <Button
                    type="button"
                    variant="ghost"
                    onClick={() => router.push('/books')}
                  >
                    キャンセル
                  </Button>
                  <Button type="submit" variant="primary" className="flex-1">
                    次へ
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        )}

        {/* Step: File Selection */}
        {step === 'files' && (
          <Card>
            <CardContent className="pt-6">
              <h1 className="text-xl lg:text-2xl font-bold mb-6">ファイルを選択</h1>
              <FileDropzone onFilesSelected={handleFilesSelected} />
            </CardContent>
          </Card>
        )}

        {/* Step: Uploading */}
        {step === 'uploading' && (
          <Card>
            <CardContent className="pt-6">
              <h1 className="text-xl lg:text-2xl font-bold mb-6">アップロード中</h1>
              <UploadProgress files={uploadFiles} onRemove={isUploading ? undefined : handleRemoveFile} />
            </CardContent>
          </Card>
        )}

        {/* Step: Completed */}
        {step === 'completed' && (
          <Card>
            <CardContent className="py-12 text-center">
              <div className="text-6xl mb-4">🎉</div>
              <h1 className="text-xl lg:text-2xl font-bold mb-4">アップロード完了！</h1>
              <p className="text-gray-600 mb-8">
                OCR処理が開始されました。処理が完了したら学習を開始できます。
              </p>
              <Button variant="primary" onClick={handleGoToBooks}>
                マイ本へ
              </Button>
            </CardContent>
          </Card>
        )}
      </div>
    </AppLayout>
  );
}
