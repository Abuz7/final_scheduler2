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

	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using system environment variables: %v", err)
	}
	err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	api.Init()
	fmt.Println("Database initialized successfully.")
	server.Start()
}
