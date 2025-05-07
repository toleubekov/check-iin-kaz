package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/toleubekov/kaspiCheckIIN/internal/api"
	"github.com/toleubekov/kaspiCheckIIN/internal/repository"
	"github.com/toleubekov/kaspiCheckIIN/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "qwerty")
	dbName := getEnv("DB_NAME", "postgres")
	dbHost := getEnv("DB_HOST", "postgres")
	dbPort := getEnv("DB_PORT", "5432")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dbConnectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	db, err := repository.InitDB(dbConnectionString)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	personRepo := repository.NewPersonRepository(db)
	iinService := service.NewIINService()
	handler := api.NewHandler(iinService, personRepo)

	router := api.SetupRouter(handler)

	port := getEnv("SERVER_PORT", "8080")

	log.Printf("Server is running on port %s...", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
