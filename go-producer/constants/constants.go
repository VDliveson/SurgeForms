package constants

import (
	"github.com/gofiber/fiber/v2"
)

const DatabaseName string = "producer"
const Exchange string = "daisy1"

type Response struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}

type QuestionBody struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type FormBody struct {
	Title       string         `json:"title" validate:"required"`
	Description string         `json:"description"`
	Questions   []QuestionBody `json:"questions"`
}
