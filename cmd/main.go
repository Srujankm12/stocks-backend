package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Srujankm12/SRstocks/pkg/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("unable to load env: %v", err)
	}

	// Create a new database connection
	conn := NewConnection()
	defer conn.Close() // FIX: Use `conn.Close()` instead of `conn.DB.Close()`

	// Create tables if they don't exist
	query := database.NewQuery(conn)
	if err := query.CreateTables(); err != nil {
		log.Fatalf("Unable to create database tables: %v", err)
	}

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}
	server := &http.Server{
		Addr:    ":" + port, // FIX: Ensure valid port format
		Handler: registerRouter(conn),
	}

	log.Printf("Server is running at port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Unable to start the server: %v", err)
	}
}
