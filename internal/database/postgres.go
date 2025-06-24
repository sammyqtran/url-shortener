package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// NewPostgresConnection creates a new PostgreSQL database connection
func NewPostgresConnection(config Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DatabaseName, config.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.MaxLifetime)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// RunMigrations applies database migrations
func RunMigrations(db *sqlx.DB) error {
	// Simple migration runner - in production, use a proper migration tool like golang-migrate
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS urls (
            id BIGSERIAL PRIMARY KEY,
            user_id VARCHAR(50) NOT NULL,
            short_code VARCHAR(10) UNIQUE NOT NULL,
            original_url TEXT NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            click_count BIGINT DEFAULT 0,
            expires_at TIMESTAMP WITH TIME ZONE
        )`,
		`CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls (short_code)`,
		`CREATE INDEX IF NOT EXISTS idx_urls_created_at ON urls (created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_urls_user_id ON urls (user_id)`,
		`CREATE OR REPLACE FUNCTION update_updated_at_column()
         RETURNS TRIGGER AS $$
         BEGIN
             NEW.updated_at = CURRENT_TIMESTAMP;
             RETURN NEW;
         END;
         $$ language 'plpgsql'`,
		`DROP TRIGGER IF EXISTS update_urls_updated_at ON urls`,
		`CREATE TRIGGER update_urls_updated_at 
         BEFORE UPDATE ON urls 
         FOR EACH ROW 
         EXECUTE FUNCTION update_updated_at_column()`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to run migration: %w", err)
		}
	}

	return nil
}
