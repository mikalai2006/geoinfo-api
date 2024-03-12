package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodedataVote struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id"`

	NodedataID primitive.ObjectID `json:"nodedataId" bson:"nodedata_id"`
	Value      int                `json:"value" bson:"value"`
	// Status     int64              `json:"status" bson:"status"` // 1 - view, 100 - new, -1 - not view(remove)
	// User  User            `json:"user,omitempty" bson:"user,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
type NodedataVoteMongo struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id"`

	NodedataID primitive.ObjectID `json:"nodedataId" bson:"nodedata_id"`
	Value      int                `json:"value" bson:"value"`
	// Status     int64              `json:"status" bson:"status"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type NodedataVoteInput struct {
	NodedataID string `json:"nodedataId" bson:"nodedata_id"`
	Value      int    `json:"value" bson:"value"`
	// Status     int64  `json:"status" bson:"status"`
}
