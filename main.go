package main

import (
	"log"
	"net/http"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/joho/godotenv"

	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"

	"api.beerlund.com/m/db"
	"api.beerlund.com/m/handlers"
	"api.beerlund.com/m/middleware"
	"api.beerlund.com/m/logger"
)

func main() {
	_ = godotenv.Load()

	logEndpoint := os.Getenv("LOG_ENDPOINT")
	if logEndpoint == "" {
		log.Fatal("LOG_ENDPOINT environment variable is not set")
	}

	// Initialize the logger
	logger.InitLogger(
		logEndpoint,
		100,
		"beerlund-api",
	)

	logger.Info("Starting BeerLund API server...", nil)

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	clerkApiKey := os.Getenv("CLERK_API_KEY")
	if clerkApiKey == "" {
		log.Fatal("CLERK_API_KEY environment variable is not set")
	}

	clerk.SetKey(clerkApiKey)

	store := db.NewPostgresStore()

	if err := store.Init(dbUrl); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer store.Close()

	handler := &handlers.Handler{Store: store}

	mux := http.NewServeMux()

	mux.Handle("/events", http.HandlerFunc(handler.ListEvents))
	mux.Handle("/events/", http.HandlerFunc(handler.GetEvent))

    joinHandler := http.HandlerFunc(handler.JoinEvent)
    protectedJoin := clerkhttp.WithHeaderAuthorization()(joinHandler)
    mux.Handle("/participate", protectedJoin)

	unassignHandler := http.HandlerFunc(handler.LeaveEvent)
	protectedUnassign := clerkhttp.WithHeaderAuthorization()(unassignHandler)
	mux.Handle("/leave", protectedUnassign)

	wrapped := middleware.CorsMiddleware(mux)

	log.Fatal(http.ListenAndServe(":8000", wrapped))
}