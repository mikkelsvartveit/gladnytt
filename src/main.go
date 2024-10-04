package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

type PageData struct {
	NextPage int
	Articles []Article
}

func main() {
	// Load the .env file
	godotenv.Load()

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

	// Load HTML template
	indexTmpl := template.Must(template.ParseFiles("templates/index.html"))

	// Serve front page
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")

		var page int
		if pageStr == "" {
			page = 1
		} else {
			var err error
			page, err = strconv.Atoi(pageStr)

			if err != nil {
				page = 1
			}
		}

		articles := listArticles(1, 10*page)

		pageData := PageData{
			NextPage: page + 1,
			Articles: articles,
		}

		err := indexTmpl.Execute(w, pageData)
		if err != nil {
			fmt.Println("Error executing template:", err)
		}
	})

	// Serve articles fragment for htmx lazy loading
	r.Get("/articles/{page}", func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(chi.URLParam(r, "page"))
		if err != nil {
			fmt.Println("Error parsing page number:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		articles := listArticles(page, 10)

		if len(articles) == 0 {
			w.WriteHeader(http.StatusOK)
			return
		}

		pageData := PageData{
			NextPage: page + 1,
			Articles: articles,
		}

		err = indexTmpl.ExecuteTemplate(w, "articles", pageData)
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
