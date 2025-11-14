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
