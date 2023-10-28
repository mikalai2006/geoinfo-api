package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`
	OsmID  string             `json:"osmId" bson:"osm_id"`

	Review string `json:"review" bson:"review"`
	Rate   int    `json:"rate" bson:"rate"`

	Publish   bool      `json:"publish" bson:"publish"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type ReviewInputData struct {
	OsmID string `json:"osmId" bson:"osm_id"  form:"osmId"`

	Review string `json:"review" bson:"review"  form:"review"`
	Rate   int    `json:"rate" bson:"rate" form:"rate"`

	Publish   bool      `json:"publish" bson:"publish" form:"publish"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
