package repository

import (
	"context"

	"github.com/sammyqtran/url-shortener/internal/models"
)

type URLRepository interface {
	// Create stores a new URL mapping
	Create(ctx context.Context, url *models.URL) error

	// GetByShortCode retrieves URL by short code
	GetByShortCode(ctx context.Context, shortCode string) (*models.URL, error)

	// GetByID retrieves URL by ID
	GetByID(ctx context.Context, id int64) (*models.URL, error)

	// Update modifies an existing URL
	Update(ctx context.Context, url *models.URL) error

	// Delete removes a URL by short code
	Delete(ctx context.Context, shortCode string) error

	// IncrementClickCount increments the click counter
	IncrementClickCount(ctx context.Context, shortCode string) error

	// GetStats returns URL statistics
	GetStats(ctx context.Context, shortCode string) (*models.URL, error)

	// ListURLs returns paginated list of URLs
	ListURLs(ctx context.Context, limit, offset int) ([]*models.URL, error)

	// IsShortCodeExists checks if short code already exists
	IsShortCodeExists(ctx context.Context, shortCode string) (bool, error)
}
