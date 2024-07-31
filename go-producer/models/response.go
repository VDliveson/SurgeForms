package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ResponseSchema struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Form        primitive.ObjectID `bson:"form,omitempty" validate:"required"`
	User        primitive.ObjectID `bson:"user,omitempty" validate:"required"`
	SubmittedAt primitive.DateTime `bson:"submittedAt,omitempty"` // Default needs to be date.now()
}
