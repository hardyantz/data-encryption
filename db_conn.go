package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	MaxIdleConnection = 10
	MaxOpenConnection = 10
	connectionPattern = "host=%s user=%s dbname=%s sslmode=disable"
)

// ReadCronPostgresDB function for creating database connection for read-access
func ConnectDB() *sql.DB {
	return CreateDBConnection(fmt.Sprintf(connectionPattern,
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_NAME")))
}

// CreateDBConnection function for creating database connection
func CreateDBConnection(descriptor string) *sql.DB {
	db, err := sql.Open("postgres", descriptor)
	if err != nil {
		defer db.Close()
		return db
	}

	db.SetMaxIdleConns(MaxIdleConnection)
	db.SetMaxOpenConns(MaxOpenConnection)

	return db
}

// CloseDb function for closing database connection
func CloseDb(db *sql.DB) {
	if db != nil {
		db.Close()
		db = nil
	}
}