package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/sammyqtran/url-shortener/internal/models"
	"github.com/sammyqtran/url-shortener/internal/repository"
)

type postgresURLRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewPostgresURLRepository creates a new PostgreSQL URL repository
func NewPostgresURLRepository(db *sqlx.DB, logger *zap.Logger) repository.URLRepository {
	return &postgresURLRepository{
		db:     db,
		logger: logger,
	}
}

func (r *postgresURLRepository) Create(ctx context.Context, url *models.URL) error {
	query := `
        INSERT INTO urls (user_id, short_code, original_url, expires_at) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id, created_at, updated_at, click_count
    `

	err := r.db.QueryRowxContext(ctx, query, url.UserID, url.ShortCode, url.OriginalURL, url.ExpiresAt).
		Scan(&url.ID, &url.CreatedAt, &url.UpdatedAt, &url.ClickCount)

	if err != nil {
		r.logger.Error("Failed to create row in database", zap.Error(err))
		return fmt.Errorf("failed to create URL: %w", err)
	}

	return nil
}

func (r *postgresURLRepository) GetByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	var url models.URL
	query := `
        SELECT id, short_code, original_url, created_at, updated_at, click_count, expires_at
        FROM urls 
        WHERE short_code = $1 AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
    `

	err := r.db.GetContext(ctx, &url, query, shortCode)
	if err != nil {
		r.logger.Error("Error retrieving shortCode", zap.String("shortCode", shortCode), zap.Error(err))
		if err == sql.ErrNoRows {
			return nil, repository.ErrURLNotFound
		}
		return nil, fmt.Errorf("failed to get URL by short code: %w", err)
	}

	return &url, nil
}

func (r *postgresURLRepository) GetByID(ctx context.Context, id int64) (*models.URL, error) {
	var url models.URL
	query := `
        SELECT id, short_code, original_url, created_at, updated_at, click_count, expires_at
        FROM urls 
        WHERE id = $1
    `

	err := r.db.GetContext(ctx, &url, query, id)
	if err != nil {
		r.logger.Error("Error retrieving rows by ID", zap.Int64("id", id), zap.Error(err))
		if err == sql.ErrNoRows {
			return nil, repository.ErrURLNotFound
		}
		return nil, fmt.Errorf("failed to get URL by ID: %w", err)
	}

	return &url, nil
}

func (r *postgresURLRepository) Update(ctx context.Context, url *models.URL) error {
	query := `
        UPDATE urls 
        SET original_url = $1, expires_at = $2, updated_at = CURRENT_TIMESTAMP
        WHERE short_code = $3
    `

	result, err := r.db.ExecContext(ctx, query, url.OriginalURL, url.ExpiresAt, url.ShortCode)
	if err != nil {
		r.logger.Error("Error updating URL", zap.Error(err))
		return fmt.Errorf("failed to update URL: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Error getting affected rows", zap.Error(err))
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrURLNotFound
	}

	return nil
}

func (r *postgresURLRepository) Delete(ctx context.Context, shortCode string) error {
	query := `DELETE FROM urls WHERE short_code = $1`

	result, err := r.db.ExecContext(ctx, query, shortCode)
	if err != nil {
		r.logger.Error("Error deleting URL", zap.Error(err))
		return fmt.Errorf("failed to delete URL: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Error getting affected rows", zap.Error(err))
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrURLNotFound
	}

	return nil
}

func (r *postgresURLRepository) IncrementClickCount(ctx context.Context, shortCode string) error {
	query := `
        UPDATE urls 
        SET click_count = click_count + 1, updated_at = CURRENT_TIMESTAMP
        WHERE short_code = $1
    `

	result, err := r.db.ExecContext(ctx, query, shortCode)
	if err != nil {
		r.logger.Error("Error incrementing click count", zap.Error(err))
		return fmt.Errorf("failed to increment click count: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Error getting affected rows", zap.Error(err))
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrURLNotFound
	}

	return nil
}

func (r *postgresURLRepository) GetStats(ctx context.Context, shortCode string) (*models.URL, error) {
	return r.GetByShortCode(ctx, shortCode)
}

func (r *postgresURLRepository) ListURLs(ctx context.Context, limit, offset int) ([]*models.URL, error) {
	var urls []*models.URL
	query := `
        SELECT id, short_code, original_url, created_at, updated_at, click_count, expires_at
        FROM urls 
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	err := r.db.SelectContext(ctx, &urls, query, limit, offset)
	if err != nil {
		r.logger.Error("Error getting rows", zap.Error(err))
		return nil, fmt.Errorf("failed to list URLs: %w", err)
	}

	return urls, nil
}

func (r *postgresURLRepository) IsShortCodeExists(ctx context.Context, shortCode string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`

	err := r.db.GetContext(ctx, &exists, query, shortCode)
	if err != nil {
		r.logger.Error("Error checking existence", zap.Error(err))
		return false, fmt.Errorf("failed to check if short code exists: %w", err)
	}

	return exists, nil
}
