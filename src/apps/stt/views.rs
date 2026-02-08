//! STT ViewSet - Pronunciation evaluation endpoints

use uuid::Uuid;

use super::dto::{
    AttemptSummaryResponse, PronunciationRequest, PronunciationResponse, PronunciationStatsResponse,
    WeakWordResponse, WordScoreResponse,
};
use super::models::{AttemptStatus, PronunciationAttempt, WordFeedback};
use crate::services::stt::{PronunciationEvaluator, SttProvider, SttRequest};

/// Pronunciation evaluation result
#[derive(Debug)]
pub enum EvaluateResult {
    Success(PronunciationResponse),
    InvalidInput(String),
    ServiceError(String),
    Unauthorized,
}

/// Attempt retrieval result
#[derive(Debug)]
pub enum GetAttemptResult {
    Success(AttemptSummaryResponse),
    NotFound,
    Unauthorized,
}

/// STT ViewSet
pub struct SttViewSet;

impl SttViewSet {
    /// Evaluate pronunciation by transcribing audio and comparing with expected text
    pub async fn evaluate(
        user_id: Uuid,
        request: PronunciationRequest,
        audio_data: Vec<u8>,
        provider: &dyn SttProvider,
    ) -> EvaluateResult {
        // Validate input
        if request.expected_text.is_empty() {
            return EvaluateResult::InvalidInput("Expected text is required".to_string());
        }

        if audio_data.is_empty() {
            return EvaluateResult::InvalidInput("Audio data is required".to_string());
        }

        if !provider.supports_language(&request.language) {
            return EvaluateResult::InvalidInput(format!(
                "Language '{}' is not supported",
                request.language
            ));
        }

        // Create the attempt record
        let mut attempt = PronunciationAttempt::new(
            user_id,
            request.expected_text.clone(),
            request.language.clone(),
        );
        if let Some(page_id) = request.page_id {
            attempt = attempt.with_page(page_id);
        }
        attempt.status = AttemptStatus::Processing;

        // Build STT request
        let stt_request = SttRequest::new(audio_data, request.language)
            .with_format(request.audio_format)
            .with_expected_text(request.expected_text.clone());

        // Transcribe
        let stt_response = match provider.transcribe(stt_request).await {
            Ok(response) => response,
            Err(e) => {
                return EvaluateResult::ServiceError(e.to_string());
            }
        };

        // Evaluate pronunciation
        let evaluation =
            PronunciationEvaluator::evaluate(&stt_response, &request.expected_text);

        // Build word scores
        let word_scores: Vec<WordScoreResponse> = evaluation
            .word_scores
            .iter()
            .map(|ws| WordScoreResponse {
                word: ws.word.clone(),
                score: ws.score,
                feedback: ws.feedback.clone(),
            })
            .collect();

        EvaluateResult::Success(PronunciationResponse {
            attempt_id: attempt.id,
            overall_score: evaluation.overall_score,
            recognized_text: evaluation.recognized_text,
            expected_text: evaluation.expected_text,
            feedback: evaluation.feedback,
            word_scores,
            duration_ms: stt_response.duration_ms,
        })
    }

    /// List pronunciation attempts for a user
    pub fn list_attempts(
        user_id: Uuid,
        attempts: &[PronunciationAttempt],
    ) -> Vec<AttemptSummaryResponse> {
        attempts
            .iter()
            .filter(|a| a.user_id == user_id)
            .map(|a| AttemptSummaryResponse {
                id: a.id,
                expected_text: a.expected_text.clone(),
                recognized_text: a.recognized_text.clone(),
                overall_score: a.overall_score,
                language: a.language.clone(),
                status: a.status,
                created_at: a.created_at,
            })
            .collect()
    }

    /// Get a single attempt
    pub fn retrieve(
        user_id: Uuid,
        attempt: Option<&PronunciationAttempt>,
    ) -> GetAttemptResult {
        match attempt {
            Some(a) if a.user_id == user_id => {
                GetAttemptResult::Success(AttemptSummaryResponse {
                    id: a.id,
                    expected_text: a.expected_text.clone(),
                    recognized_text: a.recognized_text.clone(),
                    overall_score: a.overall_score,
                    language: a.language.clone(),
                    status: a.status,
                    created_at: a.created_at,
                })
            }
            Some(_) => GetAttemptResult::Unauthorized,
            None => GetAttemptResult::NotFound,
        }
    }

    /// Get pronunciation stats for a user
    pub fn stats(
        user_id: Uuid,
        attempts: &[PronunciationAttempt],
        word_feedbacks: &[WordFeedback],
    ) -> PronunciationStatsResponse {
        let user_attempts: Vec<&PronunciationAttempt> = attempts
            .iter()
            .filter(|a| a.user_id == user_id && a.status == AttemptStatus::Completed)
            .collect();

        let total_attempts = user_attempts.len() as u64;

        let (average_score, best_score) = if user_attempts.is_empty() {
            (0.0, 0)
        } else {
            let scores: Vec<u8> = user_attempts
                .iter()
                .filter_map(|a| a.overall_score)
                .collect();
            if scores.is_empty() {
                (0.0, 0)
            } else {
                let sum: u32 = scores.iter().map(|&s| s as u32).sum();
                let avg = sum as f32 / scores.len() as f32;
                let best = *scores.iter().max().unwrap_or(&0);
                (avg, best)
            }
        };

        // Aggregate weak words from feedbacks
        let attempt_ids: Vec<Uuid> = user_attempts.iter().map(|a| a.id).collect();
        let relevant_feedbacks: Vec<&WordFeedback> = word_feedbacks
            .iter()
            .filter(|f| attempt_ids.contains(&f.attempt_id) && f.score < 60)
            .collect();

        let mut weak_word_map: std::collections::HashMap<String, (f32, u32)> =
            std::collections::HashMap::new();
        for feedback in &relevant_feedbacks {
            let entry = weak_word_map
                .entry(feedback.word.clone())
                .or_insert((0.0, 0));
            entry.0 += feedback.score as f32;
            entry.1 += 1;
        }

        let mut weak_words: Vec<WeakWordResponse> = weak_word_map
            .into_iter()
            .map(|(word, (total, count))| WeakWordResponse {
                word,
                language: String::new(), // Would be populated from attempt context
                average_score: total / count as f32,
                attempt_count: count,
            })
            .collect();
        weak_words.sort_by(|a, b| a.average_score.partial_cmp(&b.average_score).unwrap());

        PronunciationStatsResponse {
            total_attempts,
            average_score,
            best_score,
            weak_words,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::services::stt::MockSttProvider;

    #[tokio::test]
    async fn test_evaluate_success() {
        let user_id = Uuid::new_v4();
        let provider = MockSttProvider::new();
        let request = PronunciationRequest {
            expected_text: "hello world".to_string(),
            language: "en".to_string(),
            audio_format: crate::services::stt::SttAudioFormat::Wav,
            page_id: None,
        };

        let result =
            SttViewSet::evaluate(user_id, request, vec![0u8; 100], &provider).await;

        match result {
            EvaluateResult::Success(response) => {
                assert_eq!(response.expected_text, "hello world");
                assert_eq!(response.recognized_text, "hello world");
                assert!(response.overall_score >= 70);
                assert_eq!(response.word_scores.len(), 2);
            }
            other => panic!("Expected Success, got {:?}", other),
        }
    }

    #[tokio::test]
    async fn test_evaluate_empty_text() {
        let user_id = Uuid::new_v4();
        let provider = MockSttProvider::new();
        let request = PronunciationRequest {
            expected_text: "".to_string(),
            language: "en".to_string(),
            audio_format: crate::services::stt::SttAudioFormat::Wav,
            page_id: None,
        };

        let result =
            SttViewSet::evaluate(user_id, request, vec![0u8; 100], &provider).await;

        assert!(matches!(result, EvaluateResult::InvalidInput(_)));
    }

    #[tokio::test]
    async fn test_evaluate_empty_audio() {
        let user_id = Uuid::new_v4();
        let provider = MockSttProvider::new();
        let request = PronunciationRequest {
            expected_text: "hello".to_string(),
            language: "en".to_string(),
            audio_format: crate::services::stt::SttAudioFormat::Wav,
            page_id: None,
        };

        let result = SttViewSet::evaluate(user_id, request, vec![], &provider).await;

        assert!(matches!(result, EvaluateResult::InvalidInput(_)));
    }

    #[tokio::test]
    async fn test_evaluate_unsupported_language() {
        let user_id = Uuid::new_v4();
        let provider = MockSttProvider::new();
        let request = PronunciationRequest {
            expected_text: "test".to_string(),
            language: "xyz".to_string(),
            audio_format: crate::services::stt::SttAudioFormat::Wav,
            page_id: None,
        };

        let result =
            SttViewSet::evaluate(user_id, request, vec![0u8; 100], &provider).await;

        assert!(matches!(result, EvaluateResult::InvalidInput(_)));
    }

    #[test]
    fn test_list_attempts() {
        let user_id = Uuid::new_v4();
        let other_user_id = Uuid::new_v4();

        let attempts = vec![
            PronunciationAttempt::new(user_id, "hello".to_string(), "en".to_string()),
            PronunciationAttempt::new(user_id, "world".to_string(), "en".to_string()),
            PronunciationAttempt::new(other_user_id, "bonjour".to_string(), "fr".to_string()),
        ];

        let result = SttViewSet::list_attempts(user_id, &attempts);

        assert_eq!(result.len(), 2);
        assert!(result.iter().all(|a| a.expected_text != "bonjour"));
    }

    #[test]
    fn test_retrieve_attempt_success() {
        let user_id = Uuid::new_v4();
        let attempt =
            PronunciationAttempt::new(user_id, "test".to_string(), "en".to_string());

        let result = SttViewSet::retrieve(user_id, Some(&attempt));
        assert!(matches!(result, GetAttemptResult::Success(_)));
    }

    #[test]
    fn test_retrieve_attempt_unauthorized() {
        let owner_id = Uuid::new_v4();
        let other_id = Uuid::new_v4();
        let attempt =
            PronunciationAttempt::new(owner_id, "test".to_string(), "en".to_string());

        let result = SttViewSet::retrieve(other_id, Some(&attempt));
        assert!(matches!(result, GetAttemptResult::Unauthorized));
    }

    #[test]
    fn test_retrieve_attempt_not_found() {
        let user_id = Uuid::new_v4();
        let result = SttViewSet::retrieve(user_id, None);
        assert!(matches!(result, GetAttemptResult::NotFound));
    }

    #[test]
    fn test_stats_empty() {
        let user_id = Uuid::new_v4();
        let stats = SttViewSet::stats(user_id, &[], &[]);

        assert_eq!(stats.total_attempts, 0);
        assert_eq!(stats.average_score, 0.0);
        assert_eq!(stats.best_score, 0);
        assert!(stats.weak_words.is_empty());
    }

    #[test]
    fn test_stats_with_attempts() {
        let user_id = Uuid::new_v4();

        let mut attempt1 =
            PronunciationAttempt::new(user_id, "hello".to_string(), "en".to_string());
        attempt1.status = AttemptStatus::Completed;
        attempt1.overall_score = Some(80);

        let mut attempt2 =
            PronunciationAttempt::new(user_id, "world".to_string(), "en".to_string());
        attempt2.status = AttemptStatus::Completed;
        attempt2.overall_score = Some(90);

        let feedbacks = vec![WordFeedback {
            id: Uuid::new_v4(),
            attempt_id: attempt1.id,
            word: "hello".to_string(),
            score: 50,
            feedback: Some("Needs improvement".to_string()),
            start_ms: Some(0),
            end_ms: Some(300),
        }];

        let stats = SttViewSet::stats(user_id, &[attempt1, attempt2], &feedbacks);

        assert_eq!(stats.total_attempts, 2);
        assert_eq!(stats.average_score, 85.0);
        assert_eq!(stats.best_score, 90);
        assert_eq!(stats.weak_words.len(), 1);
        assert_eq!(stats.weak_words[0].word, "hello");
    }
}
