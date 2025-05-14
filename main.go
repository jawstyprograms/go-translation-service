package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"expense-tracker/config"
	"expense-tracker/routes" // Import the routes package
)

func main() {
	conn, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer conn.Close(context.Background())

	fmt.Println("Database connection successful!")

	router := routes.SetupRoutes() // Get the router from routes.go

	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router)) // Start the server
}

// To run this code properly you have to set the Environment Variable so copy and past this
// into the terminal: export DATABASE_URL="postgres://postgres:@localhost:5432/expenses"   then run
