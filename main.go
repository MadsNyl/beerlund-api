package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"api.beerlund.com/m/db"
	"api.beerlund.com/m/handlers"
)

func main() {
	_ = godotenv.Load()

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	store := db.NewPostgresStore()

	if err := store.Init(dbUrl); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer store.Close()

	handler := &handlers.Handler{Store: store}

	http.Handle("/events", http.HandlerFunc(handler.ListEvents))

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}