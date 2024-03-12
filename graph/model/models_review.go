package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaginationReview struct {
	Total int       `json:"total,omitempty"`
	Limit int       `json:"limit,omitempty"`
	Skip  int       `json:"skip,omitempty"`
	Data  []*Review `json:"data,omitempty"`
}

// type Review struct {
// 	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
// 	UserID string             `json:"userId" bson:"user_id"`
// 	OsmID  string             `json:"osmId" bson:"osm_id"`
// 	Review string             `json:"review" bson:"review"`
// 	Rate   int                `json:"rate" bson:"rate"`

// 	User User `json:"user" bson:"user,omitempty"`

// 	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
// 	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
// }

type Review struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`
	OsmID  string             `json:"osmId" bson:"osm_id"`

	Review string `json:"review" bson:"review"`
	Rate   int    `json:"rate" bson:"rate"`

	Publish   bool      `json:"publish" bson:"publish"`
	User      User      `json:"user,omitempty" bson:"user,omitempty"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type ReviewInput struct {
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

// type ReviewRate struct {
// 	r1 int `json:"r1" bson:"r1"`
// 	r2 int `json:"r2" bson:"r2"`
// 	r3 int `json:"r3" bson:"r3"`
// 	r4 int `json:"r4" bson:"r4"`
// 	r5 int `json:"r5" bson:"r5"`
// }

// type ReviewInfo struct {
// 	Count int `json:"count" bson:"count"`
// }
