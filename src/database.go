package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type ArticleModel struct {
	ID          int    `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Timestamp   int64  `db:"timestamp"`
	ArticleUrl  string `db:"articleUrl"`
	ImageUrl    string `db:"imageUrl"`
	Sentiment   string `db:"sentiment"`
}

type Article struct {
	Title       string
	Description string
	Time        time.Time
	ArticleUrl  string
	ImageUrl    string
	Sentiment   string
}

var db *sqlx.DB

func initializeDatabase() error {
	dbPath := "./data/database.db"

	// Create the directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	db, err = sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	err = createTablesIfNotExist()
	if err != nil {
		return fmt.Errorf("error creating tables: %w", err)
	}

	return nil
}

func createTablesIfNotExist() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS articles (
			id INTEGER PRIMARY KEY,
			title TEXT,
			description TEXT,
			timestamp INTEGER,
			articleUrl TEXT,
			imageUrl TEXT,
			sentiment TEXT
		);
	`)

	return err
}

func articleExists(articleUrl string) bool {
	var articleId int

	err := db.Get(&articleId, `
		SELECT id FROM articles
		WHERE articleUrl = ?
	`, articleUrl)

	return err == nil
}

func insertArticle(article Article) {
	_, dbInsertErr := db.Exec(`
		INSERT INTO articles (title, description, articleUrl, imageUrl, timestamp, sentiment)
		VALUES (?, ?, ?, ?, ?, ?)
	`, article.Title, article.Description, article.ArticleUrl, article.ImageUrl, article.Time.Unix(), article.Sentiment)

	if dbInsertErr != nil {
		fmt.Println("Error adding article:", dbInsertErr)
	}
}

func listArticles(page int, limit int) []Article {
	var articleRows []ArticleModel

	err := db.Select(&articleRows, `
		SELECT id, title, description, articleUrl, imageUrl, timestamp, sentiment
		FROM articles
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`, limit, (page-1)*limit)

	if err != nil {
		fmt.Println("Error getting articles:", err)
		return nil
	}

	var articles []Article
	for _, row := range articleRows {
		article := Article{
			Title:       row.Title,
			Description: row.Description,
			ArticleUrl:  row.ArticleUrl,
			ImageUrl:    row.ImageUrl,
			Time:        time.Unix(row.Timestamp, 0),
			Sentiment:   row.Sentiment,
		}

		articles = append(articles, article)
	}

	return articles
}
