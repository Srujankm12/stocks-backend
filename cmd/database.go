package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// NewConnection initializes a new database connection
func NewConnection() *sql.DB {
	conn, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("unable to create database: %v", err)
	}
	return conn
}
