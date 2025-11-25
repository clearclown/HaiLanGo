package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
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

	// rootユーザーの作成
	if err := createRootUser(db); err != nil {
		log.Fatalf("Failed to create root user: %v", err)
	}

	log.Println("✓ Seed completed successfully")
}

func createRootUser(db *sql.DB) error {
	// rootユーザーが既に存在するか確認
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", "root@hailango.dev").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("⏭  Root user already exists, skipping")
		return nil
	}

	// パスワード "passwd" をbcryptでハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("passwd"), 10)
	if err != nil {
		return err
	}

	// rootユーザーを作成
	userID := uuid.New().String()
	_, err = db.Exec(`
		INSERT INTO users (id, email, password_hash, display_name, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`, userID, "root@hailango.dev", string(hashedPassword), "Root User", true)

	if err != nil {
		return err
	}

	log.Println("✓ Root user created successfully")
	log.Println("  Email: root@hailango.dev")
	log.Println("  Password: passwd")

	return nil
}
