package repository

import (
	"errors"
	"time"
)

var (
	ErrURLNotFound = errors.New("URL Not Found")
	ErrURLExpired  = errors.New("URL Has Expired")
)

type URLRepository interface {
	SaveURL(code, longURL string)
	GetURL(code string) (string, bool)
	GetURLByShortCode(shortCode string) (string, error)
}

type URLData struct {
	originalURL string
	CreatedAt	time.Time
}

type InMemoryURLRepository struct {
	URLs map[string]URLData
}

func NewURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{
		URLs: make(map[string]URLData),
	}
}

func (repo *InMemoryURLRepository) SaveURL(code, longURL string) {
	repo.URLs[code] = URLData{
		originalURL: longURL,
		CreatedAt: time.Now(),
	}
}

func (repo *InMemoryURLRepository) GetURL(code string) (string, bool) {
	data, exists := repo.URLs[code]

	if !exists {
		return "", false
	}
	
	return data.originalURL, true
}

func (repo *InMemoryURLRepository) GetURLByShortCode(shortCode string) (string, error) {
	data, exists := repo.URLs[shortCode]

	if !exists {
		return "", ErrURLNotFound
	}

	if time.Since(data.CreatedAt) > 7*24*time.Hour {
		delete(repo.URLs, shortCode)
		return "", ErrURLExpired
	}

	return data.originalURL, nil
}