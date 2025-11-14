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
