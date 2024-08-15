package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AnswerSchema struct {
	Id       primitive.ObjectID `json:"_id,omitempty"`
	Question primitive.ObjectID `json:"question,omitempty" validate:"required"`
	Response primitive.ObjectID `json:"response,omitempty" validate:"required"`
	Text     string             `json:"text,omitempty"`
}
