//! Speech-to-Text and pronunciation evaluation service abstraction

use async_trait::async_trait;
use serde::{Deserialize, Serialize};
use thiserror::Error;

#[derive(Error, Debug)]
pub enum SttError {
    #[error("Service unavailable")]
    ServiceUnavailable,
    #[error("Unsupported audio format")]
    UnsupportedFormat,
    #[error("Unsupported language: {0}")]
    UnsupportedLanguage(String),
    #[error("Audio too long: {0}ms (max: {1}ms)")]
    AudioTooLong(u64, u64),
    #[error("Audio too short: minimum 100ms required")]
    AudioTooShort,
    #[error("Recognition failed: {0}")]
    RecognitionFailed(String),
    #[error("Rate limit exceeded")]
    RateLimitExceeded,
}

/// Audio format for STT input
#[derive(Debug, Clone, Copy, PartialEq, Eq, Default, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum SttAudioFormat {
    #[default]
    Wav,
    Mp3,
    Ogg,
    Webm,
}

/// STT recognition request
#[derive(Debug, Clone)]
pub struct SttRequest {
    pub audio_data: Vec<u8>,
    pub format: SttAudioFormat,
    pub language: String,
    /// Expected text for pronunciation evaluation (if any)
    pub expected_text: Option<String>,
}

impl SttRequest {
    pub fn new(audio_data: Vec<u8>, language: String) -> Self {
        Self {
            audio_data,
            format: SttAudioFormat::default(),
            language,
            expected_text: None,
        }
    }

    pub fn with_format(mut self, format: SttAudioFormat) -> Self {
        self.format = format;
        self
    }

    pub fn with_expected_text(mut self, text: String) -> Self {
        self.expected_text = Some(text);
        self
    }
}

/// Word-level timing and confidence
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WordResult {
    pub word: String,
    pub start_ms: u64,
    pub end_ms: u64,
    pub confidence: f32,
}

/// STT recognition response
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SttResponse {
    pub text: String,
    pub language: String,
    pub confidence: f32,
    pub words: Vec<WordResult>,
    pub duration_ms: u64,
}

/// Pronunciation score for a single word
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WordScore {
    pub word: String,
    /// 0-100 pronunciation score
    pub score: u8,
    /// Specific feedback for this word
    pub feedback: Option<String>,
}

/// Pronunciation evaluation result
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PronunciationResult {
    /// Overall score (0-100)
    pub overall_score: u8,
    /// Per-word scores
    pub word_scores: Vec<WordScore>,
    /// General feedback message
    pub feedback: String,
    /// Recognized text from audio
    pub recognized_text: String,
    /// Expected text for comparison
    pub expected_text: String,
}

/// STT Provider trait - implement for different providers
#[async_trait]
pub trait SttProvider: Send + Sync {
    /// Transcribe audio to text
    async fn transcribe(&self, request: SttRequest) -> Result<SttResponse, SttError>;

    /// Check if language is supported
    fn supports_language(&self, language: &str) -> bool;

    /// Get list of supported languages
    fn supported_languages(&self) -> Vec<&'static str>;
}

/// Pronunciation evaluator - compares recognized speech against expected text
pub struct PronunciationEvaluator;

impl PronunciationEvaluator {
    /// Evaluate pronunciation by comparing recognized text with expected text
    pub fn evaluate(response: &SttResponse, expected_text: &str) -> PronunciationResult {
        let expected_words: Vec<&str> = expected_text.split_whitespace().collect();
        let recognized_words: Vec<&str> = response.text.split_whitespace().collect();

        let mut word_scores = Vec::new();
        let mut total_score: u32 = 0;

        for (i, expected) in expected_words.iter().enumerate() {
            let (score, feedback) = if let Some(recognized) = recognized_words.get(i) {
                Self::score_word(expected, recognized, &response.words, i)
            } else {
                // Word was not spoken
                (0, Some("Word was not detected".to_string()))
            };

            total_score += score as u32;
            word_scores.push(WordScore {
                word: expected.to_string(),
                score,
                feedback,
            });
        }

        let overall_score = if expected_words.is_empty() {
            0
        } else {
            (total_score / expected_words.len() as u32).min(100) as u8
        };

        let feedback = Self::generate_feedback(overall_score, &word_scores);

        PronunciationResult {
            overall_score,
            word_scores,
            feedback,
            recognized_text: response.text.clone(),
            expected_text: expected_text.to_string(),
        }
    }

    fn score_word(
        expected: &str,
        recognized: &str,
        word_results: &[WordResult],
        index: usize,
    ) -> (u8, Option<String>) {
        let expected_lower = expected.to_lowercase();
        let recognized_lower = recognized.to_lowercase();

        if expected_lower == recognized_lower {
            // Exact match - use confidence from word results if available
            let confidence_score = word_results
                .get(index)
                .map(|w| (w.confidence * 100.0) as u8)
                .unwrap_or(90);
            (confidence_score.max(70), None)
        } else {
            // Partial match - calculate similarity
            let similarity = Self::levenshtein_similarity(&expected_lower, &recognized_lower);
            let score = (similarity * 100.0) as u8;
            let feedback = format!("Expected '{}', heard '{}'", expected, recognized);
            (score, Some(feedback))
        }
    }

    fn levenshtein_similarity(a: &str, b: &str) -> f32 {
        let a_chars: Vec<char> = a.chars().collect();
        let b_chars: Vec<char> = b.chars().collect();
        let a_len = a_chars.len();
        let b_len = b_chars.len();

        if a_len == 0 && b_len == 0 {
            return 1.0;
        }
        if a_len == 0 || b_len == 0 {
            return 0.0;
        }

        let mut matrix = vec![vec![0usize; b_len + 1]; a_len + 1];

        for (i, row) in matrix.iter_mut().enumerate().take(a_len + 1) {
            row[0] = i;
        }
        for (j, cell) in matrix[0].iter_mut().enumerate().take(b_len + 1) {
            *cell = j;
        }

        for i in 1..=a_len {
            for j in 1..=b_len {
                let cost = if a_chars[i - 1] == b_chars[j - 1] {
                    0
                } else {
                    1
                };
                matrix[i][j] = (matrix[i - 1][j] + 1)
                    .min(matrix[i][j - 1] + 1)
                    .min(matrix[i - 1][j - 1] + cost);
            }
        }

        let distance = matrix[a_len][b_len];
        let max_len = a_len.max(b_len);
        1.0 - (distance as f32 / max_len as f32)
    }

    fn generate_feedback(overall_score: u8, word_scores: &[WordScore]) -> String {
        let weak_words: Vec<&str> = word_scores
            .iter()
            .filter(|w| w.score < 60)
            .map(|w| w.word.as_str())
            .collect();

        match overall_score {
            90..=100 => "Excellent pronunciation! Keep it up!".to_string(),
            75..=89 => {
                if weak_words.is_empty() {
                    "Good pronunciation! Minor improvements possible.".to_string()
                } else {
                    format!(
                        "Good overall! Practice these words: {}",
                        weak_words.join(", ")
                    )
                }
            }
            50..=74 => format!(
                "Fair attempt. Focus on: {}",
                if weak_words.is_empty() {
                    "speaking more clearly".to_string()
                } else {
                    weak_words.join(", ")
                }
            ),
            _ => "Keep practicing! Try speaking slowly and clearly.".to_string(),
        }
    }
}

/// Mock STT provider for development and testing
pub struct MockSttProvider;

impl MockSttProvider {
    pub fn new() -> Self {
        Self
    }
}

impl Default for MockSttProvider {
    fn default() -> Self {
        Self::new()
    }
}

#[async_trait]
impl SttProvider for MockSttProvider {
    async fn transcribe(&self, request: SttRequest) -> Result<SttResponse, SttError> {
        if !self.supports_language(&request.language) {
            return Err(SttError::UnsupportedLanguage(request.language));
        }

        if request.audio_data.is_empty() {
            return Err(SttError::AudioTooShort);
        }

        // In mock mode, return the expected text if provided, otherwise a generic response
        let text = request
            .expected_text
            .unwrap_or_else(|| "Mock transcribed text".to_string());

        let words: Vec<WordResult> = text
            .split_whitespace()
            .enumerate()
            .map(|(i, word)| WordResult {
                word: word.to_string(),
                start_ms: (i as u64) * 300,
                end_ms: (i as u64) * 300 + 250,
                confidence: 0.92,
            })
            .collect();

        let duration_ms = words.last().map(|w| w.end_ms).unwrap_or(0);

        Ok(SttResponse {
            text,
            language: request.language,
            confidence: 0.92,
            words,
            duration_ms,
        })
    }

    fn supports_language(&self, language: &str) -> bool {
        self.supported_languages().contains(&language)
    }

    fn supported_languages(&self) -> Vec<&'static str> {
        vec![
            "en", "ja", "zh", "ko", "es", "fr", "de", "ru", "ar", "he", "fa",
        ]
    }
}

/// OpenAI Whisper STT provider
pub struct WhisperSttProvider {
    api_key: String,
    endpoint: String,
}

impl WhisperSttProvider {
    /// Create a new Whisper STT provider
    pub fn new(api_key: &str) -> Self {
        Self {
            api_key: api_key.to_string(),
            endpoint: "https://api.openai.com/v1/audio/transcriptions".to_string(),
        }
    }

    /// Create from environment variable
    pub fn from_env() -> Option<Self> {
        std::env::var("OPENAI_API_KEY")
            .ok()
            .map(|key| Self::new(&key))
    }

    /// Map language code to Whisper's ISO-639-1 language code
    fn map_language<'a>(&self, language: &'a str) -> &'a str {
        // Whisper uses ISO-639-1 codes, which matches our convention
        language
    }

    /// Map SttAudioFormat to MIME type
    fn mime_type(format: SttAudioFormat) -> &'static str {
        match format {
            SttAudioFormat::Wav => "audio/wav",
            SttAudioFormat::Mp3 => "audio/mpeg",
            SttAudioFormat::Ogg => "audio/ogg",
            SttAudioFormat::Webm => "audio/webm",
        }
    }

    /// Map SttAudioFormat to file extension
    fn file_extension(format: SttAudioFormat) -> &'static str {
        match format {
            SttAudioFormat::Wav => "wav",
            SttAudioFormat::Mp3 => "mp3",
            SttAudioFormat::Ogg => "ogg",
            SttAudioFormat::Webm => "webm",
        }
    }
}

#[async_trait]
impl SttProvider for WhisperSttProvider {
    async fn transcribe(&self, request: SttRequest) -> Result<SttResponse, SttError> {
        const MAX_AUDIO_DURATION_MS: u64 = 300_000; // 5 minutes

        if request.audio_data.is_empty() {
            return Err(SttError::AudioTooShort);
        }

        let language = self.map_language(&request.language);
        let file_name = format!("audio.{}", Self::file_extension(request.format));
        let mime = Self::mime_type(request.format);

        // Build multipart form
        let audio_part = reqwest::multipart::Part::bytes(request.audio_data)
            .file_name(file_name)
            .mime_str(mime)
            .map_err(|e: reqwest::Error| SttError::RecognitionFailed(e.to_string()))?;

        let form = reqwest::multipart::Form::new()
            .part("file", audio_part)
            .text("model", "whisper-1")
            .text("language", language.to_string())
            .text("response_format", "verbose_json")
            .text("timestamp_granularities[]", "word");

        let client = reqwest::Client::new();
        let response: reqwest::Response = client
            .post(&self.endpoint)
            .bearer_auth(&self.api_key)
            .multipart(form)
            .send()
            .await
            .map_err(|_| SttError::ServiceUnavailable)?;

        if response.status() == 429 {
            return Err(SttError::RateLimitExceeded);
        }

        if !response.status().is_success() {
            let status = response.status();
            let body: String = response
                .text()
                .await
                .unwrap_or_else(|_| "Unknown error".to_string());
            return Err(SttError::RecognitionFailed(format!(
                "API returned status {}: {}",
                status, body
            )));
        }

        let result: serde_json::Value = response
            .json::<serde_json::Value>()
            .await
            .map_err(|e: reqwest::Error| SttError::RecognitionFailed(e.to_string()))?;

        let text = result["text"].as_str().unwrap_or("").to_string();
        let duration_ms = (result["duration"].as_f64().unwrap_or(0.0) * 1000.0) as u64;

        if duration_ms > MAX_AUDIO_DURATION_MS {
            return Err(SttError::AudioTooLong(duration_ms, MAX_AUDIO_DURATION_MS));
        }

        // Parse word-level timestamps
        let words: Vec<WordResult> = result["words"]
            .as_array()
            .map(|word_array| {
                word_array
                    .iter()
                    .map(|w| WordResult {
                        word: w["word"].as_str().unwrap_or("").trim().to_string(),
                        start_ms: (w["start"].as_f64().unwrap_or(0.0) * 1000.0) as u64,
                        end_ms: (w["end"].as_f64().unwrap_or(0.0) * 1000.0) as u64,
                        confidence: 0.9, // Whisper doesn't expose per-word confidence
                    })
                    .collect()
            })
            .unwrap_or_default();

        Ok(SttResponse {
            text,
            language: language.to_string(),
            confidence: 0.9, // Whisper doesn't expose overall confidence
            words,
            duration_ms,
        })
    }

    fn supports_language(&self, language: &str) -> bool {
        self.supported_languages().contains(&language)
    }

    fn supported_languages(&self) -> Vec<&'static str> {
        // Whisper supports 50+ languages; listing primary ones
        vec![
            "en", "ja", "zh", "ko", "es", "fr", "de", "ru", "ar", "he", "fa", "it", "pt", "nl",
            "pl", "tr", "vi", "th", "id", "sv", "da", "fi", "no", "uk", "hi",
        ]
    }
}

/// Factory function to create the appropriate STT provider
pub fn create_stt_provider() -> Box<dyn SttProvider> {
    if let Some(provider) = WhisperSttProvider::from_env() {
        Box::new(provider)
    } else {
        tracing::warn!("OPENAI_API_KEY not set, using mock STT provider");
        Box::new(MockSttProvider::new())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_stt_request_creation() {
        let request = SttRequest::new(vec![0u8; 100], "en".to_string());

        assert_eq!(request.language, "en");
        assert_eq!(request.format, SttAudioFormat::Wav);
        assert!(request.expected_text.is_none());
    }

    #[test]
    fn test_stt_request_builder() {
        let request = SttRequest::new(vec![0u8; 100], "ja".to_string())
            .with_format(SttAudioFormat::Mp3)
            .with_expected_text("hello world".to_string());

        assert_eq!(request.language, "ja");
        assert_eq!(request.format, SttAudioFormat::Mp3);
        assert_eq!(request.expected_text.as_deref(), Some("hello world"));
    }

    #[tokio::test]
    async fn test_mock_stt_transcribe() {
        let provider = MockSttProvider::new();
        let request = SttRequest::new(vec![0u8; 100], "en".to_string())
            .with_expected_text("hello world".to_string());

        let response = provider.transcribe(request).await.unwrap();

        assert_eq!(response.text, "hello world");
        assert_eq!(response.language, "en");
        assert!(response.confidence > 0.0);
        assert_eq!(response.words.len(), 2);
    }

    #[tokio::test]
    async fn test_mock_stt_without_expected_text() {
        let provider = MockSttProvider::new();
        let request = SttRequest::new(vec![0u8; 100], "en".to_string());

        let response = provider.transcribe(request).await.unwrap();

        assert_eq!(response.text, "Mock transcribed text");
    }

    #[tokio::test]
    async fn test_mock_stt_unsupported_language() {
        let provider = MockSttProvider::new();
        let request = SttRequest::new(vec![0u8; 100], "xyz".to_string());

        let result = provider.transcribe(request).await;

        assert!(matches!(result, Err(SttError::UnsupportedLanguage(_))));
    }

    #[tokio::test]
    async fn test_mock_stt_empty_audio() {
        let provider = MockSttProvider::new();
        let request = SttRequest::new(vec![], "en".to_string());

        let result = provider.transcribe(request).await;

        assert!(matches!(result, Err(SttError::AudioTooShort)));
    }

    #[test]
    fn test_mock_stt_supported_languages() {
        let provider = MockSttProvider::new();

        assert!(provider.supports_language("en"));
        assert!(provider.supports_language("ja"));
        assert!(provider.supports_language("he"));
        assert!(!provider.supports_language("xyz"));
    }

    #[test]
    fn test_stt_audio_format_default() {
        let format = SttAudioFormat::default();
        assert_eq!(format, SttAudioFormat::Wav);
    }

    #[test]
    fn test_pronunciation_evaluator_perfect_match() {
        let response = SttResponse {
            text: "hello world".to_string(),
            language: "en".to_string(),
            confidence: 0.95,
            words: vec![
                WordResult {
                    word: "hello".to_string(),
                    start_ms: 0,
                    end_ms: 300,
                    confidence: 0.95,
                },
                WordResult {
                    word: "world".to_string(),
                    start_ms: 300,
                    end_ms: 600,
                    confidence: 0.93,
                },
            ],
            duration_ms: 600,
        };

        let result = PronunciationEvaluator::evaluate(&response, "hello world");

        assert!(result.overall_score >= 90);
        assert_eq!(result.word_scores.len(), 2);
        assert!(result.word_scores[0].feedback.is_none());
        assert!(result.word_scores[1].feedback.is_none());
    }

    #[test]
    fn test_pronunciation_evaluator_partial_match() {
        let response = SttResponse {
            text: "helo warld".to_string(),
            language: "en".to_string(),
            confidence: 0.70,
            words: vec![
                WordResult {
                    word: "helo".to_string(),
                    start_ms: 0,
                    end_ms: 300,
                    confidence: 0.70,
                },
                WordResult {
                    word: "warld".to_string(),
                    start_ms: 300,
                    end_ms: 600,
                    confidence: 0.65,
                },
            ],
            duration_ms: 600,
        };

        let result = PronunciationEvaluator::evaluate(&response, "hello world");

        assert!(result.overall_score < 90);
        assert_eq!(result.word_scores.len(), 2);
        assert!(result.word_scores[0].feedback.is_some());
    }

    #[test]
    fn test_pronunciation_evaluator_missing_words() {
        let response = SttResponse {
            text: "hello".to_string(),
            language: "en".to_string(),
            confidence: 0.80,
            words: vec![WordResult {
                word: "hello".to_string(),
                start_ms: 0,
                end_ms: 300,
                confidence: 0.80,
            }],
            duration_ms: 300,
        };

        let result = PronunciationEvaluator::evaluate(&response, "hello world");

        assert_eq!(result.word_scores.len(), 2);
        assert_eq!(result.word_scores[1].score, 0);
        assert!(result.word_scores[1].feedback.is_some());
    }

    #[test]
    fn test_levenshtein_similarity() {
        assert_eq!(
            PronunciationEvaluator::levenshtein_similarity("hello", "hello"),
            1.0
        );
        assert_eq!(PronunciationEvaluator::levenshtein_similarity("", ""), 1.0);
        assert_eq!(
            PronunciationEvaluator::levenshtein_similarity("hello", ""),
            0.0
        );

        let sim = PronunciationEvaluator::levenshtein_similarity("hello", "helo");
        assert!(sim > 0.7 && sim < 1.0);
    }

    #[test]
    fn test_stt_error_display() {
        let error = SttError::ServiceUnavailable;
        assert_eq!(error.to_string(), "Service unavailable");

        let error = SttError::UnsupportedLanguage("xyz".to_string());
        assert_eq!(error.to_string(), "Unsupported language: xyz");

        let error = SttError::AudioTooLong(400_000, 300_000);
        assert_eq!(
            error.to_string(),
            "Audio too long: 400000ms (max: 300000ms)"
        );

        let error = SttError::AudioTooShort;
        assert_eq!(error.to_string(), "Audio too short: minimum 100ms required");
    }

    #[test]
    fn test_generate_feedback_excellent() {
        let feedback = PronunciationEvaluator::generate_feedback(95, &[]);
        assert!(feedback.contains("Excellent"));
    }

    #[test]
    fn test_generate_feedback_good_with_weak_words() {
        let word_scores = vec![WordScore {
            word: "difficult".to_string(),
            score: 50,
            feedback: Some("Try again".to_string()),
        }];
        let feedback = PronunciationEvaluator::generate_feedback(80, &word_scores);
        assert!(feedback.contains("difficult"));
    }

    #[test]
    fn test_generate_feedback_low_score() {
        let feedback = PronunciationEvaluator::generate_feedback(30, &[]);
        assert!(feedback.contains("practicing"));
    }
}
