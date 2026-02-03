//! Internationalization (i18n) Service
//!
//! Provides multi-language support for the language learning platform.

use std::collections::HashMap;

/// Supported languages for the application
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum Language {
    English,
    Japanese,
    Chinese,
    Korean,
    Spanish,
    French,
    German,
}

impl Language {
    /// Get the ISO 639-1 code
    pub fn code(&self) -> &'static str {
        match self {
            Language::English => "en",
            Language::Japanese => "ja",
            Language::Chinese => "zh",
            Language::Korean => "ko",
            Language::Spanish => "es",
            Language::French => "fr",
            Language::German => "de",
        }
    }

    /// Get the language name
    pub fn name(&self) -> &'static str {
        match self {
            Language::English => "English",
            Language::Japanese => "Japanese",
            Language::Chinese => "Chinese",
            Language::Korean => "Korean",
            Language::Spanish => "Spanish",
            Language::French => "French",
            Language::German => "German",
        }
    }

    /// Parse from ISO code
    pub fn from_code(code: &str) -> Option<Self> {
        match code.to_lowercase().as_str() {
            "en" | "en-us" | "en-gb" => Some(Language::English),
            "ja" | "ja-jp" => Some(Language::Japanese),
            "zh" | "zh-cn" | "zh-tw" => Some(Language::Chinese),
            "ko" | "ko-kr" => Some(Language::Korean),
            "es" | "es-es" => Some(Language::Spanish),
            "fr" | "fr-fr" => Some(Language::French),
            "de" | "de-de" => Some(Language::German),
            _ => None,
        }
    }

    /// Get all supported languages
    pub fn all() -> Vec<Self> {
        vec![
            Language::English,
            Language::Japanese,
            Language::Chinese,
            Language::Korean,
            Language::Spanish,
            Language::French,
            Language::German,
        ]
    }
}

/// Translation keys for UI messages
pub mod keys {
    pub const WELCOME: &str = "welcome";
    pub const LOGIN: &str = "login";
    pub const REGISTER: &str = "register";
    pub const LOGOUT: &str = "logout";
    pub const EMAIL: &str = "email";
    pub const PASSWORD: &str = "password";
    pub const SUBMIT: &str = "submit";
    pub const CANCEL: &str = "cancel";
    pub const SAVE: &str = "save";
    pub const DELETE: &str = "delete";
    pub const CONFIRM: &str = "confirm";
    pub const LOADING: &str = "loading";
    pub const ERROR: &str = "error";
    pub const SUCCESS: &str = "success";
    pub const BOOKS: &str = "books";
    pub const LEARN: &str = "learn";
    pub const REVIEW: &str = "review";
    pub const PROFILE: &str = "profile";
    pub const SETTINGS: &str = "settings";
    pub const START_LEARNING: &str = "start_learning";
    pub const UPLOAD_BOOK: &str = "upload_book";
    pub const NO_BOOKS: &str = "no_books";
    pub const WORDS_LEARNED: &str = "words_learned";
    pub const STUDY_STREAK: &str = "study_streak";
    pub const REVIEW_QUEUE: &str = "review_queue";
    pub const SHOW_ANSWER: &str = "show_answer";
    pub const AGAIN: &str = "again";
    pub const HARD: &str = "hard";
    pub const GOOD: &str = "good";
    pub const EASY: &str = "easy";
}

/// Get translations for a language
pub fn get_translations(lang: Language) -> HashMap<&'static str, &'static str> {
    let mut translations = HashMap::new();

    match lang {
        Language::English => {
            translations.insert(keys::WELCOME, "Welcome to HaiLanGo");
            translations.insert(keys::LOGIN, "Sign In");
            translations.insert(keys::REGISTER, "Sign Up");
            translations.insert(keys::LOGOUT, "Sign Out");
            translations.insert(keys::EMAIL, "Email");
            translations.insert(keys::PASSWORD, "Password");
            translations.insert(keys::SUBMIT, "Submit");
            translations.insert(keys::CANCEL, "Cancel");
            translations.insert(keys::SAVE, "Save");
            translations.insert(keys::DELETE, "Delete");
            translations.insert(keys::CONFIRM, "Confirm");
            translations.insert(keys::LOADING, "Loading...");
            translations.insert(keys::ERROR, "An error occurred");
            translations.insert(keys::SUCCESS, "Success!");
            translations.insert(keys::BOOKS, "Books");
            translations.insert(keys::LEARN, "Learn");
            translations.insert(keys::REVIEW, "Review");
            translations.insert(keys::PROFILE, "Profile");
            translations.insert(keys::SETTINGS, "Settings");
            translations.insert(keys::START_LEARNING, "Start Learning");
            translations.insert(keys::UPLOAD_BOOK, "Upload Book");
            translations.insert(keys::NO_BOOKS, "No books yet");
            translations.insert(keys::WORDS_LEARNED, "Words Learned");
            translations.insert(keys::STUDY_STREAK, "Study Streak");
            translations.insert(keys::REVIEW_QUEUE, "Review Queue");
            translations.insert(keys::SHOW_ANSWER, "Show Answer");
            translations.insert(keys::AGAIN, "Again");
            translations.insert(keys::HARD, "Hard");
            translations.insert(keys::GOOD, "Good");
            translations.insert(keys::EASY, "Easy");
        }
        Language::Japanese => {
            translations.insert(keys::WELCOME, "HaiLanGoへようこそ");
            translations.insert(keys::LOGIN, "ログイン");
            translations.insert(keys::REGISTER, "新規登録");
            translations.insert(keys::LOGOUT, "ログアウト");
            translations.insert(keys::EMAIL, "メールアドレス");
            translations.insert(keys::PASSWORD, "パスワード");
            translations.insert(keys::SUBMIT, "送信");
            translations.insert(keys::CANCEL, "キャンセル");
            translations.insert(keys::SAVE, "保存");
            translations.insert(keys::DELETE, "削除");
            translations.insert(keys::CONFIRM, "確認");
            translations.insert(keys::LOADING, "読み込み中...");
            translations.insert(keys::ERROR, "エラーが発生しました");
            translations.insert(keys::SUCCESS, "成功しました！");
            translations.insert(keys::BOOKS, "本");
            translations.insert(keys::LEARN, "学習");
            translations.insert(keys::REVIEW, "復習");
            translations.insert(keys::PROFILE, "プロフィール");
            translations.insert(keys::SETTINGS, "設定");
            translations.insert(keys::START_LEARNING, "学習を開始");
            translations.insert(keys::UPLOAD_BOOK, "本をアップロード");
            translations.insert(keys::NO_BOOKS, "本がありません");
            translations.insert(keys::WORDS_LEARNED, "習得した単語");
            translations.insert(keys::STUDY_STREAK, "連続学習日数");
            translations.insert(keys::REVIEW_QUEUE, "復習キュー");
            translations.insert(keys::SHOW_ANSWER, "答えを見る");
            translations.insert(keys::AGAIN, "もう一度");
            translations.insert(keys::HARD, "難しい");
            translations.insert(keys::GOOD, "良い");
            translations.insert(keys::EASY, "簡単");
        }
        Language::Chinese => {
            translations.insert(keys::WELCOME, "欢迎使用HaiLanGo");
            translations.insert(keys::LOGIN, "登录");
            translations.insert(keys::REGISTER, "注册");
            translations.insert(keys::LOGOUT, "退出");
            translations.insert(keys::EMAIL, "邮箱");
            translations.insert(keys::PASSWORD, "密码");
            translations.insert(keys::SUBMIT, "提交");
            translations.insert(keys::CANCEL, "取消");
            translations.insert(keys::SAVE, "保存");
            translations.insert(keys::DELETE, "删除");
            translations.insert(keys::CONFIRM, "确认");
            translations.insert(keys::LOADING, "加载中...");
            translations.insert(keys::ERROR, "发生错误");
            translations.insert(keys::SUCCESS, "成功！");
            translations.insert(keys::BOOKS, "书籍");
            translations.insert(keys::LEARN, "学习");
            translations.insert(keys::REVIEW, "复习");
            translations.insert(keys::PROFILE, "个人资料");
            translations.insert(keys::SETTINGS, "设置");
            translations.insert(keys::START_LEARNING, "开始学习");
            translations.insert(keys::UPLOAD_BOOK, "上传书籍");
            translations.insert(keys::NO_BOOKS, "暂无书籍");
            translations.insert(keys::WORDS_LEARNED, "已学单词");
            translations.insert(keys::STUDY_STREAK, "连续学习天数");
            translations.insert(keys::REVIEW_QUEUE, "复习队列");
            translations.insert(keys::SHOW_ANSWER, "显示答案");
            translations.insert(keys::AGAIN, "重来");
            translations.insert(keys::HARD, "困难");
            translations.insert(keys::GOOD, "良好");
            translations.insert(keys::EASY, "简单");
        }
        _ => {
            // Fall back to English for other languages
            return get_translations(Language::English);
        }
    }

    translations
}

/// Translate a key for a given language
pub fn translate(lang: Language, key: &'static str) -> &'static str {
    get_translations(lang).get(key).copied().unwrap_or(key)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_language_codes() {
        assert_eq!(Language::English.code(), "en");
        assert_eq!(Language::Japanese.code(), "ja");
        assert_eq!(Language::Chinese.code(), "zh");
    }

    #[test]
    fn test_language_from_code() {
        assert_eq!(Language::from_code("en"), Some(Language::English));
        assert_eq!(Language::from_code("ja"), Some(Language::Japanese));
        assert_eq!(Language::from_code("zh-cn"), Some(Language::Chinese));
        assert_eq!(Language::from_code("invalid"), None);
    }

    #[test]
    fn test_translations() {
        let en = get_translations(Language::English);
        let ja = get_translations(Language::Japanese);

        assert_eq!(en.get(keys::LOGIN), Some(&"Sign In"));
        assert_eq!(ja.get(keys::LOGIN), Some(&"ログイン"));
    }

    #[test]
    fn test_translate_function() {
        assert_eq!(
            translate(Language::English, keys::WELCOME),
            "Welcome to HaiLanGo"
        );
        assert_eq!(
            translate(Language::Japanese, keys::WELCOME),
            "HaiLanGoへようこそ"
        );
    }

    #[test]
    fn test_all_languages() {
        let languages = Language::all();
        assert_eq!(languages.len(), 7);
        assert!(languages.contains(&Language::English));
        assert!(languages.contains(&Language::Japanese));
    }
}
