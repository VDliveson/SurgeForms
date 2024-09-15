package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/VDliveson/SurgeForms/go-producer/routes"
	"github.com/VDliveson/SurgeForms/go-producer/utils"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	var wg sync.WaitGroup
	errChan := make(chan error, 2) // Buffer size of 2 to handle both errors

	// Connect to RabbitMQ in a separate goroutine
	wg.Add(1)
	go func() {
		log.Println("Connecting to RabbitMQ ...")
		defer wg.Done()
		if err := utils.ConnectQueue(); err != nil {
			errChan <- fmt.Errorf("error connecting to RabbitMQ: %v", err)
		}
	}()

	// Connect to MongoDB in a separate goroutine
	wg.Add(1)
	go func() {
		log.Println("Connecting to MongoDB...")
		defer wg.Done()
		if err := utils.ConnectDB(); err != nil {
			errChan <- fmt.Errorf("error connecting to MongoDB: %v", err)
		}
	}()

	// Wait for both goroutines to complete
	wg.Wait()
	close(errChan)

	// Check if there were any errors
	for err := range errChan {
		if err != nil {
			log.Fatal(err)
		}
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
