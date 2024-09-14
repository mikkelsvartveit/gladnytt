package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Initialize database
	dbErr := initializeDatabase()
	if dbErr != nil {
		fmt.Println("Error initializing database:", dbErr)
		return
	}

	// Run Goroutine to fetch data periodically
	go runPeriodically(time.Minute*5, fetchData)

	// Serve index.html
	indexTmpl := template.Must(template.ParseFiles("templates/index.html"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		articles := getArticles(1, 10)
		err := indexTmpl.Execute(w, articles)
		if err != nil {
			fmt.Println("Error executing template:", err)
		}
	})

	// Serve all files from the "static" directory
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("static"))))

	httpErr := http.ListenAndServe(":8080", r)
	if httpErr != nil {
		fmt.Printf("Error starting server: %s\n", httpErr)
	}
}

func runPeriodically(interval time.Duration, f func()) {
	// Run the function immediately on startup
	f()

	// Run the function on the specified interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		f()
	}
}
