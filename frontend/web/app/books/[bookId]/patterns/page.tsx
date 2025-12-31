'use client';

import { PatternList } from '@/components/patterns/PatternList';
import type { Pattern, PatternType } from '@/lib/types/pattern';
import { getPatterns } from '@/services/patternsApi';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

interface PageProps {
  params: {
    bookId: string;
  };
}

export default function PatternsPage({ params }: PageProps) {
  const router = useRouter();
  const [patterns, setPatterns] = useState<Pattern[]>([]);
  const [filteredPatterns, setFilteredPatterns] = useState<Pattern[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const fetchPatterns = async () => {
      try {
        setLoading(true);
        const data = await getPatterns(params.bookId);
        setPatterns(data);
        setFilteredPatterns(data);
      } catch (err) {
        setError(err instanceof Error ? err : new Error('Failed to fetch patterns'));
      } finally {
        setLoading(false);
      }
    };

    fetchPatterns();
  }, [params.bookId]);

  const handlePatternClick = (pattern: Pattern) => {
    router.push(`/patterns/${pattern.id}/practice`);
  };

  const handleFilterChange = (type: PatternType | 'all') => {
    if (type === 'all') {
      setFilteredPatterns(patterns);
    } else {
      setFilteredPatterns(patterns.filter((p) => p.type === type));
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-lg">Loading patterns...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-red-500">
          <p>Error: {error.message}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 py-6">
          <h1 className="text-2xl font-bold text-gray-900">Conversation Patterns</h1>
          <p className="text-gray-600 mt-1">Learn common patterns from this book</p>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-8">
        <PatternList
          patterns={filteredPatterns}
          onPatternClick={handlePatternClick}
          onFilterChange={handleFilterChange}
          sortBy="frequency"
        />
      </main>
    </div>
  );
}
