package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"ProjectNotification/internal/config"
	_ "github.com/lib/pq"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func NewDB() *sql.DB {
	config.OverloadDotEnv()

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "notification")
	password := getEnv("DB_PASSWORD", "notification")
	dbname := getEnv("DB_NAME", "notification_db")
	sslmode := getEnv("DB_SSLMODE", "disable")

	log.Printf("db config: host=%s port=%s user=%s dbname=%s sslmode=%s", host, port, user, dbname, sslmode)

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("db open error:", err)
	}

	var pingErr error
	for i := 0; i < 15; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		pingErr = db.PingContext(ctx)
		cancel()
		if pingErr == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if pingErr != nil {
		log.Fatal("db ping error:", pingErr)
	}

	if err := ensureNotificationTable(db); err != nil {
		log.Fatal("db init error:", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("✅  connected to postgres")

	return db
}

func ensureNotificationTable(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS notification (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  type TEXT NOT NULL,
  payload JSONB NOT NULL,
  status TEXT NOT NULL,
  created_at BIGINT NOT NULL
);
`)
	return err
}
