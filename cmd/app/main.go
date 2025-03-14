package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	calc "github.com/Tawxyn/goStockScraper/cmd/app/dcf"
	handlers "github.com/Tawxyn/goStockScraper/cmd/app/handlers"
	database "github.com/Tawxyn/goStockScraper/pkg"
	"github.com/joho/godotenv"
)

func main() {

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Ensure that the context is canceled when main returns

	//  Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatalf("DATABASE_URL was not found in the .env or is empty")
	}

	// Database import (pkg/database.go)
	pgInstance, err := database.InitDatabase(ctx, connString)
	if err != nil {
		log.Fatalf("Error initializing database post .env load: %v\n", err)
	}
	defer pgInstance.Close() // Close database after main exits

	// Create an instance of FinancialService
	financialService := calc.NewFinancialService(pgInstance)

	// Pass the database instance to the handlers
	handler := handlers.NewHandler(pgInstance, financialService)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Define HTTP routes
	http.HandleFunc("/user", handler.UserHandler)
	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/analyze", handler.AnalyzeHandler)
	http.HandleFunc("/CalculateWAAC", handler.CalculateWAAC)

	// State HTTP Server
	log.Fatal(http.ListenAndServe(":8080", nil))

}
