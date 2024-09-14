package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type ArticleTable struct {
	ID          int          `db:"id"`
	Title       string       `db:"title"`
	Description string       `db:"description"`
	ArticleUrl  string       `db:"articleUrl"`
	ImageUrl    string       `db:"imageUrl"`
	IsShown     sql.NullBool `db:"isShown"`
	Datetime    int64        `db:"datetime"`
}

type Article struct {
	Title       string
	Description string
	ArticleUrl  string
	ImageUrl    string
	Time        time.Time
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
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			articleUrl TEXT UNIQUE NOT NULL,
			imageUrl TEXT NOT NULL,
			datetime INTEGER NOT NULL,
			isShown BOOLEAN
		);
	`)

	return err
}

func addArticleIfNotExists(title string, description string, articleUrl string, imageUrl string, pubDate string) {
	// Check if article already exists in database
	var articleId int
	err := db.Get(&articleId, `
		SELECT id FROM articles
		WHERE articleUrl = ?
	`, articleUrl)

	if err == nil {
		// Article already exists, skip
		return
	}

	// Convert pubDate to Unix timestamp
	time, err := time.Parse(time.RFC1123, pubDate)
	if err != nil {
		fmt.Println("Error parsing pubDate:", err)
		return
	}
	datetime := time.Unix()

	// TODO: Run article through LLM filter to get isShown

	_, dbInsertErr := db.Exec(`
		INSERT INTO articles (title, description, articleUrl, imageUrl, datetime)
		VALUES (?, ?, ?, ?, ?)
	`, title, description, articleUrl, imageUrl, datetime)

	if dbInsertErr != nil {
		fmt.Println("Error adding article:", dbInsertErr)
	}
}

func getArticles(page int, limit int) []Article {
	var rawArticles []ArticleTable

	err := db.Select(&rawArticles, `
		SELECT id, title, description, articleUrl, imageUrl, datetime, isShown
		FROM articles
		ORDER BY datetime DESC
		LIMIT ? OFFSET ?
	`, limit, (page-1)*limit)

	if err != nil {
		fmt.Println("Error getting articles:", err)
		return nil
	}

	var articles []Article
	for _, rawArticle := range rawArticles {
		article := Article{
			Title:       rawArticle.Title,
			Description: rawArticle.Description,
			ArticleUrl:  rawArticle.ArticleUrl,
			ImageUrl:    rawArticle.ImageUrl,
			Time:        time.Unix(rawArticle.Datetime, 0),
		}

		articles = append(articles, article)

	}

	return articles
}
