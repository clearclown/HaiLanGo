package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	var command string
	flag.StringVar(&command, "command", "up", "Migration command: up, down, version")
	flag.Parse()

	// データベース接続文字列を環境変数から取得
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://HaiLanGo:password@localhost:5432/HaiLanGo_dev?sslmode=disable"
		log.Printf("DATABASE_URL not set, using default: %s", dbURL)
	}

	// データベース接続
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 接続確認
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✓ Database connection successful")

	// マイグレーション実行
	switch command {
	case "up":
		if err := migrateUp(db); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("✓ All migrations applied successfully")
	case "down":
		if err := migrateDown(db); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("✓ All migrations rolled back successfully")
	case "version":
		version, err := getVersion(db)
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		fmt.Printf("Current migration version: %d\n", version)
	default:
		log.Fatalf("Unknown command: %s (use: up, down, version)", command)
	}
}

func migrateUp(db *sql.DB) error {
	// Create schema_migrations table if not exists
	if err := createMigrationsTable(db); err != nil {
		return err
	}

	// Read all migration files
	migrations := []struct {
		version int
		name    string
		sql     string
	}{
		{1, "create_users_table", getSQL("001_create_users_table.up.sql")},
		{2, "create_refresh_tokens_table", getSQL("002_create_refresh_tokens_table.up.sql")},
		{3, "create_books_table", getSQL("003_create_books_table.up.sql")},
		{4, "create_pages_table", getSQL("004_create_pages_table.up.sql")},
		{5, "create_review_tables", getSQL("005_create_review_tables.up.sql")},
		{6, "create_learning_tables", getSQL("006_create_learning_tables.up.sql")},
		{7, "create_ocr_tables", getSQL("007_create_ocr_tables.up.sql")},
		{8, "create_tts_stt_tables", getSQL("008_create_tts_stt_tables.up.sql")},
		{9, "create_dictionary_tables", getSQL("009_create_dictionary_tables.up.sql")},
		{10, "create_pattern_tables", getSQL("010_create_pattern_tables.up.sql")},
		{11, "create_teacher_mode_tables", getSQL("011_create_teacher_mode_tables.up.sql")},
	}

	// Also include subscription and stats tables
	subscriptionSQL := getSQL("001_create_subscription_tables.up.sql")
	statsSQL := getSQL("001_create_stats_tables.sql")

	// Apply subscription tables first if they exist
	if subscriptionSQL != "" {
		log.Println("Applying subscription tables migration...")
		if _, err := db.Exec(subscriptionSQL); err != nil {
			log.Printf("Warning: Subscription tables migration failed (may already exist): %v", err)
		}
	}

	// Apply stats tables
	if statsSQL != "" {
		log.Println("Applying stats tables migration...")
		if _, err := db.Exec(statsSQL); err != nil {
			log.Printf("Warning: Stats tables migration failed (may already exist): %v", err)
		}
	}

	// Apply each migration in order
	for _, m := range migrations {
		// Check if migration already applied
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", m.version).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration version %d: %w", m.version, err)
		}

		if count > 0 {
			log.Printf("⏭  Migration %03d_%s already applied, skipping", m.version, m.name)
			continue
		}

		// Apply migration
		log.Printf("⚙  Applying migration %03d_%s...", m.version, m.name)
		if _, err := db.Exec(m.sql); err != nil {
			return fmt.Errorf("failed to apply migration %03d_%s: %w", m.version, m.name, err)
		}

		// Record migration
		_, err = db.Exec("INSERT INTO schema_migrations (version, name) VALUES ($1, $2)", m.version, m.name)
		if err != nil {
			return fmt.Errorf("failed to record migration %d: %w", m.version, err)
		}

		log.Printf("✓ Migration %03d_%s applied successfully", m.version, m.name)
	}

	return nil
}

func migrateDown(db *sql.DB) error {
	log.Println("⚠  Warning: This will drop all tables and data!")
	log.Println("Rolling back all migrations...")

	migrations := []struct {
		version int
		name    string
		sql     string
	}{
		{11, "create_teacher_mode_tables", getSQL("011_create_teacher_mode_tables.down.sql")},
		{10, "create_pattern_tables", getSQL("010_create_pattern_tables.down.sql")},
		{9, "create_dictionary_tables", getSQL("009_create_dictionary_tables.down.sql")},
		{8, "create_tts_stt_tables", getSQL("008_create_tts_stt_tables.down.sql")},
		{7, "create_ocr_tables", getSQL("007_create_ocr_tables.down.sql")},
		{6, "create_learning_tables", getSQL("006_create_learning_tables.down.sql")},
		{5, "create_review_tables", getSQL("005_create_review_tables.down.sql")},
		{4, "create_pages_table", getSQL("004_create_pages_table.down.sql")},
		{3, "create_books_table", getSQL("003_create_books_table.down.sql")},
		{2, "create_refresh_tokens_table", getSQL("002_create_refresh_tokens_table.down.sql")},
		{1, "create_users_table", getSQL("001_create_users_table.down.sql")},
	}

	// Apply down migrations in reverse order
	for _, m := range migrations {
		log.Printf("⚙  Rolling back migration %03d_%s...", m.version, m.name)
		if _, err := db.Exec(m.sql); err != nil {
			log.Printf("Warning: Failed to rollback migration %03d_%s: %v", m.version, m.name, err)
		}

		// Delete migration record
		_, err := db.Exec("DELETE FROM schema_migrations WHERE version = $1", m.version)
		if err != nil {
			log.Printf("Warning: Failed to delete migration record %d: %v", m.version, err)
		}

		log.Printf("✓ Migration %03d_%s rolled back", m.version, m.name)
	}

	// Also rollback subscription and stats tables
	subscriptionDownSQL := getSQL("001_create_subscription_tables.down.sql")
	if subscriptionDownSQL != "" {
		log.Println("Rolling back subscription tables...")
		if _, err := db.Exec(subscriptionDownSQL); err != nil {
			log.Printf("Warning: Subscription tables rollback failed: %v", err)
		}
	}

	return nil
}

func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func getVersion(db *sql.DB) (int, error) {
	var version int
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	return version, err
}

func getSQL(filename string) string {
	content, err := os.ReadFile("migrations/" + filename)
	if err != nil {
		log.Printf("Warning: Failed to read %s: %v", filename, err)
		return ""
	}
	return string(content)
}
