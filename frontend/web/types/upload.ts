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
