package server

import (
	"log"
	"net/http"
	"os"
)

func Start() {

	port := "7540"

	if customPort := os.Getenv("TODO_PORT"); customPort != "" {
		port = customPort
	}

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер запущен на http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
