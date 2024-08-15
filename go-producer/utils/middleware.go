package utils

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestLogger(c *fiber.Ctx) error {
	// Start a timer
	start := time.Now()

	// Process the request
	err := c.Next()

	// Calculate the duration
	duration := time.Since(start)

	// Log the response status and time taken
	log.Printf("Response: status=%d, method=%s, url=%s, ip=%s, duration=%v", c.Response().StatusCode(), c.Method(), c.OriginalURL(), c.IP(), duration)

	return err
}
