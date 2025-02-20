package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/VDliveson/SurgeForms/go-producer/internal/dependencies"
	"github.com/VDliveson/SurgeForms/go-producer/routes"
	"github.com/VDliveson/SurgeForms/go-producer/utils"
	"github.com/gofiber/fiber/v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	di, err := dependencies.New(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	// Fiber App
	app := fiber.New()
	app.Use(utils.RequestLogger)
	routes.APIRoute(app, di)

	port := utils.GetEnv("PORT", "3000")
	address := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", address)

	go func() {
		if err := app.Listen(address); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal. Cleaning up...")

	di.Shutdown()
	log.Println("Server shut down gracefully.")
}
