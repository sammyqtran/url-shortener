package repository

import "errors"

var (
	ErrURLNotFound     = errors.New("URL not found")
	ErrShortCodeExists = errors.New("short code already exists")
	ErrInvalidURL      = errors.New("invalid URL format")
	ErrExpiredURL      = errors.New("URL has expired")
)
