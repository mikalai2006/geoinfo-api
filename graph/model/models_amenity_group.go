package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AmenityGroup struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id"`
	// Key         string                       `json:"key" bson:"key"`
	Title       string                       `json:"title" bson:"title"`
	Description string                       `json:"description" bson:"description"`
	Props       map[string]interface{}       `json:"props" bson:"props"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	// Amenitys    []interface{}                `json:"amenitys" bson:"amenitys"`
	Status    int64     `json:"status" bson:"status"` // 1 - view, 100 - new, -1 - not view(remove)
	SortOrder int64     `json:"sortOrder" bson:"sort_order"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type AmenityGroupInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID string             `json:"userId" bson:"user_id" form:"userId"`
	// Key         string                       `json:"key" bson:"key" form:"key"`
	Title       string                       `json:"title" bson:"title" form:"title"`
	Description string                       `json:"description" bson:"description" form:"description"`
	Props       map[string]interface{}       `json:"props" bson:"props"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	SortOrder   int64                        `json:"sortOrder" bson:"sort_order"`
	// Amenitys    []interface{}                `json:"amenitys" bson:"amenitys" form:"amenitys"`
	Status int64 `json:"status" bson:"status"`
}
