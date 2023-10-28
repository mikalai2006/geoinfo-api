package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID          primitive.ObjectID           `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID           `json:"userId" bson:"user_id"`
	Key         string                       `json:"key" bson:"key"`
	Type        string                       `json:"type" bson:"type"`
	MultiOpt    int64                        `json:"multiOpt" bson:"multi_opt"`
	Title       string                       `json:"title" bson:"title"`
	Description string                       `json:"description" bson:"description"`
	Props       map[string]interface{}       `json:"props" bson:"props"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	// Filter        int                          `json:"filter" bson:"filter"`
	Multilanguage bool `json:"multilanguage" bson:"multilanguage"`
	// TagoptID      []string  `json:"tagoptId" bson:"tagopt_id"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type TagInput struct {
	ID          primitive.ObjectID           `json:"id" bson:"_id,omitempty"`
	UserID      string                       `json:"userId" bson:"user_id" form:"userId"`
	Key         string                       `json:"key" bson:"key" form:"key"`
	Type        string                       `json:"type" bson:"type" form:"type"`
	MultiOpt    int64                        `json:"multiOpt" bson:"multi_opt"`
	Title       string                       `json:"title" bson:"title" form:"title"`
	Description string                       `json:"description" bson:"description" form:"description"`
	Props       map[string]interface{}       `json:"props" bson:"props"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	// Filter        int                          `json:"filter" bson:"filter"`
	Multilanguage bool `json:"multilanguage" bson:"multilanguage"`
	// TagoptID      []string  `json:"tagoptId" bson:"tagopt_id"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
