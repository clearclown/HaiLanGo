import { describe, it, expect, vi, beforeEach, type Mock } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { OCRTextEditor } from './OCRTextEditor';
import { ocrApiClient } from '../../services/ocrApi';

// Mock the OCR API client
vi.mock('../../services/ocrApi', () => ({
  ocrApiClient: {
    updateOCRText: vi.fn(),
    getCorrectionHistory: vi.fn(),
  },
}));

describe('OCRTextEditor', () => {
  const mockProps = {
    bookId: 'book-123',
    pageId: 'page-456',
    originalText: 'Original OCR text',
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders the editor with original text', () => {
    render(<OCRTextEditor {...mockProps} />);

    expect(screen.getByTestId('ocr-text-editor')).toBeInTheDocument();
    expect(screen.getByTestId('text-editor-textarea')).toHaveValue('Original OCR text');
  });

  it('renders with corrected text when provided', () => {
    render(<OCRTextEditor {...mockProps} correctedText="Corrected text" />);

    expect(screen.getByTestId('text-editor-textarea')).toHaveValue('Corrected text');
  });

  it('shows unsaved indicator when text changes', () => {
    render(<OCRTextEditor {...mockProps} />);

    const textarea = screen.getByTestId('text-editor-textarea');
    fireEvent.change(textarea, { target: { value: 'Modified text' } });

    expect(screen.getByTestId('unsaved-indicator')).toBeInTheDocument();
  });

  it('updates character count when text changes', () => {
    render(<OCRTextEditor {...mockProps} />);

    const textarea = screen.getByTestId('text-editor-textarea');
    fireEvent.change(textarea, { target: { value: 'Test' } });

    expect(screen.getByText('4 / 10,000 characters')).toBeInTheDocument();
  });

  it('shows error for empty text', async () => {
    render(<OCRTextEditor {...mockProps} />);

    const textarea = screen.getByTestId('text-editor-textarea');
    fireEvent.change(textarea, { target: { value: '   ' } });

    const saveButton = screen.getByTestId('save-button');
    fireEvent.click(saveButton);

    await waitFor(() => {
      expect(screen.getByTestId('error-message')).toHaveTextContent(
        'Text cannot be empty'
      );
    });
  });

  it('shows error for text exceeding max length', async () => {
    render(<OCRTextEditor {...mockProps} />);

    const longText = 'a'.repeat(10001);
    const textarea = screen.getByTestId('text-editor-textarea');
    fireEvent.change(textarea, { target: { value: longText } });

    const saveButton = screen.getByTestId('save-button');
    fireEvent.click(saveButton);

    await waitFor(() => {
      expect(screen.getByTestId('error-message')).toHaveTextContent(
        'Text exceeds maximum length of 10,000 characters'
      );
    });
  });

  it('calls API and shows success message on save', async () => {
    const mockCorrection = {
      id: 'correction-789',
      book_id: 'book-123',
      page_id: 'page-456',
      original_text: 'Original OCR text',
      corrected_text: 'Modified text',
      user_id: 'user-001',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };

    (ocrApiClient.updateOCRText as Mock).mockResolvedValue({
      success: true,
      correction: mockCorrection,
      message: 'OCR text updated successfully',
    });

    const onSave = vi.fn();
    render(<OCRTextEditor {...mockProps} onSave={onSave} />);

    const textarea = screen.getByTestId('text-editor-textarea');
    fireEvent.change(textarea, { target: { value: 'Modified text' } });

    const saveButton = screen.getByTestId('save-button');
    fireEvent.click(saveButton);

    await waitFor(() => {
      expect(ocrApiClient.updateOCRText).toHaveBeenCalledWith(
        'book-123',
        'page-456',
        'Modified text'
      );
      expect(screen.getByTestId('success-message')).toBeInTheDocument();
      expect(onSave).toHaveBeenCalledWith(mockCorrection);
    });
  });

  it('shows error message on API failure', async () => {
    (ocrApiClient.updateOCRText as Mock).mockRejectedValue(
      new Error('API Error')
    );

    const onError = vi.fn();
    render(<OCRTextEditor {...mockProps} onError={onError} />);

    const textarea = screen.getByTestId('text-editor-textarea');
    fireEvent.change(textarea, { target: { value: 'Modified text' } });

    const saveButton = screen.getByTestId('save-button');
    fireEvent.click(saveButton);

    await waitFor(() => {
      expect(screen.getByTestId('error-message')).toHaveTextContent('API Error');
      expect(onError).toHaveBeenCalledWith(expect.any(Error));
    });
  });

  it('resets text to original when reset button is clicked', () => {
    render(<OCRTextEditor {...mockProps} />);

    const textarea = screen.getByTestId('text-editor-textarea');
    fireEvent.change(textarea, { target: { value: 'Modified text' } });

    const resetButton = screen.getByTestId('reset-button');
    fireEvent.click(resetButton);

    expect(textarea).toHaveValue('Original OCR text');
    expect(screen.queryByTestId('unsaved-indicator')).not.toBeInTheDocument();
  });

  it('resets to corrected text when available', () => {
    render(<OCRTextEditor {...mockProps} correctedText="Corrected text" />);

    const textarea = screen.getByTestId('text-editor-textarea');
    fireEvent.change(textarea, { target: { value: 'Modified again' } });

    const resetButton = screen.getByTestId('reset-button');
    fireEvent.click(resetButton);

    expect(textarea).toHaveValue('Corrected text');
  });

  it('disables buttons when no changes', () => {
    render(<OCRTextEditor {...mockProps} />);

    const saveButton = screen.getByTestId('save-button');
    const resetButton = screen.getByTestId('reset-button');

    expect(saveButton).toBeDisabled();
    expect(resetButton).toBeDisabled();
  });

  it('disables buttons while saving', async () => {
    (ocrApiClient.updateOCRText as Mock).mockImplementation(
      () => new Promise((resolve) => setTimeout(resolve, 100))
    );

    render(<OCRTextEditor {...mockProps} />);

    const textarea = screen.getByTestId('text-editor-textarea');
    fireEvent.change(textarea, { target: { value: 'Modified text' } });

    const saveButton = screen.getByTestId('save-button');
    fireEvent.click(saveButton);

    expect(saveButton).toBeDisabled();
    expect(screen.getByText('Saving...')).toBeInTheDocument();
  });
});
