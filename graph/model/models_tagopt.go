package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tagopt struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id"`
	TagID  primitive.ObjectID `json:"tagId" bson:"tag_id" binding:"required" primitive:"true"`
	// OsmID     string             `json:"osmId" bson:"osm_id" binding:"required"`
	// Key       string                       `json:"key" bson:"key"`
	Value       string                       `json:"value" bson:"value"`
	Title       string                       `json:"title" bson:"title"`
	Description string                       `json:"description" bson:"description"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	Props       map[string]interface{}       `json:"props" bson:"props"`
	CountItem   int                          `json:"countItem" bson:"countItem"`
	CreatedAt   time.Time                    `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time                    `json:"updatedAt" bson:"updated_at"`
}

type TagoptInput struct {
	ID          primitive.ObjectID           `json:"id" bson:"_id,omitempty"`
	UserID      string                       `json:"userId" bson:"user_id"`
	Value       string                       `json:"value" bson:"value" form:"value"`
	Title       string                       `json:"title" bson:"title"`
	Description string                       `json:"description" bson:"description"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	Props       map[string]interface{}       `json:"props" bson:"props"`
	TagID       string                       `json:"tagId" bson:"tag_id" form:"tagId" primitive:"true"`
	CreatedAt   time.Time                    `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time                    `json:"updatedAt" bson:"updated_at"`
}

// TagID  string             `json:"tagId" bson:"tag_id"`
// OsmID     string             `json:"osmId" bson:"osm_id" binding:"required"`
// Key       string                       `json:"key" bson:"key"`
