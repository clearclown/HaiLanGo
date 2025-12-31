'use client';

import { AppLayout } from '@/components/layout';
import {
  Badge,
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Input,
  Progress,
} from '@/components/ui';
import { apiClient } from '@/lib/api/client';
import type { Word, WordStats } from '@/types/vocabulary';
import { useEffect, useState } from 'react';

export default function VocabularyPage() {
  const [words, setWords] = useState<Word[]>([]);
  const [stats, setStats] = useState<WordStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [sortBy, setSortBy] = useState<'created_at' | 'mastery' | 'review_count'>('created_at');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
  const [showAddForm, setShowAddForm] = useState(false);
  const [newWord, setNewWord] = useState({
    text: '',
    meaning: '',
    pronunciation: '',
    language: 'en',
    part_of_speech: '',
    example: '',
  });

  // biome-ignore lint/correctness/useExhaustiveDependencies: intentionally reload when sort changes
  useEffect(() => {
    loadVocabulary();
  }, [sortBy, sortOrder]);

  const loadVocabulary = async () => {
    try {
      setIsLoading(true);
      setError(null);

      const [wordsData, statsData] = await Promise.all([
        apiClient.vocabulary.list({
          query: searchQuery || undefined,
          sort_by: sortBy,
          sort_order: sortOrder,
          limit: 100,
        }),
        apiClient.vocabulary.getStats(),
      ]);

      setWords(wordsData.words || []);
      setStats(statsData);
    } catch (err) {
      console.error('Failed to load vocabulary:', err);
      setError('単語帳の読み込みに失敗しました');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSearch = () => {
    loadVocabulary();
  };

  const handleAddWord = async () => {
    if (!newWord.text.trim() || !newWord.language) {
      return;
    }

    try {
      await apiClient.vocabulary.add({
        text: newWord.text,
        meaning: newWord.meaning,
        pronunciation: newWord.pronunciation,
        language: newWord.language,
        part_of_speech: newWord.part_of_speech,
        example: newWord.example,
      });

      setNewWord({
        text: '',
        meaning: '',
        pronunciation: '',
        language: 'en',
        part_of_speech: '',
        example: '',
      });
      setShowAddForm(false);
      loadVocabulary();
    } catch (err) {
      console.error('Failed to add word:', err);
      setError('単語の追加に失敗しました');
    }
  };

  const handleDeleteWord = async (wordId: string) => {
    if (!confirm('この単語を削除しますか？')) {
      return;
    }

    try {
      await apiClient.vocabulary.delete(wordId);
      loadVocabulary();
    } catch (err) {
      console.error('Failed to delete word:', err);
      setError('単語の削除に失敗しました');
    }
  };

  const handleExportCSV = async () => {
    try {
      const blob = await apiClient.vocabulary.exportCSV();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'vocabulary.csv';
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (err) {
      console.error('Failed to export CSV:', err);
      setError('CSVのエクスポートに失敗しました');
    }
  };

  const getMasteryColor = (mastery: number): string => {
    if (mastery >= 80) return 'bg-green-500';
    if (mastery >= 50) return 'bg-yellow-500';
    if (mastery >= 20) return 'bg-orange-500';
    return 'bg-red-500';
  };

  const getMasteryLabel = (mastery: number): string => {
    if (mastery >= 80) return '習得済み';
    if (mastery >= 50) return '学習中';
    if (mastery >= 20) return '復習必要';
    return '未学習';
  };

  if (isLoading) {
    return (
      <AppLayout>
        <div className="container-app py-6 lg:py-8 flex items-center justify-center min-h-[60vh]">
          <div className="text-center">
            <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-primary mx-auto mb-4" />
            <p className="text-gray-600">読み込み中...</p>
          </div>
        </div>
      </AppLayout>
    );
  }

  return (
    <AppLayout>
      <div className="container-app py-6 lg:py-8">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 mb-6">
          <div>
            <h1 className="text-2xl lg:text-3xl font-bold text-gray-900">単語帳</h1>
            <p className="text-gray-600 mt-1">学習した単語を管理</p>
          </div>
          <div className="flex gap-2">
            <Button variant="secondary" onClick={handleExportCSV}>
              CSVエクスポート
            </Button>
            <Button variant="primary" onClick={() => setShowAddForm(!showAddForm)}>
              {showAddForm ? 'キャンセル' : '単語を追加'}
            </Button>
          </div>
        </div>

        {/* Error */}
        {error && (
          <div className="bg-red-50 text-red-700 p-4 rounded-lg mb-6">
            {error}
            <Button variant="ghost" size="sm" className="ml-4" onClick={() => setError(null)}>
              閉じる
            </Button>
          </div>
        )}

        {/* Stats */}
        {stats && (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            <Card>
              <CardContent className="pt-4">
                <p className="text-sm text-gray-500">総単語数</p>
                <p className="text-2xl font-bold text-primary">{stats.total_words}</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-4">
                <p className="text-sm text-gray-500">習得済み</p>
                <p className="text-2xl font-bold text-green-600">{stats.mastered_words}</p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-4">
                <p className="text-sm text-gray-500">平均習得度</p>
                <p className="text-2xl font-bold text-secondary">
                  {stats.average_mastery.toFixed(1)}%
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardContent className="pt-4">
                <p className="text-sm text-gray-500">総復習回数</p>
                <p className="text-2xl font-bold text-gray-700">{stats.total_reviews}</p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Add Word Form */}
        {showAddForm && (
          <Card className="mb-6">
            <CardHeader>
              <CardTitle>新しい単語を追加</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Input
                  label="単語"
                  value={newWord.text}
                  onChange={(e) => setNewWord({ ...newWord, text: e.target.value })}
                  placeholder="例: hello"
                  required
                />
                <Input
                  label="意味"
                  value={newWord.meaning}
                  onChange={(e) => setNewWord({ ...newWord, meaning: e.target.value })}
                  placeholder="例: こんにちは"
                />
                <Input
                  label="発音"
                  value={newWord.pronunciation}
                  onChange={(e) => setNewWord({ ...newWord, pronunciation: e.target.value })}
                  placeholder="例: /həˈloʊ/"
                />
                <div>
                  <label
                    htmlFor="language-select"
                    className="block text-sm font-medium text-gray-700 mb-1"
                  >
                    言語
                  </label>
                  <select
                    id="language-select"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-primary"
                    value={newWord.language}
                    onChange={(e) => setNewWord({ ...newWord, language: e.target.value })}
                  >
                    <option value="en">英語</option>
                    <option value="ja">日本語</option>
                    <option value="zh">中国語</option>
                    <option value="ru">ロシア語</option>
                    <option value="es">スペイン語</option>
                    <option value="fr">フランス語</option>
                    <option value="de">ドイツ語</option>
                    <option value="pt">ポルトガル語</option>
                    <option value="it">イタリア語</option>
                    <option value="fa">ペルシャ語</option>
                    <option value="he">ヘブライ語</option>
                    <option value="tr">トルコ語</option>
                  </select>
                </div>
                <Input
                  label="品詞"
                  value={newWord.part_of_speech}
                  onChange={(e) => setNewWord({ ...newWord, part_of_speech: e.target.value })}
                  placeholder="例: 名詞"
                />
                <Input
                  label="例文"
                  value={newWord.example}
                  onChange={(e) => setNewWord({ ...newWord, example: e.target.value })}
                  placeholder="例: Hello, how are you?"
                />
              </div>
              <div className="flex justify-end gap-2 mt-4">
                <Button variant="secondary" onClick={() => setShowAddForm(false)}>
                  キャンセル
                </Button>
                <Button variant="primary" onClick={handleAddWord}>
                  追加
                </Button>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Search and Filter */}
        <Card className="mb-6">
          <CardContent className="pt-4">
            <div className="flex flex-col sm:flex-row gap-4">
              <div className="flex-1">
                <Input
                  placeholder="単語を検索..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                />
              </div>
              <select
                aria-label="並び替え項目"
                className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-primary"
                value={sortBy}
                onChange={(e) =>
                  setSortBy(e.target.value as 'created_at' | 'mastery' | 'review_count')
                }
              >
                <option value="created_at">追加日順</option>
                <option value="mastery">習得度順</option>
                <option value="review_count">復習回数順</option>
              </select>
              <select
                aria-label="並び替え順序"
                className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-primary"
                value={sortOrder}
                onChange={(e) => setSortOrder(e.target.value as 'asc' | 'desc')}
              >
                <option value="desc">降順</option>
                <option value="asc">昇順</option>
              </select>
              <Button variant="primary" onClick={handleSearch}>
                検索
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Word List */}
        {words.length === 0 ? (
          <div className="text-center py-12">
            <div className="text-6xl mb-4">📚</div>
            <h2 className="text-2xl font-bold mb-2">単語がありません</h2>
            <p className="text-gray-600 mb-4">
              単語を追加して学習を始めましょう
              <br />
              または学習ページから自動収集できます
            </p>
            <Button variant="primary" onClick={() => setShowAddForm(true)}>
              最初の単語を追加
            </Button>
          </div>
        ) : (
          <div className="grid gap-4">
            {words.map((word) => (
              <Card key={word.id} className="hover:shadow-md transition-shadow">
                <CardContent className="pt-4">
                  <div className="flex flex-col md:flex-row md:items-center gap-4">
                    {/* Word Info */}
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-2">
                        <h3 className="text-xl font-bold text-gray-900">{word.text}</h3>
                        {word.pronunciation && (
                          <span className="text-gray-500 text-sm">{word.pronunciation}</span>
                        )}
                        {word.part_of_speech && (
                          <Badge variant="secondary">{word.part_of_speech}</Badge>
                        )}
                      </div>
                      <p className="text-gray-700 mb-2">{word.meaning || '（意味未設定）'}</p>
                      {word.example && (
                        <p className="text-gray-500 text-sm italic">"{word.example}"</p>
                      )}
                      {word.tags && word.tags.length > 0 && (
                        <div className="flex gap-1 mt-2">
                          {word.tags.map((tag) => (
                            <Badge key={tag} variant="default">
                              {tag}
                            </Badge>
                          ))}
                        </div>
                      )}
                    </div>

                    {/* Mastery */}
                    <div className="flex items-center gap-4">
                      <div className="w-32">
                        <div className="flex justify-between text-sm mb-1">
                          <span className="text-gray-500">習得度</span>
                          <span className="font-medium">{word.mastery.toFixed(0)}%</span>
                        </div>
                        <Progress
                          value={word.mastery}
                          className={getMasteryColor(word.mastery)}
                          size="sm"
                        />
                        <p className="text-xs text-gray-500 mt-1 text-center">
                          {getMasteryLabel(word.mastery)}
                        </p>
                      </div>
                      <div className="text-center">
                        <p className="text-2xl font-bold text-primary">{word.review_count}</p>
                        <p className="text-xs text-gray-500">復習</p>
                      </div>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleDeleteWord(word.id)}
                        className="text-red-500 hover:text-red-700"
                      >
                        削除
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </AppLayout>
  );
}
