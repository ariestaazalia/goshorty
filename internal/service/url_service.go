package service

import (
	"math/rand"

	"github.com/ariestaazalia/goshorty/internal/repository"
)

type URLService struct {
	Repository repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return URLService{Repository: repo}
}

func (svc *URLService) ShortenURL(longURL string) string {
	code := generateShortCode(10)

	svc.Repository.SaveURL(code, longURL)

	return code
}

func (svc *URLService) GetURL(code string) (string, bool) {
	return svc.Repository.GetURL(code)
}

func generateShortCode(length int) string {
	const charset = "abcdefghjklmnopqrstuvwxyz1234567890"
	code := make([]byte, length)

	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}

	return string(code)
}

func (svc *URLService) GetOriginalURL(shortCode string) (string, error) {
	originalURL, err := svc.Repository.GetURLByShortCode(shortCode)

	if err != nil {
		return "", err
	}

	return originalURL, nil
}