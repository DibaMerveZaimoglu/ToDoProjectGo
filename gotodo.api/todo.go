package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Todo item structure.
type Todo struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title    string             `json:"title,omitempty"`
	IsActive bool               `json:"isactive,omitempty"`
}
