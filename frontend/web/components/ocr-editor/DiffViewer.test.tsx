import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { DiffViewer } from './DiffViewer';

describe('DiffViewer', () => {
  const originalText = 'This is the original text from OCR';
  const correctedText = 'This is the corrected text from manual edit';

  it('renders the diff viewer', () => {
    render(<DiffViewer originalText={originalText} correctedText={correctedText} />);

    expect(screen.getByTestId('diff-viewer')).toBeInTheDocument();
  });

  it('displays original and corrected text', () => {
    render(<DiffViewer originalText={originalText} correctedText={correctedText} />);

    expect(screen.getByTestId('original-text')).toHaveTextContent(originalText);
    expect(screen.getByTestId('corrected-text')).toHaveTextContent(correctedText);
  });

  it('shows diff stats when texts are different', () => {
    render(<DiffViewer originalText={originalText} correctedText={correctedText} />);

    const stats = screen.getByTestId('diff-stats');
    expect(stats).toBeInTheDocument();

    // Should show added and removed words
    expect(stats).toHaveTextContent(/[+-]\d+ words/);
  });

  it('shows "No changes" when texts are identical', () => {
    render(<DiffViewer originalText={originalText} correctedText={originalText} />);

    expect(screen.getByTestId('no-changes')).toBeInTheDocument();
    expect(screen.queryByTestId('diff-stats')).not.toBeInTheDocument();
  });

  it('handles empty original text', () => {
    render(<DiffViewer originalText="" correctedText={correctedText} />);

    const originalTextElement = screen.getByTestId('original-text');
    expect(originalTextElement.querySelector('.empty-text')).toHaveTextContent('No text');
  });

  it('handles empty corrected text', () => {
    render(<DiffViewer originalText={originalText} correctedText="" />);

    const correctedTextElement = screen.getByTestId('corrected-text');
    expect(correctedTextElement.querySelector('.empty-text')).toHaveTextContent('No text');
  });

  it('calculates word additions correctly', () => {
    const original = 'Hello world';
    const corrected = 'Hello world from testing';

    render(<DiffViewer originalText={original} correctedText={corrected} />);

    const stats = screen.getByTestId('diff-stats');
    expect(stats).toHaveTextContent('+2 words');
  });

  it('calculates word removals correctly', () => {
    const original = 'Hello world from testing';
    const corrected = 'Hello world';

    render(<DiffViewer originalText={original} correctedText={corrected} />);

    const stats = screen.getByTestId('diff-stats');
    expect(stats).toHaveTextContent('-2 words');
  });

  it('handles multiline text', () => {
    const multilineOriginal = 'Line 1\nLine 2\nLine 3';
    const multilineCorrected = 'Line 1\nModified Line 2\nLine 3';

    render(<DiffViewer originalText={multilineOriginal} correctedText={multilineCorrected} />);

    expect(screen.getByTestId('original-text')).toHaveTextContent('Line 1');
    expect(screen.getByTestId('original-text')).toHaveTextContent('Line 2');
    expect(screen.getByTestId('corrected-text')).toHaveTextContent('Modified Line 2');
  });

  it('preserves whitespace in text display', () => {
    const textWithSpaces = 'Text   with   multiple   spaces';

    render(<DiffViewer originalText={textWithSpaces} correctedText={textWithSpaces} />);

    const originalTextElement = screen.getByTestId('original-text');
    expect(originalTextElement.textContent).toContain('Text   with   multiple   spaces');
  });
});
