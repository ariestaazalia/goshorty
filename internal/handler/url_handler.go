package handler

import (
	"html/template"
	"net/http"

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
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}