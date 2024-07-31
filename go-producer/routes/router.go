package routes

import (
	"net/http"

	"github.com/VDliveson/SurgeForms/go-producer/constants"
	"github.com/VDliveson/SurgeForms/go-producer/controllers"
	"github.com/gofiber/fiber/v2"
)

func APIRoute(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(constants.Response{
			Status:  http.StatusOK,
			Message: "Welcome to the producer API service.",
			Data: &fiber.Map{
				"api":         "Producer API Home route",
				"version":     "1.0",
				"description": "This is the home route of the Producer API",
			},
		})
	})
	app.Get("/api/forms/", controllers.HomeRoute)
	app.Get("/api/forms/get/:id", controllers.GetForm)
	app.Get("/api/forms/question/:id", controllers.GetQuestion)
	app.Post("/api/forms/create", controllers.CreateForm)
	// app.Post("/api/forms/response", controllers.AddResponse)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(http.StatusNotFound).JSON(constants.Response{
			Status:  http.StatusNotFound,
			Message: "Route not found",
			Data:    nil,
		})
	})
}
