package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Like struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"userId" bson:"user_id"`
	NodeID    primitive.ObjectID `json:"nodeId" bson:"node_id"`
	Status    int                `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updated_at"`
}

type LikeInput struct {
	NodeID string `json:"nodeId" bson:"node_id" form:"nodeId"`
	Status int    `json:"status" bson:"status" form:"status"`
}
