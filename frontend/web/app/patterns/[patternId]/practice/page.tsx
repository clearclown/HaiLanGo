'use client';

import { PatternPractice } from '@/components/patterns/PatternPractice';
import type { Pattern, PatternPractice as PatternPracticeType } from '@/lib/types/pattern';
import { getPatternPractice } from '@/services/patternsApi';
import { useEffect, useState } from 'react';

interface PageProps {
  params: {
    patternId: string;
  };
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export default function PracticePage({ params }: PageProps) {
  const [pattern, setPattern] = useState<Pattern | null>(null);
  const [practices, setPractices] = useState<PatternPracticeType[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);

        // Fetch pattern details
        const patternResponse = await fetch(`${API_BASE_URL}/api/v1/patterns/${params.patternId}`);
        if (!patternResponse.ok) {
          throw new Error('Failed to fetch pattern');
        }
        const patternData = await patternResponse.json();
        setPattern(patternData);

        // Fetch practice questions
        const practiceData = await getPatternPractice(params.patternId);
        setPractices(practiceData);
      } catch (err) {
        setError(err instanceof Error ? err : new Error('Failed to fetch data'));
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [params.patternId]);

  const handleComplete = (score: { correct: number; total: number }) => {
    // Save score to API
    console.log('Practice completed:', score);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-lg">Loading practice...</div>
      </div>
    );
  }

  if (error || !pattern) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-red-500">
          <p>Error: {error?.message || 'Pattern not found'}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <PatternPractice pattern={pattern} practices={practices} onComplete={handleComplete} />
    </div>
  );
}
