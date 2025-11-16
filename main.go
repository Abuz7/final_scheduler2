package main

import (
	"final/pkg/api"
	"final/pkg/db"
	"final/pkg/server"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	dbFile := "scheduler.db"
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		dbFile = envDBFile
	}
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using system environment variables: %v", err)
	}

	db, err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer db.Close() // Закрываем подключение при выходе из main

	api.Init()
	fmt.Println("Database initialized successfully.")
	server.Start()
}
