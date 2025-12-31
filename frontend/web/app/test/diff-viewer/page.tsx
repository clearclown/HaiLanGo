'use client';

import { DiffViewer } from '@/components/ocr-editor/DiffViewer';
import { useState } from 'react';

/**
 * Diff Viewer Test Page
 * This page is used for E2E testing of the DiffViewer component
 */
export default function DiffViewerTestPage() {
  const [originalText, setOriginalText] = useState('The quick brown fox jumps over the lazy dog.');
  const [correctedText, setCorrectedText] = useState(
    'The quick brown cat jumps over the lazy dog.'
  );

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4">
        <header className="mb-8">
          <h1 className="text-2xl font-bold text-gray-900">Diff Viewer Test Page</h1>
          <p className="text-gray-600 mt-2">This page is for testing the DiffViewer component</p>
        </header>

        {/* Controls for testing different scenarios */}
        <section className="mb-8 bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold mb-4">Test Controls</h2>

          <div className="space-y-4">
            <div>
              <label
                htmlFor="original-input"
                className="block text-sm font-medium text-gray-700 mb-1"
              >
                Original Text
              </label>
              <textarea
                id="original-input"
                value={originalText}
                onChange={(e) => setOriginalText(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                rows={3}
              />
            </div>

            <div>
              <label
                htmlFor="corrected-input"
                className="block text-sm font-medium text-gray-700 mb-1"
              >
                Corrected Text
              </label>
              <textarea
                id="corrected-input"
                value={correctedText}
                onChange={(e) => setCorrectedText(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                rows={3}
              />
            </div>

            <div className="flex gap-3">
              <button
                type="button"
                onClick={() => {
                  setOriginalText('The quick brown fox jumps over the lazy dog.');
                  setCorrectedText('The quick brown cat jumps over the lazy dog.');
                }}
                className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
              >
                Reset to Difference Example
              </button>

              <button
                type="button"
                onClick={() => {
                  setOriginalText('Same text');
                  setCorrectedText('Same text');
                }}
                className="px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600"
              >
                Set Identical Texts
              </button>

              <button
                type="button"
                onClick={() => {
                  setOriginalText('');
                  setCorrectedText('');
                }}
                className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600"
              >
                Clear All
              </button>
            </div>
          </div>
        </section>

        {/* DiffViewer Component */}
        <section className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold mb-4">DiffViewer Component</h2>
          <DiffViewer originalText={originalText} correctedText={correctedText} />
        </section>
      </div>
    </div>
  );
}
