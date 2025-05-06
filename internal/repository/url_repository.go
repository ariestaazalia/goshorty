package repository

import "errors"

type URLRepository interface {
	SaveURL(code, longURL string)
	GetURL(code string) (string, bool)
	GetURLByShortCode(shortCode string) (string, error)
}

type InMemoryURLRepository struct {
	URLs map[string]string
}

func NewURLRepository() *InMemoryURLRepository {
	return &InMemoryURLRepository{URLs: make(map[string]string)}
}

func (repo *InMemoryURLRepository) SaveURL(code, longURL string) {
	repo.URLs[code] = longURL
}

func (repo *InMemoryURLRepository) GetURL(code string) (string, bool) {
	url, exists := repo.URLs[code]
	
	return url, exists
}

func (repo *InMemoryURLRepository) GetURLByShortCode(shortCode string) (string, error) {
	originalURL, exists := repo.URLs[shortCode]

	if !exists {
		return "", errors.New("URL Not Found")
	}

	return originalURL, nil
}