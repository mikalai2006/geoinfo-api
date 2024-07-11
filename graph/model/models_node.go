package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Node struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"userId" bson:"user_id"`
	AmenityID primitive.ObjectID `json:"amenityId" bson:"amenity_id"`
	Lon       float64            `json:"lon" bson:"lon"`
	Lat       float64            `json:"lat" bson:"lat"`
	Type      string             `json:"type" bson:"type"`
	Name      string             `json:"name" bson:"name"`
	CCode     string             `json:"ccode" bson:"ccode"`
	// Status int64              `json:"status" bson:"status"` // 1 - view, 100 - new, -1 - not view(remove)
	// Tags   []string           `json:"tags" bson:"tags"`
	// Like   int                `json:"like" bson:"like"`
	// Dlike  int                `json:"dlike" bson:"dlike"`
	//TagsData interface{}        `json:"tagsData" bson:"tags_data"`
	// Amenity   []string           `json:"amenity" bson:"amenity"`
	Props map[string]interface{} `json:"props" bson:"props"`
	OsmID string                 `json:"osmId" bson:"osm_id"`
	// Status int64                  `json:"status" bson:"status"` // 1 - view, 100 - new, -1 - not view(remove)
	// Like   int64                  `json:"like" bson:"like"`
	// Dlike  int64                  `json:"dlike" bson:"dlike"`
	NodeLike NodeLike `json:"nodeLike" bson:"node_like"`

	Data   []Nodedata `json:"data" bson:"data,omitempty"`
	Images []Image    `json:"images" bson:"images,omitempty"`
	User   User       `json:"user" bson:"user,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type NodeInputData struct {
	UserID    string                 `json:"userId" bson:"user_id" form:"userId" primitive:"true"`
	AmenityID string                 `json:"amenityId" bson:"amenity_id" primitive:"true"`
	Lon       float64                `json:"lon" bson:"lon" form:"lon"`
	Lat       float64                `json:"lat" bson:"lat" form:"lat"`
	Type      string                 `json:"type" bson:"type" form:"type"`
	Name      string                 `json:"name" bson:"name" form:"name"`
	CCode     string                 `json:"ccode" bson:"ccode" form:"ccode"`
	Props     map[string]interface{} `json:"props" bson:"props" form:"props"`
	OsmID     string                 `json:"osmId" bson:"osm_id" form:"osmId"`

	NodeLike NodeLike `json:"nodeLike" bson:"node_like"`
	// Amenity []string          `json:"amenity" bson:"amenity" form:"amenity"`
	// Tags    []string          `json:"tags" bson:"tags" form:"tags"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type NodeLike struct {
	// Like  int  `json:"like,omitempty"`
	// Dlike int  `json:"dlike,omitempty"`
	// Ilike Like `json:"ilike,omitempty"`

	Status int64 `json:"status" bson:"status"` // 1 - view, 100 - new, -1 - not view(remove)
	Like   int64 `json:"like" bson:"like"`
	Dlike  int64 `json:"dlike" bson:"dlike"`
}

type NodeFilterTagOption struct {
	TagID string        `json:"tagId" bson:"tag_id"`
	Value []interface{} `json:"value" bson:"value"`
}

type NodeFilterTag struct {
	Type    string                `json:"type" bson:"type"`
	Options []NodeFilterTagOption `json:"options" bson:"options"`
}

type NodeInput struct {
	ID        primitive.ObjectID     `json:"_id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID     `json:"userId" bson:"user_id"`
	AmenityID primitive.ObjectID     `json:"amenityId" bson:"amenity_id"`
	Lon       float64                `json:"lon" bson:"lon"`
	Lat       float64                `json:"lat" bson:"lat"`
	Type      string                 `json:"type" bson:"type"`
	Name      string                 `json:"name" bson:"name"`
	CCode     string                 `json:"ccode" bson:"ccode"`
	Props     map[string]interface{} `json:"props" bson:"props"`
	OsmID     string                 `json:"osmId" bson:"osm_id"`

	NodeLike NodeLike `json:"nodeLike" bson:"node_like"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
