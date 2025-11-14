export type BookStatus = 'uploading' | 'processing' | 'ready' | 'failed';
export type OCRStatus = 'pending' | 'processing' | 'completed' | 'failed';

export interface Book {
  id: string;
  user_id: string;
  title: string;
  target_language: string; // 学習先言語
  native_language: string; // 母国語
  reference_language?: string; // 参照言語（本に使用されている言語）
  cover_image_url?: string;
  total_pages: number;
  processed_pages: number;
  status: BookStatus;
  ocr_status: OCRStatus;
  created_at: string;
  updated_at: string;
}

export interface BookMetadata {
  title: string;
  target_language: string;
  native_language: string;
  reference_language?: string;
}

export interface BookFile {
  id: string;
  book_id: string;
  file_name: string;
  file_type: 'pdf' | 'png' | 'jpg' | 'heic';
  file_size: number;
  storage_path: string;
  page_number?: number;
  uploaded_at: string;
}

export interface UploadProgress {
  book_id: string;
  total_files: number;
  uploaded_files: number;
  total_bytes: number;
  uploaded_bytes: number;
  status: string;
  message?: string;
}

export interface Page {
  id: string;
  book_id: string;
  page_number: number;
  image_url: string;
  ocr_text: string;
  ocr_confidence: number;
  detected_lang: string;
  ocr_status: OCRStatus;
  ocr_error?: string;
  created_at: string;
  updated_at: string;
}
