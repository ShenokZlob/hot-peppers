package main

import (
	"fmt"
	"hotpepper/internal/ctrl"
	"hotpepper/internal/ctrl/middleware"
	"hotpepper/internal/repo"
	"log"
	"net/http"
)

func main() {
	r := repo.NewPeppers()

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Привет от Go-сервера!")
	})
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("data/static"))))
	ctrl.NewPeppersHandler(mux, r)

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: middleware.CorsMiddleware(mux),
	}

	log.Println("Server podnyalsya")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server upal: %v", err)
	}
}
