//! Review ViewSet - SRS vocabulary review endpoints

use uuid::Uuid;

use super::dto::{
    BulkReviewRequest, BulkReviewResultResponse, CreateVocabularyRequest, RecordReviewRequest,
    ReviewItemResponse, ReviewQueueResponse, ReviewResultResponse, ReviewStatsResponse,
    SrsScheduleResponse, VocabularyResponse,
};
use super::models::{SrsSchedule, Vocabulary};

/// Result for vocabulary creation
#[derive(Debug)]
pub enum CreateVocabularyResult {
    Success(VocabularyResponse),
    PageNotFound,
    DuplicateWord,
    InvalidInput(String),
}

/// Result for recording review
#[derive(Debug)]
pub enum RecordReviewResult {
    Success(ReviewResultResponse),
    VocabularyNotFound,
    ScheduleNotFound,
    InvalidQuality,
}

/// Result for getting review queue
#[derive(Debug)]
pub enum ReviewQueueResult {
    Success(ReviewQueueResponse),
    Empty,
}

/// Review ViewSet - handles SRS review endpoints
pub struct ReviewViewSet;

impl ReviewViewSet {
    /// Add a new vocabulary word
    pub fn create_vocabulary(
        request: CreateVocabularyRequest,
        user_id: Uuid,
        page_exists: bool,
        word_exists: bool,
    ) -> CreateVocabularyResult {
        if !page_exists {
            return CreateVocabularyResult::PageNotFound;
        }

        if word_exists {
            return CreateVocabularyResult::DuplicateWord;
        }

        // Validate input
        if request.word.trim().is_empty() {
            return CreateVocabularyResult::InvalidInput("Word cannot be empty".to_string());
        }

        if request.meaning.trim().is_empty() {
            return CreateVocabularyResult::InvalidInput("Meaning cannot be empty".to_string());
        }

        // Create vocabulary
        let mut vocab = Vocabulary::new(
            request.page_id,
            user_id,
            request.word.trim().to_string(),
            request.meaning.trim().to_string(),
        );

        vocab.reading = request.reading;
        vocab.part_of_speech = request.part_of_speech;
        vocab.example_sentence = request.example_sentence;

        CreateVocabularyResult::Success(Self::vocab_to_response(&vocab))
    }

    /// Record a review result
    pub fn record_review(
        request: RecordReviewRequest,
        schedule: Option<&mut SrsSchedule>,
        user_id: Uuid,
    ) -> RecordReviewResult {
        // Validate quality
        if request.quality > 5 {
            return RecordReviewResult::InvalidQuality;
        }

        let schedule = match schedule {
            Some(s) => s,
            None => return RecordReviewResult::ScheduleNotFound,
        };

        // Verify ownership
        if schedule.user_id != user_id {
            return RecordReviewResult::ScheduleNotFound;
        }

        // Update schedule with SM-2 algorithm
        schedule.update_after_review(request.quality);

        RecordReviewResult::Success(ReviewResultResponse {
            vocabulary_id: schedule.vocabulary_id,
            was_correct: request.quality >= 3,
            new_interval_days: schedule.interval_days,
            next_review_date: schedule.next_review_date,
        })
    }

    /// Record multiple reviews at once
    pub fn record_bulk_reviews(
        request: BulkReviewRequest,
        schedules: &mut [(Uuid, SrsSchedule)],
        user_id: Uuid,
    ) -> BulkReviewResultResponse {
        let mut results = Vec::new();
        let mut correct_count = 0;
        let mut incorrect_count = 0;

        for review in request.reviews {
            // Find matching schedule
            if let Some((_, schedule)) = schedules
                .iter_mut()
                .find(|(vocab_id, s)| *vocab_id == review.vocabulary_id && s.user_id == user_id)
            {
                if review.quality <= 5 {
                    let was_correct = review.quality >= 3;
                    schedule.update_after_review(review.quality);

                    if was_correct {
                        correct_count += 1;
                    } else {
                        incorrect_count += 1;
                    }

                    results.push(ReviewResultResponse {
                        vocabulary_id: review.vocabulary_id,
                        was_correct,
                        new_interval_days: schedule.interval_days,
                        next_review_date: schedule.next_review_date,
                    });
                }
            }
        }

        BulkReviewResultResponse {
            results,
            correct_count,
            incorrect_count,
        }
    }

    /// Get review queue (due items)
    pub fn get_review_queue(
        vocabularies: &[Vocabulary],
        schedules: &[SrsSchedule],
        user_id: Uuid,
        limit: usize,
    ) -> ReviewQueueResult {
        // Filter user's due schedules
        let due_schedules: Vec<_> = schedules
            .iter()
            .filter(|s| s.user_id == user_id && s.is_due())
            .collect();

        let due_count = due_schedules.len();
        let total_vocabulary = vocabularies.iter().filter(|v| v.user_id == user_id).count();

        if due_schedules.is_empty() {
            return ReviewQueueResult::Success(ReviewQueueResponse {
                items: vec![],
                due_count: 0,
                total_vocabulary,
            });
        }

        // Build review items
        let items: Vec<_> = due_schedules
            .iter()
            .take(limit)
            .filter_map(|schedule| {
                vocabularies
                    .iter()
                    .find(|v| v.id == schedule.vocabulary_id)
                    .map(|vocab| ReviewItemResponse {
                        vocabulary: Self::vocab_to_response(vocab),
                        schedule: Self::schedule_to_response(schedule),
                    })
            })
            .collect();

        ReviewQueueResult::Success(ReviewQueueResponse {
            items,
            due_count,
            total_vocabulary,
        })
    }

    /// Get review statistics for a user
    pub fn get_stats(
        vocabularies: &[Vocabulary],
        schedules: &[SrsSchedule],
        user_id: Uuid,
    ) -> ReviewStatsResponse {
        let user_vocabs: Vec<_> = vocabularies
            .iter()
            .filter(|v| v.user_id == user_id)
            .collect();

        let user_schedules: Vec<_> = schedules.iter().filter(|s| s.user_id == user_id).collect();

        let due_today = user_schedules.iter().filter(|s| s.is_due()).count();

        // Count overdue (due date < today, i.e., more than 0 days ago)
        let today = chrono::Utc::now().date_naive();
        let overdue = user_schedules
            .iter()
            .filter(|s| s.next_review_date < today)
            .count();

        // Learned = items with repetitions > 0
        let learned_count = user_schedules.iter().filter(|s| s.repetitions > 0).count();

        // Calculate average retention rate
        let total_retention: f32 = user_schedules.iter().map(|s| s.retention_rate()).sum();

        let average_retention_rate = if user_schedules.is_empty() {
            0.0
        } else {
            total_retention / user_schedules.len() as f32
        };

        ReviewStatsResponse {
            total_vocabulary: user_vocabs.len(),
            due_today,
            overdue,
            learned_count,
            average_retention_rate,
            streak_days: 0, // TODO: Implement streak tracking
        }
    }

    /// Get vocabulary by ID
    pub fn retrieve_vocabulary(
        vocabulary: Option<&Vocabulary>,
        user_id: Uuid,
    ) -> Option<VocabularyResponse> {
        vocabulary
            .filter(|v| v.user_id == user_id)
            .map(Self::vocab_to_response)
    }

    /// List vocabularies for a user
    pub fn list_vocabularies(
        vocabularies: &[Vocabulary],
        user_id: Uuid,
    ) -> Vec<VocabularyResponse> {
        vocabularies
            .iter()
            .filter(|v| v.user_id == user_id)
            .map(Self::vocab_to_response)
            .collect()
    }

    /// Convert vocabulary to response DTO
    fn vocab_to_response(vocab: &Vocabulary) -> VocabularyResponse {
        VocabularyResponse {
            id: vocab.id,
            page_id: vocab.page_id,
            word: vocab.word.clone(),
            reading: vocab.reading.clone(),
            meaning: vocab.meaning.clone(),
            part_of_speech: vocab.part_of_speech.clone(),
            example_sentence: vocab.example_sentence.clone(),
            frequency: vocab.frequency,
            created_at: vocab.created_at,
        }
    }

    /// Convert schedule to response DTO
    fn schedule_to_response(schedule: &SrsSchedule) -> SrsScheduleResponse {
        SrsScheduleResponse {
            id: schedule.id,
            vocabulary_id: schedule.vocabulary_id,
            next_review_date: schedule.next_review_date,
            interval_days: schedule.interval_days,
            easiness_factor: schedule.easiness_factor,
            repetitions: schedule.repetitions,
            retention_rate: schedule.retention_rate(),
            is_due: schedule.is_due(),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_create_vocabulary_success() {
        let user_id = Uuid::new_v4();
        let request = CreateVocabularyRequest {
            page_id: Uuid::new_v4(),
            word: "hello".to_string(),
            meaning: "a greeting".to_string(),
            reading: None,
            part_of_speech: Some("noun".to_string()),
            example_sentence: None,
        };

        let result = ReviewViewSet::create_vocabulary(request, user_id, true, false);

        match result {
            CreateVocabularyResult::Success(response) => {
                assert_eq!(response.word, "hello");
                assert_eq!(response.meaning, "a greeting");
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_create_vocabulary_page_not_found() {
        let request = CreateVocabularyRequest {
            page_id: Uuid::new_v4(),
            word: "test".to_string(),
            meaning: "test".to_string(),
            reading: None,
            part_of_speech: None,
            example_sentence: None,
        };

        let result = ReviewViewSet::create_vocabulary(request, Uuid::new_v4(), false, false);
        assert!(matches!(result, CreateVocabularyResult::PageNotFound));
    }

    #[test]
    fn test_create_vocabulary_duplicate() {
        let request = CreateVocabularyRequest {
            page_id: Uuid::new_v4(),
            word: "existing".to_string(),
            meaning: "already exists".to_string(),
            reading: None,
            part_of_speech: None,
            example_sentence: None,
        };

        let result = ReviewViewSet::create_vocabulary(request, Uuid::new_v4(), true, true);
        assert!(matches!(result, CreateVocabularyResult::DuplicateWord));
    }

    #[test]
    fn test_create_vocabulary_empty_word() {
        let request = CreateVocabularyRequest {
            page_id: Uuid::new_v4(),
            word: "   ".to_string(),
            meaning: "valid meaning".to_string(),
            reading: None,
            part_of_speech: None,
            example_sentence: None,
        };

        let result = ReviewViewSet::create_vocabulary(request, Uuid::new_v4(), true, false);
        assert!(matches!(result, CreateVocabularyResult::InvalidInput(_)));
    }

    #[test]
    fn test_record_review_success() {
        let user_id = Uuid::new_v4();
        let vocab_id = Uuid::new_v4();
        let mut schedule = SrsSchedule::new(user_id, vocab_id);

        let request = RecordReviewRequest {
            vocabulary_id: vocab_id,
            quality: 5,
        };

        let result = ReviewViewSet::record_review(request, Some(&mut schedule), user_id);

        match result {
            RecordReviewResult::Success(response) => {
                assert!(response.was_correct);
                assert_eq!(response.new_interval_days, 1);
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_record_review_fail() {
        let user_id = Uuid::new_v4();
        let vocab_id = Uuid::new_v4();
        let mut schedule = SrsSchedule::new(user_id, vocab_id);

        let request = RecordReviewRequest {
            vocabulary_id: vocab_id,
            quality: 2, // Fail
        };

        let result = ReviewViewSet::record_review(request, Some(&mut schedule), user_id);

        match result {
            RecordReviewResult::Success(response) => {
                assert!(!response.was_correct);
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_record_review_invalid_quality() {
        let user_id = Uuid::new_v4();
        let mut schedule = SrsSchedule::new(user_id, Uuid::new_v4());

        let request = RecordReviewRequest {
            vocabulary_id: Uuid::new_v4(),
            quality: 6, // Invalid
        };

        let result = ReviewViewSet::record_review(request, Some(&mut schedule), user_id);
        assert!(matches!(result, RecordReviewResult::InvalidQuality));
    }

    #[test]
    fn test_record_bulk_reviews() {
        let user_id = Uuid::new_v4();
        let vocab1 = Uuid::new_v4();
        let vocab2 = Uuid::new_v4();

        let mut schedules = vec![
            (vocab1, SrsSchedule::new(user_id, vocab1)),
            (vocab2, SrsSchedule::new(user_id, vocab2)),
        ];

        let request = BulkReviewRequest {
            reviews: vec![
                RecordReviewRequest {
                    vocabulary_id: vocab1,
                    quality: 5,
                },
                RecordReviewRequest {
                    vocabulary_id: vocab2,
                    quality: 2,
                },
            ],
        };

        let result = ReviewViewSet::record_bulk_reviews(request, &mut schedules, user_id);

        assert_eq!(result.correct_count, 1);
        assert_eq!(result.incorrect_count, 1);
        assert_eq!(result.results.len(), 2);
    }

    #[test]
    fn test_get_review_queue() {
        let user_id = Uuid::new_v4();
        let page_id = Uuid::new_v4();

        let vocab1 = Vocabulary::new(
            page_id,
            user_id,
            "word1".to_string(),
            "meaning1".to_string(),
        );
        let vocab2 = Vocabulary::new(
            page_id,
            user_id,
            "word2".to_string(),
            "meaning2".to_string(),
        );

        let schedule1 = SrsSchedule::new(user_id, vocab1.id);
        let schedule2 = SrsSchedule::new(user_id, vocab2.id);

        let vocabularies = vec![vocab1, vocab2];
        let schedules = vec![schedule1, schedule2];

        let result = ReviewViewSet::get_review_queue(&vocabularies, &schedules, user_id, 10);

        match result {
            ReviewQueueResult::Success(response) => {
                assert_eq!(response.due_count, 2);
                assert_eq!(response.total_vocabulary, 2);
                assert_eq!(response.items.len(), 2);
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_get_review_queue_empty() {
        let user_id = Uuid::new_v4();

        let result = ReviewViewSet::get_review_queue(&[], &[], user_id, 10);

        match result {
            ReviewQueueResult::Success(response) => {
                assert_eq!(response.due_count, 0);
                assert!(response.items.is_empty());
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_get_stats() {
        let user_id = Uuid::new_v4();
        let page_id = Uuid::new_v4();

        let vocab = Vocabulary::new(page_id, user_id, "word".to_string(), "meaning".to_string());
        let mut schedule = SrsSchedule::new(user_id, vocab.id);
        schedule.update_after_review(5); // Mark as learned

        let vocabularies = vec![vocab];
        let schedules = vec![schedule];

        let stats = ReviewViewSet::get_stats(&vocabularies, &schedules, user_id);

        assert_eq!(stats.total_vocabulary, 1);
        assert_eq!(stats.learned_count, 1);
    }

    #[test]
    fn test_retrieve_vocabulary() {
        let user_id = Uuid::new_v4();
        let vocab = Vocabulary::new(
            Uuid::new_v4(),
            user_id,
            "test".to_string(),
            "test".to_string(),
        );

        let result = ReviewViewSet::retrieve_vocabulary(Some(&vocab), user_id);
        assert!(result.is_some());

        // Wrong user
        let result = ReviewViewSet::retrieve_vocabulary(Some(&vocab), Uuid::new_v4());
        assert!(result.is_none());
    }

    #[test]
    fn test_list_vocabularies() {
        let user_id = Uuid::new_v4();
        let other_user = Uuid::new_v4();
        let page_id = Uuid::new_v4();

        let vocabularies = vec![
            Vocabulary::new(
                page_id,
                user_id,
                "word1".to_string(),
                "meaning1".to_string(),
            ),
            Vocabulary::new(
                page_id,
                user_id,
                "word2".to_string(),
                "meaning2".to_string(),
            ),
            Vocabulary::new(
                page_id,
                other_user,
                "word3".to_string(),
                "meaning3".to_string(),
            ),
        ];

        let result = ReviewViewSet::list_vocabularies(&vocabularies, user_id);
        assert_eq!(result.len(), 2);
    }

    #[test]
    fn test_create_vocabulary_with_full_data() {
        let user_id = Uuid::new_v4();
        let request = CreateVocabularyRequest {
            page_id: Uuid::new_v4(),
            word: "読む".to_string(),
            meaning: "to read".to_string(),
            reading: Some("よむ".to_string()),
            part_of_speech: Some("verb".to_string()),
            example_sentence: Some("本を読む".to_string()),
        };

        let result = ReviewViewSet::create_vocabulary(request, user_id, true, false);

        match result {
            CreateVocabularyResult::Success(response) => {
                assert_eq!(response.word, "読む");
                assert_eq!(response.reading, Some("よむ".to_string()));
                assert_eq!(response.part_of_speech, Some("verb".to_string()));
            }
            _ => panic!("Expected success"),
        }
    }

    #[test]
    fn test_record_review_wrong_user() {
        let owner = Uuid::new_v4();
        let attacker = Uuid::new_v4();
        let mut schedule = SrsSchedule::new(owner, Uuid::new_v4());

        let request = RecordReviewRequest {
            vocabulary_id: Uuid::new_v4(),
            quality: 5,
        };

        let result = ReviewViewSet::record_review(request, Some(&mut schedule), attacker);
        assert!(matches!(result, RecordReviewResult::ScheduleNotFound));
    }

    #[test]
    fn test_get_stats_empty() {
        let user_id = Uuid::new_v4();
        let stats = ReviewViewSet::get_stats(&[], &[], user_id);

        assert_eq!(stats.total_vocabulary, 0);
        assert_eq!(stats.due_today, 0);
        assert_eq!(stats.average_retention_rate, 0.0);
    }
}
