package constants

import (
	"github.com/gofiber/fiber/v2"
)

const DatabaseName string = "producer"
const Exchange string = "daisy1"

type Response struct {
	Success bool       `json:"success"`
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

type AnswerBody struct {
	Question string `json:"question"`
	Text     string `json:"text"`
}

type ResponseBody struct {
	Form     string                 `json:"form"`
	User     string                 `json:"user"`
	Answers  []AnswerBody           `json:"answers"`
	Metadata map[string]interface{} `json: "metadata"`
}
