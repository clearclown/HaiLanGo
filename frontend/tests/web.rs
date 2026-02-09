//! Frontend WASM integration tests

use wasm_bindgen_test::*;

wasm_bindgen_test_configure!(run_in_browser);

use hailango_frontend::api::*;

// ── API type serialization tests ──

#[wasm_bindgen_test]
fn test_login_payload_serialization() {
    #[derive(serde::Serialize)]
    struct LoginPayload<'a> {
        email: &'a str,
        password: &'a str,
    }

    let payload = LoginPayload {
        email: "test@example.com",
        password: "secret123",
    };
    let json = serde_json::to_string(&payload).unwrap();
    assert!(json.contains("test@example.com"));
    assert!(json.contains("secret123"));
}

#[wasm_bindgen_test]
fn test_auth_response_deserialization() {
    let json = r#"{
        "user": {
            "id": "user-1",
            "email": "test@example.com",
            "display_name": "Test User",
            "native_language": "en",
            "email_verified": true
        },
        "tokens": {
            "access_token": "token-abc",
            "refresh_token": "refresh-xyz",
            "expires_in": 3600
        }
    }"#;

    let auth: AuthResponse = serde_json::from_str(json).unwrap();
    assert_eq!(auth.user.id, "user-1");
    assert_eq!(auth.user.email, "test@example.com");
    assert_eq!(auth.user.display_name, "Test User");
    assert!(auth.user.email_verified);
    assert_eq!(auth.tokens.access_token, "token-abc");
    assert_eq!(auth.tokens.refresh_token, "refresh-xyz");
    assert_eq!(auth.tokens.expires_in, 3600);
}

#[wasm_bindgen_test]
fn test_auth_response_default_fields() {
    let json = r#"{
        "user": {
            "id": "u1",
            "email": "a@b.com",
            "display_name": "A"
        },
        "tokens": {
            "access_token": "t",
            "refresh_token": "r",
            "expires_in": 60
        }
    }"#;

    let auth: AuthResponse = serde_json::from_str(json).unwrap();
    assert_eq!(auth.user.native_language, "");
    assert!(!auth.user.email_verified);
}

#[wasm_bindgen_test]
fn test_book_item_deserialization() {
    let json = r#"{
        "id": "book-1",
        "title": "HSK Level 3",
        "author": "Hanban",
        "language": "zh",
        "total_pages": 200,
        "progress": 0.45
    }"#;

    let book: BookItem = serde_json::from_str(json).unwrap();
    assert_eq!(book.id, "book-1");
    assert_eq!(book.title, "HSK Level 3");
    assert_eq!(book.author, "Hanban");
    assert_eq!(book.language, "zh");
    assert_eq!(book.total_pages, 200);
    assert!((book.progress - 0.45).abs() < f32::EPSILON);
}

#[wasm_bindgen_test]
fn test_book_item_default_fields() {
    let json = r#"{
        "id": "b2",
        "title": "Test",
        "author": "Author",
        "total_pages": 10
    }"#;

    let book: BookItem = serde_json::from_str(json).unwrap();
    assert_eq!(book.language, "");
    assert!((book.progress - 0.0).abs() < f32::EPSILON);
}

#[wasm_bindgen_test]
fn test_page_content_deserialization() {
    let json = r#"{
        "page_number": 5,
        "text": "Hello World",
        "book_title": "Test Book",
        "total_pages": 100
    }"#;

    let page: PageContent = serde_json::from_str(json).unwrap();
    assert_eq!(page.page_number, 5);
    assert_eq!(page.text, "Hello World");
    assert_eq!(page.book_title, "Test Book");
    assert_eq!(page.total_pages, 100);
}

#[wasm_bindgen_test]
fn test_review_card_deserialization() {
    let json = r#"{
        "id": "card-1",
        "word": "你好",
        "reading": "nǐ hǎo",
        "meaning": "hello",
        "sentence": "你好世界"
    }"#;

    let card: ReviewCard = serde_json::from_str(json).unwrap();
    assert_eq!(card.id, "card-1");
    assert_eq!(card.word, "你好");
    assert_eq!(card.reading, "nǐ hǎo");
    assert_eq!(card.meaning, "hello");
    assert_eq!(card.sentence, Some("你好世界".to_string()));
}

#[wasm_bindgen_test]
fn test_review_card_optional_sentence() {
    let json = r#"{
        "id": "c2",
        "word": "学",
        "reading": "xué",
        "meaning": "study"
    }"#;

    let card: ReviewCard = serde_json::from_str(json).unwrap();
    assert!(card.sentence.is_none());
}

#[wasm_bindgen_test]
fn test_review_stats_deserialization() {
    let json = r#"{
        "total_cards": 150,
        "due_today": 25,
        "streak": 7,
        "accuracy": 0.85
    }"#;

    let stats: ReviewStats = serde_json::from_str(json).unwrap();
    assert_eq!(stats.total_cards, 150);
    assert_eq!(stats.due_today, 25);
    assert_eq!(stats.streak, 7);
    assert!((stats.accuracy - 0.85).abs() < f32::EPSILON);
}

#[wasm_bindgen_test]
fn test_submit_review_serialization() {
    let req = SubmitReviewRequest {
        card_id: "card-42".to_string(),
        rating: 3,
    };
    let json = serde_json::to_string(&req).unwrap();
    assert!(json.contains("card-42"));
    assert!(json.contains("3"));
}

#[wasm_bindgen_test]
fn test_tts_synthesize_request_serialization() {
    let req = TtsSynthesizeRequest {
        text: "你好世界".to_string(),
        language: "zh".to_string(),
    };
    let json = serde_json::to_string(&req).unwrap();
    assert!(json.contains("你好世界"));
    assert!(json.contains("zh"));
}

#[wasm_bindgen_test]
fn test_start_lesson_request_serialization() {
    let req = StartLessonRequest {
        book_id: "book-1".to_string(),
        start_page: 1,
        end_page: 10,
        speed: Some(1.5),
        page_interval: Some(5),
        repeat_count: None,
    };
    let json = serde_json::to_string(&req).unwrap();
    assert!(json.contains("book-1"));
    assert!(json.contains("1.5"));
    assert!(json.contains("\"page_interval\":5"));
    assert!(!json.contains("repeat_count"));
}

#[wasm_bindgen_test]
fn test_start_lesson_request_skip_none_fields() {
    let req = StartLessonRequest {
        book_id: "b1".to_string(),
        start_page: 1,
        end_page: 5,
        speed: None,
        page_interval: None,
        repeat_count: None,
    };
    let json = serde_json::to_string(&req).unwrap();
    assert!(!json.contains("speed"));
    assert!(!json.contains("page_interval"));
    assert!(!json.contains("repeat_count"));
}

#[wasm_bindgen_test]
fn test_update_teacher_config_serialization() {
    let config = UpdateTeacherConfig {
        speed: Some(2.0),
        page_interval: None,
        repeat_count: Some(3),
        auto_advance: Some(true),
    };
    let json = serde_json::to_string(&config).unwrap();
    assert!(json.contains("2.0"));
    assert!(json.contains("true"));
    assert!(json.contains("repeat_count"));
    assert!(!json.contains("page_interval"));
}

#[wasm_bindgen_test]
fn test_lesson_status_response_deserialization() {
    let json = r#"{
        "session_id": "sess-1",
        "status": "playing",
        "current_page": 3,
        "total_pages": 10,
        "pages_completed": 2
    }"#;

    let status: LessonStatusResponse = serde_json::from_str(json).unwrap();
    assert_eq!(status.session_id, "sess-1");
    assert_eq!(status.status, "playing");
    assert_eq!(status.current_page, 3);
    assert_eq!(status.total_pages, 10);
    assert_eq!(status.pages_completed, 2);
}

#[wasm_bindgen_test]
fn test_lesson_status_default_pages_completed() {
    let json = r#"{
        "session_id": "s1",
        "status": "idle",
        "current_page": 1,
        "total_pages": 5
    }"#;

    let status: LessonStatusResponse = serde_json::from_str(json).unwrap();
    assert_eq!(status.pages_completed, 0);
}

#[wasm_bindgen_test]
fn test_oauth_redirect_response_deserialization() {
    let json = r#"{
        "auth_url": "https://accounts.google.com/o/oauth2/v2/auth?...",
        "state": "random-state-123"
    }"#;

    let resp: OAuthRedirectResponse = serde_json::from_str(json).unwrap();
    assert!(resp.auth_url.starts_with("https://"));
    assert_eq!(resp.state, "random-state-123");
}

#[wasm_bindgen_test]
fn test_oauth_provider_info_deserialization() {
    let json = r#"{"name": "google", "configured": true}"#;
    let info: OAuthProviderInfo = serde_json::from_str(json).unwrap();
    assert_eq!(info.name, "google");
    assert!(info.configured);
}

#[wasm_bindgen_test]
fn test_learning_session_deserialization() {
    let json = r#"{
        "id": "ls-1",
        "book_id": "book-1",
        "book_title": "HSK 3",
        "current_page": 5,
        "total_pages": 100
    }"#;

    let session: LearningSession = serde_json::from_str(json).unwrap();
    assert_eq!(session.id, "ls-1");
    assert_eq!(session.book_id, "book-1");
    assert_eq!(session.book_title, "HSK 3");
    assert_eq!(session.current_page, 5);
    assert_eq!(session.total_pages, 100);
}

#[wasm_bindgen_test]
fn test_tts_language_deserialization() {
    let json = r#"{"code": "zh-CN", "name": "Chinese (Simplified)"}"#;
    let lang: TtsLanguage = serde_json::from_str(json).unwrap();
    assert_eq!(lang.code, "zh-CN");
    assert_eq!(lang.name, "Chinese (Simplified)");
}

#[wasm_bindgen_test]
fn test_review_submit_response_deserialization() {
    let json = r#"{
        "next_review": "2025-05-20T10:00:00Z",
        "interval_days": 4
    }"#;

    let resp: ReviewSubmitResponse = serde_json::from_str(json).unwrap();
    assert_eq!(resp.interval_days, 4);
    assert!(resp.next_review.contains("2025"));
}

// ── Token management tests ──

#[wasm_bindgen_test]
fn test_token_management() {
    // Initially no token
    ApiClient::clear_token();
    assert!(!ApiClient::is_authenticated());
    assert!(ApiClient::get_token().is_none());

    // Set token
    ApiClient::set_token("test-token-123");
    assert!(ApiClient::is_authenticated());
    assert_eq!(ApiClient::get_token(), Some("test-token-123".to_string()));

    // Clear token
    ApiClient::clear_token();
    assert!(!ApiClient::is_authenticated());
    assert!(ApiClient::get_token().is_none());
}

#[wasm_bindgen_test]
fn test_refresh_token_management() {
    ApiClient::clear_token();

    ApiClient::set_refresh_token("refresh-abc");
    // Refresh token alone doesn't count as authenticated
    assert!(!ApiClient::is_authenticated());

    ApiClient::set_token("access-xyz");
    assert!(ApiClient::is_authenticated());

    // Clear removes both
    ApiClient::clear_token();
    assert!(!ApiClient::is_authenticated());
}

// ── Component rendering tests ──

#[wasm_bindgen_test]
fn test_app_router_renders() {
    use leptos::*;
    use hailango_frontend::AppRouter;

    let _ = create_runtime();
    let view = AppRouter();
    // Verify the component produces a view without panicking
    let _ = view.into_view();
}
