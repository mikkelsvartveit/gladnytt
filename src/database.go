package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func initializeDatabase() error {
	// Specify the database file path
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
			articleUrl TEXT UNIQUE NOT NULL,
			imageUrl TEXT NOT NULL,
			datetime INTEGER NOT NULL,
			title TEXT,
			summary TEXT
		);
	`)

	return err
}

func addArticleIfNotExists(articleUrl string, imageUrl string, pubDate string) {
	// Check if article already exists in database
	var articleId int
	err := db.Get(&articleId, `
		SELECT id FROM articles
		WHERE articleUrl = ?
	`, articleUrl)

	if err == nil {
		// Article already exists, skip
		fmt.Println("Article already exists:", articleUrl)
		return
	}

	// Convert pubDate to Unix timestamp
	datetime, err := time.Parse(time.RFC1123, pubDate)
	if err != nil {
		fmt.Println("Error parsing pubDate:", err)
		return
	}

	_, dbInsertErr := db.Exec(`
		INSERT INTO articles (articleUrl, imageUrl, datetime)
		VALUES (?, ?, ?)
	`, articleUrl, imageUrl, datetime)

	if dbInsertErr != nil {
		fmt.Println("Error adding article:", dbInsertErr)
	}
}
