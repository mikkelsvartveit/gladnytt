package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	indexTmpl := template.Must(template.ParseFiles("templates/index.html"))

	// Serve index.html at the root URL
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		indexTmpl.Execute(w, nil)
	})

	// Serve static files from the "static" directory
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("static"))))

	fmt.Println("Starting server on http://localhost:8080")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
