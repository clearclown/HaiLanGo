package tts

import (
	"bytes"
	"context"
	"errors"
	"io"
)

// MockTTSClient はモックTTSクライアント
type MockTTSClient struct {
	dataDir   string
	mockAudio []byte
}

// NewMockTTSClient は新しいモックTTSクライアントを作成する
func NewMockTTSClient() *MockTTSClient {
	return &MockTTSClient{
		dataDir:   "./mocks/data",
		mockAudio: nil,
	}
}

// Synthesize はテキストを音声に変換する（モック）
func (m *MockTTSClient) Synthesize(ctx context.Context, text string, language string, voice string) (io.ReadCloser, error) {
	if text == "" {
		return nil, errors.New("text cannot be empty")
	}

	// カスタム音声データが設定されている場合はそれを返す
	if m.mockAudio != nil {
		return io.NopCloser(bytes.NewReader(m.mockAudio)), nil
	}

	// デフォルトのモック音声データを生成
	audioData := m.generateMockAudio(text, language, voice)
	return io.NopCloser(bytes.NewReader(audioData)), nil
}

// generateMockAudio はモック音声データを生成する
func (m *MockTTSClient) generateMockAudio(text string, language string, voice string) []byte {
	// MP3ヘッダーを模した疑似データ
	header := []byte{0x49, 0x44, 0x33} // "ID3" - MP3 ID3タグヘッダー

	// テキストの長さに基づいたダミーデータ
	dataSize := len(text) * 100
	if dataSize < 1000 {
		dataSize = 1000
	}
	data := make([]byte, dataSize)
	for i := range data {
		data[i] = byte(i % 256)
	}

	return append(header, data...)
}

// GetVoices は指定言語で利用可能な音声一覧を取得（モック）
func (m *MockTTSClient) GetVoices(ctx context.Context, language string) ([]Voice, error) {
	// 言語ごとのモック音声データ
	allVoices := []Voice{
		// 英語
		{ID: "en-US-JennyNeural", Name: "Jenny", Language: "en-US", Gender: "female", Type: "neural", SampleRate: 24000},
		{ID: "en-US-GuyNeural", Name: "Guy", Language: "en-US", Gender: "male", Type: "neural", SampleRate: 24000},
		{ID: "en-GB-SoniaNeural", Name: "Sonia", Language: "en-GB", Gender: "female", Type: "neural", SampleRate: 24000},
		// 日本語
		{ID: "ja-JP-NanamiNeural", Name: "Nanami", Language: "ja-JP", Gender: "female", Type: "neural", SampleRate: 24000},
		{ID: "ja-JP-KeitaNeural", Name: "Keita", Language: "ja-JP", Gender: "male", Type: "neural", SampleRate: 24000},
		// ロシア語
		{ID: "ru-RU-SvetlanaNeural", Name: "Svetlana", Language: "ru-RU", Gender: "female", Type: "neural", SampleRate: 24000},
		{ID: "ru-RU-DmitryNeural", Name: "Dmitry", Language: "ru-RU", Gender: "male", Type: "neural", SampleRate: 24000},
		// 中国語
		{ID: "zh-CN-XiaoxiaoNeural", Name: "Xiaoxiao", Language: "zh-CN", Gender: "female", Type: "neural", SampleRate: 24000},
		{ID: "zh-CN-YunxiNeural", Name: "Yunxi", Language: "zh-CN", Gender: "male", Type: "neural", SampleRate: 24000},
		// スペイン語
		{ID: "es-ES-ElviraNeural", Name: "Elvira", Language: "es-ES", Gender: "female", Type: "neural", SampleRate: 24000},
		// フランス語
		{ID: "fr-FR-DeniseNeural", Name: "Denise", Language: "fr-FR", Gender: "female", Type: "neural", SampleRate: 24000},
		// ドイツ語
		{ID: "de-DE-KatjaNeural", Name: "Katja", Language: "de-DE", Gender: "female", Type: "neural", SampleRate: 24000},
		// ポルトガル語
		{ID: "pt-BR-FranciscaNeural", Name: "Francisca", Language: "pt-BR", Gender: "female", Type: "neural", SampleRate: 24000},
		// イタリア語
		{ID: "it-IT-ElsaNeural", Name: "Elsa", Language: "it-IT", Gender: "female", Type: "neural", SampleRate: 24000},
		// クルド語（マイナー言語サポート）
		{ID: "ku-Arab-FahdNeural", Name: "Fahd", Language: "ku-Arab", Gender: "male", Type: "neural", SampleRate: 24000},
		// アムハラ語（マイナー言語サポート）
		{ID: "am-ET-MekdesNeural", Name: "Mekdes", Language: "am-ET", Gender: "female", Type: "neural", SampleRate: 24000},
	}

	// 言語フィルタリング
	if language == "" {
		return allVoices, nil
	}

	var filteredVoices []Voice
	for _, v := range allVoices {
		// 言語コードの前方一致でフィルタリング
		if len(v.Language) >= len(language) && v.Language[:len(language)] == language {
			filteredVoices = append(filteredVoices, v)
		}
	}

	return filteredVoices, nil
}

// GetSupportedLanguages はサポートされている言語一覧を取得（モック）
func (m *MockTTSClient) GetSupportedLanguages(ctx context.Context) ([]string, error) {
	// Azureは140言語をサポート（モックでは主要言語をシミュレート）
	return []string{
		"en", "ja", "zh", "ru", "es", "fr", "de", "pt", "it", "ko",
		"ar", "he", "fa", "tr", "nl", "pl", "sv", "da", "no", "fi",
		"th", "vi", "id", "ms", "hi", "bn", "ta", "te", "mr", "gu",
		"ku", "am", // マイナー言語
	}, nil
}

// GetName はプロバイダー名を返す
func (m *MockTTSClient) GetName() string {
	return "mock"
}

// SetMockAudio はカスタム音声データを設定する（テスト用）
func (m *MockTTSClient) SetMockAudio(audioData []byte) error {
	m.mockAudio = audioData
	return nil
}

// Ensure MockTTSClient implements TTSClient interface
var _ TTSClient = (*MockTTSClient)(nil)
