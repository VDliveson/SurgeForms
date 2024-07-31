package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type QuestionSchema struct {
	Id   primitive.ObjectID `bson:"_id,omitempty"`
	Form primitive.ObjectID `bson:"form,omitempty" validate:"required"`
	Text string             `bson:"text,omitempty"`
	Type string             `bson:"type,omitempty"`
}
