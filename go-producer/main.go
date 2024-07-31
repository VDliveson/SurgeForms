package main

import (
	"log"

	"github.com/VDliveson/SurgeForms/go-producer/routes"
	"github.com/VDliveson/SurgeForms/go-producer/utils"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	err := utils.ConnectQueue()
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
		return
	}

	err = utils.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
		return
	}

	routes.APIRoute(app) //add this
	app.Listen(":3000")

}
