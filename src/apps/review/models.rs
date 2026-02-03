//! SRS Review models - Vocabulary and SRS Schedule

use chrono::{DateTime, Duration, NaiveDate, Utc};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// Vocabulary word extracted from a page
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Vocabulary {
    pub id: Uuid,
    pub page_id: Uuid,
    pub user_id: Uuid,
    pub word: String,
    pub reading: Option<String>,
    pub meaning: String,
    pub part_of_speech: Option<String>,
    pub example_sentence: Option<String>,
    pub frequency: i32,
    pub created_at: DateTime<Utc>,
}

impl Vocabulary {
    pub fn new(page_id: Uuid, user_id: Uuid, word: String, meaning: String) -> Self {
        Self {
            id: Uuid::new_v4(),
            page_id,
            user_id,
            word,
            meaning,
            reading: None,
            part_of_speech: None,
            example_sentence: None,
            frequency: 1,
            created_at: Utc::now(),
        }
    }
}

/// SRS Schedule using SM-2 algorithm
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SrsSchedule {
    pub id: Uuid,
    pub user_id: Uuid,
    pub vocabulary_id: Uuid,
    pub next_review_date: NaiveDate,
    pub interval_days: i32,
    pub easiness_factor: f32,
    pub repetitions: i32,
    pub correct_count: i32,
    pub incorrect_count: i32,
    pub last_reviewed_at: Option<DateTime<Utc>>,
    pub created_at: DateTime<Utc>,
}

impl SrsSchedule {
    /// Create new SRS schedule for a vocabulary
    pub fn new(user_id: Uuid, vocabulary_id: Uuid) -> Self {
        Self {
            id: Uuid::new_v4(),
            user_id,
            vocabulary_id,
            next_review_date: Utc::now().date_naive(),
            interval_days: 1,
            easiness_factor: 2.5,
            repetitions: 0,
            correct_count: 0,
            incorrect_count: 0,
            last_reviewed_at: None,
            created_at: Utc::now(),
        }
    }

    /// SM-2 algorithm implementation
    /// quality: 0-5 (0-2 = fail, 3-5 = pass)
    pub fn update_after_review(&mut self, quality: u8) {
        let q = quality.min(5) as f32;

        // Update easiness factor (EF)
        // EF' = EF + (0.1 - (5 - q) * (0.08 + (5 - q) * 0.02))
        let ef_delta = 0.1 - (5.0 - q) * (0.08 + (5.0 - q) * 0.02);
        self.easiness_factor = (self.easiness_factor + ef_delta).max(1.3);

        if quality >= 3 {
            // Correct answer
            self.correct_count += 1;
            self.repetitions += 1;

            // Calculate new interval
            self.interval_days = match self.repetitions {
                1 => 1,
                2 => 6,
                _ => (self.interval_days as f32 * self.easiness_factor).round() as i32,
            };
        } else {
            // Incorrect answer - reset
            self.incorrect_count += 1;
            self.repetitions = 0;
            self.interval_days = 1;
        }

        // Update next review date
        self.next_review_date = Utc::now().date_naive() + Duration::days(self.interval_days as i64);
        self.last_reviewed_at = Some(Utc::now());
    }

    /// Check if review is due today
    pub fn is_due(&self) -> bool {
        self.next_review_date <= Utc::now().date_naive()
    }

    /// Get retention rate
    pub fn retention_rate(&self) -> f32 {
        let total = self.correct_count + self.incorrect_count;
        if total == 0 {
            0.0
        } else {
            self.correct_count as f32 / total as f32
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_create_vocabulary() {
        let page_id = Uuid::new_v4();
        let user_id = Uuid::new_v4();
        let vocab = Vocabulary::new(
            page_id,
            user_id,
            "hello".to_string(),
            "greeting".to_string(),
        );

        assert_eq!(vocab.word, "hello");
        assert_eq!(vocab.meaning, "greeting");
        assert_eq!(vocab.frequency, 1);
    }

    #[test]
    fn test_create_srs_schedule() {
        let user_id = Uuid::new_v4();
        let vocab_id = Uuid::new_v4();
        let schedule = SrsSchedule::new(user_id, vocab_id);

        assert_eq!(schedule.interval_days, 1);
        assert_eq!(schedule.easiness_factor, 2.5);
        assert_eq!(schedule.repetitions, 0);
    }

    #[test]
    fn test_sm2_perfect_recall() {
        let mut schedule = SrsSchedule::new(Uuid::new_v4(), Uuid::new_v4());

        // First review - perfect (5)
        schedule.update_after_review(5);
        assert_eq!(schedule.interval_days, 1);
        assert_eq!(schedule.repetitions, 1);
        assert_eq!(schedule.correct_count, 1);

        // Second review - perfect (5)
        schedule.update_after_review(5);
        assert_eq!(schedule.interval_days, 6);
        assert_eq!(schedule.repetitions, 2);

        // Third review - perfect (5)
        schedule.update_after_review(5);
        assert!(schedule.interval_days > 6);
        assert_eq!(schedule.repetitions, 3);
    }

    #[test]
    fn test_sm2_failed_review_resets() {
        let mut schedule = SrsSchedule::new(Uuid::new_v4(), Uuid::new_v4());

        // Build up some progress
        schedule.update_after_review(5);
        schedule.update_after_review(5);
        schedule.update_after_review(5);

        let old_interval = schedule.interval_days;
        assert!(old_interval > 1);

        // Fail the next review (quality = 1)
        schedule.update_after_review(1);

        assert_eq!(schedule.interval_days, 1); // Reset to 1
        assert_eq!(schedule.repetitions, 0); // Reset repetitions
        assert_eq!(schedule.incorrect_count, 1);
    }

    #[test]
    fn test_sm2_easiness_factor_bounds() {
        let mut schedule = SrsSchedule::new(Uuid::new_v4(), Uuid::new_v4());

        // Many failures should not drop EF below 1.3
        for _ in 0..20 {
            schedule.update_after_review(0);
        }

        assert!(schedule.easiness_factor >= 1.3);
    }

    #[test]
    fn test_is_due() {
        let mut schedule = SrsSchedule::new(Uuid::new_v4(), Uuid::new_v4());

        // Should be due immediately after creation
        assert!(schedule.is_due());

        // After a perfect review, should not be due immediately
        schedule.update_after_review(5);
        schedule.update_after_review(5);
        // interval_days is now 6, so next_review_date is in the future
        assert!(!schedule.is_due() || schedule.interval_days == 1);
    }

    #[test]
    fn test_retention_rate() {
        let mut schedule = SrsSchedule::new(Uuid::new_v4(), Uuid::new_v4());

        // Initially 0%
        assert_eq!(schedule.retention_rate(), 0.0);

        // 3 correct, 1 incorrect = 75%
        schedule.update_after_review(5);
        schedule.update_after_review(5);
        schedule.update_after_review(5);
        schedule.update_after_review(1);

        assert!((schedule.retention_rate() - 0.75).abs() < 0.01);
    }
}
