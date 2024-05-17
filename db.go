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

	// Create namespace table
	createNamespaceTableQuery := `
	CREATE TABLE IF NOT EXISTS namespaces (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		namespace_key TEXT,
		index_number INTEGER,
		hash TEXT,
		height INTEGER,
		UNIQUE(namespace_key, index_number)
	);`
	_, err = db.Exec(createNamespaceTableQuery)
	if err != nil {
		log.Fatalf("failed to create namespace table: %v", err)
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

func InsertNamespace(db *sql.DB, namespaceKey string, hash string, height int) (int, error) {
	// Get the maximum index for the given namespace_key
	var maxIndex int
	err := db.QueryRow("SELECT COALESCE(MAX(index_number), 0) FROM namespaces WHERE namespace_key = ?", namespaceKey).Scan(&maxIndex)
	if err != nil {
		return 0, err
	}

	// Increment the max index by 1 for the new entry
	newIndex := maxIndex + 1

	insertQuery := `INSERT INTO namespaces (namespace_key, index_number, hash, height) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(insertQuery, namespaceKey, newIndex, hash, height)
	return newIndex, err
}

func GetNamespaceData(db *sql.DB, namespaceKey string, index int) (string, int, error) {
	var hash string
	var height int
	err := db.QueryRow("SELECT hash, height FROM namespaces WHERE namespace_key = ? AND index_number = ?", namespaceKey, index).Scan(&hash, &height)
	if err != nil {
		return "", 0, err
	}
	return hash, height, nil
}
