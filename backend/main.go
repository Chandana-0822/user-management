package main

import (
	"backend/database"
	"backend/routers"
	"log"
)

func main() {
	// Initialize database
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Start router
	e := routers.InitRouter()
	log.Println("Server started on port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
