package database

import (
	"context"
	"fmt"
	"golang-auth/config"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	var err error

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.GetEnv("DB_USER", "postgres"),
		config.GetEnv("DB_PASSWORD", "password"),
		config.GetEnv("DB_HOST", "localhost"),
		config.GetEnv("DB_PORT", "5432"),
		config.GetEnv("DB_NAME", "dbname"),
	)

	log.Printf("Connecting to database with DSN: %s", dsn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DB, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to create connection pool: %v", err)
	}

	// Bağlantıyı test et
	if err := DB.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to the database successfully")
}

func CloseDB() {
	DB.Close()
	log.Println("Database connection closed")
}
