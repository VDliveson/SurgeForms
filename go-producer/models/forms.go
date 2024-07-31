package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FormSchema struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title,omitempty" validate:"required"`
	Description string             `bson:"description,omitempty"`
	CreatedAt   primitive.DateTime `bson:"createdAt,omitempty"` // Default needs to be date.now()
}
