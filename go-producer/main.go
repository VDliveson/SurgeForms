package main

import (
	"fmt"
	"log"

	"github.com/VDliveson/SurgeForms/go-producer/routes"
	"github.com/VDliveson/SurgeForms/go-producer/utils"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Connect to RabbitMQ
	if err := utils.ConnectQueue(); err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}

	// Connect to MongoDB
	if err := utils.ConnectDB(); err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Apply middleware
	app.Use(utils.RequestLogger)

	// Set up routes
	routes.APIRoute(app)

	// Get port from environment or use default
	port := utils.GetEnv("PORT", "3000")

	// Start the server
	address := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", address)
	if err := app.Listen(address); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
