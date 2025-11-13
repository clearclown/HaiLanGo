'use client';

export interface DiffViewerProps {
  originalText: string;
  correctedText: string;
}

/**
 * DiffViewer component for displaying text differences
 */
export function DiffViewer({ originalText, correctedText }: DiffViewerProps) {
  const hasChanges = originalText !== correctedText;

  // Simple word-level diff calculation
  const calculateDiff = () => {
    if (!hasChanges) {
      return {
        original: originalText,
        corrected: correctedText,
        addedWords: 0,
        removedWords: 0,
      };
    }

    const originalWords = originalText.split(/\s+/);
    const correctedWords = correctedText.split(/\s+/);

    const added = correctedWords.filter((word) => !originalWords.includes(word));
    const removed = originalWords.filter((word) => !correctedWords.includes(word));

    return {
      original: originalText,
      corrected: correctedText,
      addedWords: added.length,
      removedWords: removed.length,
    };
  };

  const diff = calculateDiff();

  return (
    <div className="diff-viewer" data-testid="diff-viewer">
      <div className="diff-header">
        <h3 className="diff-title">Text Comparison</h3>
        {hasChanges && (
          <div className="diff-stats" data-testid="diff-stats">
            {diff.removedWords > 0 && (
              <span className="stat-removed">-{diff.removedWords} words</span>
            )}
            {diff.addedWords > 0 && (
              <span className="stat-added">+{diff.addedWords} words</span>
            )}
          </div>
        )}
        {!hasChanges && (
          <span className="no-changes" data-testid="no-changes">
            No changes
          </span>
        )}
      </div>

      <div className="diff-content">
        <div className="diff-section">
          <div className="section-label original-label">Original Text</div>
          <div
            className="section-text original-text"
            data-testid="original-text"
          >
            {originalText || <em className="empty-text">No text</em>}
          </div>
        </div>

        <div className="diff-divider" />

        <div className="diff-section">
          <div className="section-label corrected-label">Corrected Text</div>
          <div
            className="section-text corrected-text"
            data-testid="corrected-text"
          >
            {correctedText || <em className="empty-text">No text</em>}
          </div>
        </div>
      </div>

      <style jsx>{`
        .diff-viewer {
          border: 1px solid #e0e6ed;
          border-radius: 8px;
          padding: 16px;
          background: #ffffff;
        }

        .diff-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 16px;
        }

        .diff-title {
          font-size: 18px;
          font-weight: 600;
          color: #2c3e50;
          margin: 0;
        }

        .diff-stats {
          display: flex;
          gap: 12px;
          font-size: 14px;
          font-weight: 500;
        }

        .stat-removed {
          color: #e74c3c;
        }

        .stat-added {
          color: #27ae60;
        }

        .no-changes {
          color: #7f8c8d;
          font-size: 14px;
        }

        .diff-content {
          display: grid;
          grid-template-columns: 1fr auto 1fr;
          gap: 16px;
        }

        .diff-section {
          display: flex;
          flex-direction: column;
          gap: 8px;
        }

        .section-label {
          font-size: 12px;
          font-weight: 600;
          text-transform: uppercase;
          letter-spacing: 0.5px;
        }

        .original-label {
          color: #e74c3c;
        }

        .corrected-label {
          color: #27ae60;
        }

        .section-text {
          padding: 12px;
          border-radius: 4px;
          font-size: 14px;
          line-height: 1.6;
          min-height: 100px;
          white-space: pre-wrap;
          word-break: break-word;
        }

        .original-text {
          background: #fff5f5;
          border: 1px solid #ffdddd;
        }

        .corrected-text {
          background: #f0fff4;
          border: 1px solid #c6f6d5;
        }

        .empty-text {
          color: #7f8c8d;
          font-style: italic;
        }

        .diff-divider {
          width: 1px;
          background: #e0e6ed;
          align-self: stretch;
        }

        @media (max-width: 768px) {
          .diff-content {
            grid-template-columns: 1fr;
          }

          .diff-divider {
            display: none;
          }
        }
      `}</style>
    </div>
  );
}
