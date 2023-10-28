package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Node struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id"`
	Lon    float64            `json:"lon" bson:"lon"`
	Lat    float64            `json:"lat" bson:"lat"`
	Type   string             `json:"type" bson:"type"`
	Name   string             `json:"name" bson:"name"`
	// Status int64              `json:"status" bson:"status"` // 1 - view, 100 - new, -1 - not view(remove)
	// Tags   []string           `json:"tags" bson:"tags"`
	// Like   int                `json:"like" bson:"like"`
	// Dlike  int                `json:"dlike" bson:"dlike"`
	//TagsData interface{}        `json:"tagsData" bson:"tags_data"`
	OsmID     string             `json:"osmId" bson:"osm_id"`
	AmenityID primitive.ObjectID `json:"amenityId" bson:"amenity_id"`
	// Amenity   []string           `json:"amenity" bson:"amenity"`
	Props     map[string]interface{} `json:"props" bson:"props"`
	CreatedAt time.Time              `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time              `json:"updatedAt" bson:"updated_at"`
}

type NodeLike struct {
	Like  int  `json:"like,omitempty"`
	Dlike int  `json:"dlike,omitempty"`
	Ilike Like `json:"ilike,omitempty"`
}
