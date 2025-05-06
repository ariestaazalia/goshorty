package main

import (
	"log"
	"net/http"

	"github.com/ariestaazalia/goshorty/internal/handler"
	"github.com/ariestaazalia/goshorty/internal/repository"
	"github.com/ariestaazalia/goshorty/internal/service"
)

func main() {
	repo := repository.NewURLRepository()
	service := service.NewURLService(repo)
	handler := handler.NewURLHandler(service)

	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/shorten", handler.Shorten)
	http.HandleFunc("/r/", handler.Redirect)

	log.Println("Server started on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}