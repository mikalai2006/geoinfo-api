package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Amenity struct {
	ID          primitive.ObjectID           `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID           `json:"userId" bson:"user_id"`
	Key         string                       `json:"key" bson:"key"`
	Group       string                       `json:"group" bson:"group"`
	Title       string                       `json:"title" bson:"title"`
	Description string                       `json:"description" bson:"description"`
	Props       map[string]interface{}       `json:"props" bson:"props"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	Type        string                       `json:"type" bson:"type"`
	Tags        []interface{}                `json:"tags" bson:"tags"`
	Status      int64                        `json:"status" bson:"status"` // 1 - view, 100 - new, -1 - not view(remove)
	CreatedAt   time.Time                    `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time                    `json:"updatedAt" bson:"updated_at"`
}

type AmenityInput struct {
	ID          primitive.ObjectID           `json:"id" bson:"_id,omitempty"`
	UserID      string                       `json:"userId" bson:"user_id" form:"userId"`
	Key         string                       `json:"key" bson:"key" form:"key"`
	Group       string                       `json:"group" bson:"group" form:"group"`
	Title       string                       `json:"title" bson:"title" form:"title"`
	Description string                       `json:"description" bson:"description" form:"description"`
	Props       map[string]interface{}       `json:"props" bson:"props"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	Type        string                       `json:"type" bson:"type"`
	Tags        []interface{}                `json:"tags" bson:"tags" form:"tags"`
	Status      int64                        `json:"status" bson:"status"`
	CreatedAt   time.Time                    `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time                    `json:"updatedAt" bson:"updated_at"`
}
