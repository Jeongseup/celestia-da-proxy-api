package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	// Create table if not exists
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS blobs (
		hash TEXT PRIMARY KEY,
		height INTEGER
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	return db, nil
}

func InsertBlob(db *sql.DB, hash string, height int) error {
	insertQuery := `INSERT INTO blobs (hash, height) VALUES (?, ?)`
	_, err := db.Exec(insertQuery, hash, height)
	return err
}

func GetBlobHeight(db *sql.DB, hash string) (int, error) {
	var height int
	err := db.QueryRow("SELECT height FROM blobs WHERE hash = ?", hash).Scan(&height)
	if err != nil {
		return 0, err
	}
	return height, nil
}
