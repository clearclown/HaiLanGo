'use client';

import { useState, useEffect } from 'react';
import { ocrApiClient, type OCRTextCorrection } from '../../services/ocrApi';

export interface OCRTextEditorProps {
  bookId: string;
  pageId: string;
  originalText: string;
  correctedText?: string;
  onSave?: (correction: OCRTextCorrection) => void;
  onError?: (error: Error) => void;
}

/**
 * OCRTextEditor component for editing OCR text
 */
export function OCRTextEditor({
  bookId,
  pageId,
  originalText,
  correctedText,
  onSave,
  onError,
}: OCRTextEditorProps) {
  const [text, setText] = useState(correctedText || originalText);
  const [isSaving, setIsSaving] = useState(false);
  const [hasChanges, setHasChanges] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  // Track changes
  useEffect(() => {
    const currentBase = correctedText || originalText;
    setHasChanges(text !== currentBase);
  }, [text, originalText, correctedText]);

  // Handle text change
  const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setText(e.target.value);
    setError(null);
    setSuccess(false);
  };

  // Validate text
  const validateText = (value: string): string | null => {
    if (value.trim() === '') {
      return 'Text cannot be empty';
    }

    if (value.length > 10000) {
      return 'Text exceeds maximum length of 10,000 characters';
    }

    return null;
  };

  // Handle save
  const handleSave = async () => {
    // Validate
    const validationError = validateText(text);
    if (validationError) {
      setError(validationError);
      return;
    }

    setIsSaving(true);
    setError(null);
    setSuccess(false);

    try {
      const response = await ocrApiClient.updateOCRText(bookId, pageId, text);

      setSuccess(true);
      setHasChanges(false);

      if (onSave) {
        onSave(response.correction);
      }

      // Clear success message after 3 seconds
      setTimeout(() => setSuccess(false), 3000);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to save changes';
      setError(errorMessage);

      if (onError && err instanceof Error) {
        onError(err);
      }
    } finally {
      setIsSaving(false);
    }
  };

  // Handle reset
  const handleReset = () => {
    setText(correctedText || originalText);
    setError(null);
    setSuccess(false);
  };

  return (
    <div className="ocr-text-editor" data-testid="ocr-text-editor">
      <div className="editor-header">
        <h3 className="editor-title">Edit OCR Text</h3>
        <div className="editor-info">
          <span className="char-count">
            {text.length} / 10,000 characters
          </span>
          {hasChanges && (
            <span className="unsaved-indicator" data-testid="unsaved-indicator">
              • Unsaved changes
            </span>
          )}
        </div>
      </div>

      <div className="editor-content">
        <textarea
          className="text-editor"
          value={text}
          onChange={handleTextChange}
          placeholder="Enter corrected OCR text..."
          rows={10}
          disabled={isSaving}
          data-testid="text-editor-textarea"
          aria-label="OCR text editor"
        />
      </div>

      {error && (
        <div className="error-message" role="alert" data-testid="error-message">
          {error}
        </div>
      )}

      {success && (
        <div className="success-message" role="status" data-testid="success-message">
          ✓ Changes saved successfully
        </div>
      )}

      <div className="editor-actions">
        <button
          type="button"
          className="btn-reset"
          onClick={handleReset}
          disabled={!hasChanges || isSaving}
          data-testid="reset-button"
        >
          Reset
        </button>
        <button
          type="button"
          className="btn-save"
          onClick={handleSave}
          disabled={!hasChanges || isSaving}
          data-testid="save-button"
        >
          {isSaving ? 'Saving...' : 'Save Changes'}
        </button>
      </div>

      <style jsx>{`
        .ocr-text-editor {
          border: 1px solid #e0e6ed;
          border-radius: 8px;
          padding: 16px;
          background: #ffffff;
        }

        .editor-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 12px;
        }

        .editor-title {
          font-size: 18px;
          font-weight: 600;
          color: #2c3e50;
          margin: 0;
        }

        .editor-info {
          display: flex;
          align-items: center;
          gap: 12px;
          font-size: 14px;
        }

        .char-count {
          color: #7f8c8d;
        }

        .unsaved-indicator {
          color: #f39c12;
          font-weight: 500;
        }

        .editor-content {
          margin-bottom: 12px;
        }

        .text-editor {
          width: 100%;
          padding: 12px;
          border: 1px solid #e0e6ed;
          border-radius: 4px;
          font-family: inherit;
          font-size: 14px;
          line-height: 1.5;
          resize: vertical;
          min-height: 200px;
        }

        .text-editor:focus {
          outline: none;
          border-color: #4a90e2;
          box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.1);
        }

        .text-editor:disabled {
          background: #f5f7fa;
          cursor: not-allowed;
        }

        .error-message {
          padding: 12px;
          background: #fee;
          border: 1px solid #fcc;
          border-radius: 4px;
          color: #c00;
          margin-bottom: 12px;
          font-size: 14px;
        }

        .success-message {
          padding: 12px;
          background: #efe;
          border: 1px solid #cfc;
          border-radius: 4px;
          color: #060;
          margin-bottom: 12px;
          font-size: 14px;
        }

        .editor-actions {
          display: flex;
          justify-content: flex-end;
          gap: 12px;
        }

        button {
          padding: 8px 16px;
          border-radius: 4px;
          font-size: 14px;
          font-weight: 500;
          cursor: pointer;
          transition: all 0.2s;
        }

        .btn-reset {
          background: #ffffff;
          border: 1px solid #e0e6ed;
          color: #2c3e50;
        }

        .btn-reset:hover:not(:disabled) {
          background: #f5f7fa;
        }

        .btn-save {
          background: #4a90e2;
          border: 1px solid #4a90e2;
          color: #ffffff;
        }

        .btn-save:hover:not(:disabled) {
          background: #3a7bc8;
        }

        button:disabled {
          opacity: 0.5;
          cursor: not-allowed;
        }
      `}</style>
    </div>
  );
}
