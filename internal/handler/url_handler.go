package handler

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/ariestaazalia/goshorty/internal/repository"
	"github.com/ariestaazalia/goshorty/internal/service"
)

type URLHandler struct {
	Service service.URLService
}

func NewURLHandler(svc service.URLService) *URLHandler {
	return &URLHandler{Service: svc}
}

func (h *URLHandler) Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/index.html")
	
	if err != nil {
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	shortCode 	:= h.Service.ShortenURL(originalURL)

	// Get Current Base URL
	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	shortened 	:= scheme + "://" + host + "/r/" + shortCode

	tmpl, err := template.ParseFiles("web/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	
	tmpl.Execute(w, map[string]interface{}{
		"Shortened": shortened,
	})
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/r/"):]

	originalURL, err := h.Service.GetOriginalURL(code)

	if err != nil {
		if errors.Is(err, repository.ErrURLExpired) {
			tmpl, _ := template.ParseFiles("web/error.html")
			w.WriteHeader(http.StatusGone)

			tmpl.Execute(w, map[string]string {
				"Title": "URL Expired",
				"Message": "This shortened URL is no longer valid. It may have expired or been removed.",
			})
			return
		}
		if errors.Is(err, repository.ErrURLNotFound) {
			tmpl, _ := template.ParseFiles("web/error.html")
			w.WriteHeader(http.StatusNotFound)

			tmpl.Execute(w, map[string]string {
				"Title": "URL Not Found",
				"Message": "We couldn’t find the URL you’re looking for",
			})

			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}