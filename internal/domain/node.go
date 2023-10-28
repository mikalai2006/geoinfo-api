package domain

import (
	"time"
)

// type Node struct {
// 	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
// 	UserID primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`

// 	OsmID   string              `json:"osmId" bson:"osm_id" binding:"required"`
// 	Lon     float64             `json:"lon" bson:"lon" binding:"required"`
// 	Lat     float64             `json:"lat" bson:"lat" binding:"required"`
// 	Type    string              `json:"type" bson:"type" binding:"required"`
// 	Tags    []string `json:"tags" bson:"tags" binding:"required"`
// 	Props   map[string]string   `json:"props" bson:"props"`     //  binding:"required"
// 	Amenity []string            `json:"amenity" bson:"amenity"` //  binding:"required"

// 	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
// 	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
// }

type NodeInputData struct {
	OsmID   string            `json:"osmId" bson:"osm_id" form:"osmId"`
	Lon     float64           `json:"lon" bson:"lon" form:"lon"`
	Lat     float64           `json:"lat" bson:"lat" form:"lat"`
	Type    string            `json:"type" bson:"type" form:"type"`
	Tags    []string          `json:"tags" bson:"tags" form:"tags"`
	Props   map[string]string `json:"props" bson:"props" form:"props"`
	Amenity []string          `json:"amenity" bson:"amenity" form:"amenity"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
