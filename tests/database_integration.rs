//! Database Integration Tests
//!
//! Uses TestContainers for PostgreSQL to test actual database operations.
//! Run with: cargo test --test database_integration -- --ignored

use sqlx::{PgPool, postgres::PgPoolOptions};
use testcontainers::{ContainerAsync, runners::AsyncRunner};
use testcontainers_modules::postgres::Postgres;
use uuid::Uuid;

/// Test helper to create a PostgreSQL container and connection pool
async fn setup_postgres() -> (ContainerAsync<Postgres>, PgPool) {
    let container = Postgres::default()
        .with_db_name("hailango_test")
        .with_user("test")
        .with_password("test")
        .start()
        .await
        .expect("Failed to start PostgreSQL container");

    let host_port = container
        .get_host_port_ipv4(5432)
        .await
        .expect("Failed to get PostgreSQL port");

    let database_url = format!(
        "postgresql://test:test@127.0.0.1:{}/hailango_test",
        host_port
    );

    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&database_url)
        .await
        .expect("Failed to connect to PostgreSQL");

    // Run schema creation (inline for testing without migrations)
    sqlx::query(
        r#"
        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            email VARCHAR(255) NOT NULL UNIQUE,
            password_hash VARCHAR(255),
            display_name VARCHAR(100),
            avatar_url TEXT,
            native_language VARCHAR(10) DEFAULT 'en',
            target_language VARCHAR(10),
            oauth_provider VARCHAR(50),
            oauth_id VARCHAR(255),
            is_active BOOLEAN NOT NULL DEFAULT TRUE,
            is_verified BOOLEAN NOT NULL DEFAULT FALSE,
            last_login_at TIMESTAMP WITH TIME ZONE,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
        );

        CREATE TABLE IF NOT EXISTS books (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            title VARCHAR(255) NOT NULL,
            author VARCHAR(255),
            source_language VARCHAR(10) NOT NULL,
            target_language VARCHAR(10) NOT NULL,
            total_pages INTEGER NOT NULL DEFAULT 0,
            status VARCHAR(50) NOT NULL DEFAULT 'pending',
            cover_image_url TEXT,
            original_file_url TEXT,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
        );

        CREATE TABLE IF NOT EXISTS pages (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
            page_number INTEGER NOT NULL,
            original_text TEXT,
            translated_text TEXT,
            audio_url TEXT,
            ocr_confidence REAL,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
            UNIQUE(book_id, page_number)
        );

        CREATE TABLE IF NOT EXISTS vocabulary (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            page_id UUID REFERENCES pages(id) ON DELETE SET NULL,
            word VARCHAR(255) NOT NULL,
            reading VARCHAR(255),
            meaning TEXT NOT NULL,
            context TEXT,
            notes TEXT,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
        );

        CREATE TABLE IF NOT EXISTS srs_schedule (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            vocabulary_id UUID NOT NULL UNIQUE REFERENCES vocabulary(id) ON DELETE CASCADE,
            next_review TIMESTAMP WITH TIME ZONE NOT NULL,
            interval_days INTEGER NOT NULL DEFAULT 1,
            easiness_factor REAL NOT NULL DEFAULT 2.5,
            repetitions INTEGER NOT NULL DEFAULT 0,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
        );

        CREATE TABLE IF NOT EXISTS learning_sessions (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
            session_type VARCHAR(50) NOT NULL DEFAULT 'reading',
            status VARCHAR(50) NOT NULL DEFAULT 'active',
            start_page INTEGER NOT NULL DEFAULT 1,
            end_page INTEGER,
            current_page INTEGER NOT NULL DEFAULT 1,
            started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
            ended_at TIMESTAMP WITH TIME ZONE,
            created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
        );
        "#,
    )
    .execute(&pool)
    .await
    .expect("Failed to create schema");

    (container, pool)
}

#[tokio::test]
#[ignore = "requires Docker"]
async fn test_user_crud_operations() {
    let (_container, pool) = setup_postgres().await;

    let user_id = Uuid::new_v4();
    let email = "test@example.com";
    let password_hash = "hashed_password_here";
    let display_name = "Test User";

    // Create user
    sqlx::query(
        r#"
        INSERT INTO users (id, email, password_hash, display_name)
        VALUES ($1, $2, $3, $4)
        "#,
    )
    .bind(user_id)
    .bind(email)
    .bind(password_hash)
    .bind(display_name)
    .execute(&pool)
    .await
    .expect("Failed to create user");

    // Read user
    let user: (Uuid, String, Option<String>, bool) = sqlx::query_as(
        r#"
        SELECT id, email, display_name, is_active
        FROM users WHERE id = $1
        "#,
    )
    .bind(user_id)
    .fetch_one(&pool)
    .await
    .expect("Failed to fetch user");

    assert_eq!(user.1, email);
    assert_eq!(user.2.as_deref(), Some(display_name));
    assert!(user.3);

    // Update user
    let new_name = "Updated User";
    sqlx::query(
        r#"
        UPDATE users SET display_name = $1 WHERE id = $2
        "#,
    )
    .bind(new_name)
    .bind(user_id)
    .execute(&pool)
    .await
    .expect("Failed to update user");

    let updated: (Option<String>,) =
        sqlx::query_as(r#"SELECT display_name FROM users WHERE id = $1"#)
            .bind(user_id)
            .fetch_one(&pool)
            .await
            .expect("Failed to fetch updated user");

    assert_eq!(updated.0.as_deref(), Some(new_name));

    // Delete user
    sqlx::query(r#"DELETE FROM users WHERE id = $1"#)
        .bind(user_id)
        .execute(&pool)
        .await
        .expect("Failed to delete user");

    let deleted: Option<(Uuid,)> = sqlx::query_as(r#"SELECT id FROM users WHERE id = $1"#)
        .bind(user_id)
        .fetch_optional(&pool)
        .await
        .expect("Failed to check deletion");

    assert!(deleted.is_none());
}

#[tokio::test]
#[ignore = "requires Docker"]
async fn test_book_with_pages() {
    let (_container, pool) = setup_postgres().await;

    // First create a user
    let user_id = Uuid::new_v4();
    sqlx::query(
        r#"
        INSERT INTO users (id, email, password_hash)
        VALUES ($1, 'bookowner@test.com', 'hash')
        "#,
    )
    .bind(user_id)
    .execute(&pool)
    .await
    .expect("Failed to create user");

    // Create book
    let book_id = Uuid::new_v4();
    sqlx::query(
        r#"
        INSERT INTO books (id, user_id, title, source_language, target_language)
        VALUES ($1, $2, 'Test Book', 'zh', 'en')
        "#,
    )
    .bind(book_id)
    .bind(user_id)
    .execute(&pool)
    .await
    .expect("Failed to create book");

    // Add pages
    for page_num in 1..=3i32 {
        let page_id = Uuid::new_v4();
        sqlx::query(
            r#"
            INSERT INTO pages (id, book_id, page_number, original_text, translated_text)
            VALUES ($1, $2, $3, $4, $5)
            "#,
        )
        .bind(page_id)
        .bind(book_id)
        .bind(page_num)
        .bind(format!("Original text page {}", page_num))
        .bind(format!("Translated text page {}", page_num))
        .execute(&pool)
        .await
        .expect("Failed to create page");
    }

    // Verify pages
    let pages: Vec<(i32, Option<String>)> = sqlx::query_as(
        r#"
        SELECT page_number, original_text
        FROM pages WHERE book_id = $1
        ORDER BY page_number
        "#,
    )
    .bind(book_id)
    .fetch_all(&pool)
    .await
    .expect("Failed to fetch pages");

    assert_eq!(pages.len(), 3);
    assert_eq!(pages[0].0, 1);
    assert_eq!(pages[2].0, 3);
}

#[tokio::test]
#[ignore = "requires Docker"]
async fn test_vocabulary_and_srs() {
    let (_container, pool) = setup_postgres().await;

    // Setup user, book, page
    let user_id = Uuid::new_v4();
    let book_id = Uuid::new_v4();
    let page_id = Uuid::new_v4();

    sqlx::query(
        r#"INSERT INTO users (id, email, password_hash) VALUES ($1, 'vocab@test.com', 'hash')"#,
    )
    .bind(user_id)
    .execute(&pool)
    .await
    .unwrap();

    sqlx::query(r#"INSERT INTO books (id, user_id, title, source_language, target_language) VALUES ($1, $2, 'Vocab Book', 'ja', 'en')"#)
        .bind(book_id)
        .bind(user_id)
        .execute(&pool)
        .await
        .unwrap();

    sqlx::query(
        r#"INSERT INTO pages (id, book_id, page_number, original_text) VALUES ($1, $2, 1, 'Test')"#,
    )
    .bind(page_id)
    .bind(book_id)
    .execute(&pool)
    .await
    .unwrap();

    // Add vocabulary
    let vocab_id = Uuid::new_v4();
    sqlx::query(
        r#"
        INSERT INTO vocabulary (id, user_id, page_id, word, reading, meaning)
        VALUES ($1, $2, $3, '食べる', 'たべる', 'to eat')
        "#,
    )
    .bind(vocab_id)
    .bind(user_id)
    .bind(page_id)
    .execute(&pool)
    .await
    .expect("Failed to create vocabulary");

    // Add SRS schedule
    let schedule_id = Uuid::new_v4();
    sqlx::query(
        r#"
        INSERT INTO srs_schedule (id, vocabulary_id, next_review, interval_days, easiness_factor, repetitions)
        VALUES ($1, $2, NOW(), 1, 2.5, 0)
        "#,
    )
    .bind(schedule_id)
    .bind(vocab_id)
    .execute(&pool)
    .await
    .expect("Failed to create SRS schedule");

    // Verify
    let schedule: (i32, f32, String) = sqlx::query_as(
        r#"
        SELECT s.interval_days, s.easiness_factor, v.word
        FROM srs_schedule s
        JOIN vocabulary v ON s.vocabulary_id = v.id
        WHERE s.id = $1
        "#,
    )
    .bind(schedule_id)
    .fetch_one(&pool)
    .await
    .expect("Failed to fetch schedule");

    assert_eq!(schedule.0, 1);
    assert!((schedule.1 - 2.5).abs() < 0.01);
    assert_eq!(schedule.2, "食べる");
}

#[tokio::test]
#[ignore = "requires Docker"]
async fn test_learning_session() {
    let (_container, pool) = setup_postgres().await;

    // Setup
    let user_id = Uuid::new_v4();
    let book_id = Uuid::new_v4();

    sqlx::query(
        r#"INSERT INTO users (id, email, password_hash) VALUES ($1, 'learner@test.com', 'hash')"#,
    )
    .bind(user_id)
    .execute(&pool)
    .await
    .unwrap();

    sqlx::query(r#"INSERT INTO books (id, user_id, title, source_language, target_language, total_pages) VALUES ($1, $2, 'Learning Book', 'ko', 'en', 10)"#)
        .bind(book_id)
        .bind(user_id)
        .execute(&pool)
        .await
        .unwrap();

    // Create learning session
    let session_id = Uuid::new_v4();
    sqlx::query(
        r#"
        INSERT INTO learning_sessions (id, user_id, book_id, session_type, start_page, end_page)
        VALUES ($1, $2, $3, 'reading', 1, 5)
        "#,
    )
    .bind(session_id)
    .bind(user_id)
    .bind(book_id)
    .execute(&pool)
    .await
    .expect("Failed to create session");

    // Update session progress
    sqlx::query(
        r#"
        UPDATE learning_sessions
        SET current_page = 3, status = 'in_progress'
        WHERE id = $1
        "#,
    )
    .bind(session_id)
    .execute(&pool)
    .await
    .expect("Failed to update session");

    // Verify
    let session: (i32, String, String) = sqlx::query_as(
        r#"
        SELECT current_page, status, session_type
        FROM learning_sessions WHERE id = $1
        "#,
    )
    .bind(session_id)
    .fetch_one(&pool)
    .await
    .expect("Failed to fetch session");

    assert_eq!(session.0, 3);
    assert_eq!(session.1, "in_progress");
    assert_eq!(session.2, "reading");
}
