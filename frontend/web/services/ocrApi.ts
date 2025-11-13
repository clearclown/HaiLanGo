/**
 * OCR API Client for managing OCR text corrections
 */

export interface OCRTextCorrection {
  id: string;
  book_id: string;
  page_id: string;
  original_text: string;
  corrected_text: string;
  user_id: string;
  created_at: string;
  updated_at: string;
}

export interface OCRCorrectionHistory {
  page_id: string;
  corrections: OCRTextCorrection[];
  total_count: number;
}

export interface UpdateOCRTextRequest {
  corrected_text: string;
}

export interface UpdateOCRTextResponse {
  success: boolean;
  correction: OCRTextCorrection;
  message?: string;
}

export interface APIError {
  error: string;
}

/**
 * OCR API Client
 */
export class OCRApiClient {
  private baseUrl: string;
  private authToken: string | null = null;

  constructor(baseUrl: string = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080') {
    this.baseUrl = baseUrl;
  }

  /**
   * Set the authentication token
   */
  setAuthToken(token: string): void {
    this.authToken = token;
  }

  /**
   * Get headers for API requests
   */
  private getHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (this.authToken) {
      headers['Authorization'] = `Bearer ${this.authToken}`;
    }

    return headers;
  }

  /**
   * Update OCR text with manual corrections
   */
  async updateOCRText(
    bookId: string,
    pageId: string,
    correctedText: string,
  ): Promise<UpdateOCRTextResponse> {
    const response = await fetch(
      `${this.baseUrl}/api/v1/books/${bookId}/pages/${pageId}/ocr-text`,
      {
        method: 'PUT',
        headers: this.getHeaders(),
        body: JSON.stringify({ corrected_text: correctedText } as UpdateOCRTextRequest),
      },
    );

    if (!response.ok) {
      const error: APIError = await response.json();
      throw new Error(error.error || 'Failed to update OCR text');
    }

    return response.json();
  }

  /**
   * Get correction history for a page
   */
  async getCorrectionHistory(
    bookId: string,
    pageId: string,
    limit = 10,
    offset = 0,
  ): Promise<OCRCorrectionHistory> {
    const params = new URLSearchParams({
      limit: limit.toString(),
      offset: offset.toString(),
    });

    const response = await fetch(
      `${this.baseUrl}/api/v1/books/${bookId}/pages/${pageId}/ocr-history?${params}`,
      {
        method: 'GET',
        headers: this.getHeaders(),
      },
    );

    if (!response.ok) {
      const error: APIError = await response.json();
      throw new Error(error.error || 'Failed to get correction history');
    }

    return response.json();
  }
}

/**
 * Default OCR API client instance
 */
export const ocrApiClient = new OCRApiClient();
