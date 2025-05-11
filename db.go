package main

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "./local.db")
	if err != nil {
		log.Fatal("Failed to open db", err)
	}

	createTable := `
		CREATE TABLE IF NOT EXISTS task (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			completed BOOLEAN,
			body TEXT
		);
	`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create table", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}
}
