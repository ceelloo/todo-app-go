package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func New() *sql.DB {
	db, err := sql.Open("sqlite3", "./local.db")
	if err != nil {
		log.Fatal("Failed to open db", err)
	}

	var version string
	db.QueryRow("SELECT sqlite_version()").Scan(&version)
	
	log.Printf("SQLite Version: %s", version)
	log.Print("Connected to database")

	return db
}

func InitializeDatabase(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS task (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		completed BOOLEAN NOT NULL,
		body TEXT NOT NULL
	);
	`)
	if err != nil {
		log.Fatal("Failed to create table", err)
	}
}