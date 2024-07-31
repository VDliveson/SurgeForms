package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AnswerSchema struct {
	Id       primitive.ObjectID `json:"_id,omitempty"`
	Question string             `json:"question,omitempty" validate:"required"`
	Response string             `json:"response,omitempty" validate:"required"`
	Text     string             `json:"text,omitempty"`
}
