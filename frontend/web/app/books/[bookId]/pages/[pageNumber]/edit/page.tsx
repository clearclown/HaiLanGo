'use client';

import { OCRTextEditor } from '@/components/ocr-editor/OCRTextEditor';
import type { OCRTextCorrection } from '@/services/ocrApi';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

interface PageProps {
  params: {
    bookId: string;
    pageNumber: string;
  };
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export default function OCREditorPage({ params }: PageProps) {
  const router = useRouter();
  const [pageData, setPageData] = useState<{
    originalText: string;
    correctedText?: string;
    imageUrl?: string;
  } | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const fetchPageData = async () => {
      try {
        setLoading(true);

        // Fetch page OCR data
        const response = await fetch(
          `${API_BASE_URL}/api/v1/books/${params.bookId}/pages/${params.pageNumber}`
        );

        if (!response.ok) {
          throw new Error('Failed to fetch page data');
        }

        const data = await response.json();
        setPageData({
          originalText: data.ocrText || data.ocr_text || '',
          correctedText: data.correctedText || data.corrected_text,
          imageUrl: data.imageUrl || data.image_url,
        });
      } catch (err) {
        setError(err instanceof Error ? err : new Error('Failed to fetch page data'));
      } finally {
        setLoading(false);
      }
    };

    fetchPageData();
  }, [params.bookId, params.pageNumber]);

  const handleSave = (correction: OCRTextCorrection) => {
    console.log('Saved correction:', correction);
    // Update local state with the new corrected text
    setPageData((prev) =>
      prev
        ? {
            ...prev,
            correctedText: correction.corrected_text,
          }
        : null
    );
  };

  const handleError = (err: Error) => {
    console.error('OCR Editor error:', err);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-lg">Loading page data...</div>
      </div>
    );
  }

  if (error || !pageData) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-red-500">
          <p>Error: {error?.message || 'Page not found'}</p>
          <button
            type="button"
            onClick={() => router.back()}
            className="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Go Back
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        <header className="mb-6">
          <button
            type="button"
            onClick={() => router.back()}
            className="text-blue-600 hover:text-blue-800 mb-4 flex items-center"
          >
            <span className="mr-1">&larr;</span> Back
          </button>
          <h1 className="text-2xl font-bold text-gray-900">Edit OCR Text</h1>
          <p className="text-gray-600 mt-1">
            Review and correct the OCR text for page {params.pageNumber}
          </p>
        </header>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Page Image Preview */}
          {pageData.imageUrl && (
            <div className="bg-white rounded-lg shadow p-4">
              <h3 className="text-lg font-semibold mb-4">Original Image</h3>
              <div className="border rounded overflow-hidden">
                <img
                  src={pageData.imageUrl}
                  alt={`Page ${params.pageNumber}`}
                  className="w-full h-auto"
                />
              </div>
            </div>
          )}

          {/* OCR Text Editor */}
          <div className={pageData.imageUrl ? '' : 'lg:col-span-2'}>
            <OCRTextEditor
              bookId={params.bookId}
              pageId={params.pageNumber}
              originalText={pageData.originalText}
              correctedText={pageData.correctedText}
              onSave={handleSave}
              onError={handleError}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
