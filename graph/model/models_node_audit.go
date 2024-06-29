package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodeAudit struct {
	ID      primitive.ObjectID     `json:"_id" bson:"_id,omitempty"`
	UserID  primitive.ObjectID     `json:"userId" bson:"user_id"`
	NodeID  primitive.ObjectID     `json:"nodeId" bson:"node_id"`
	Status  int                    `json:"status" bson:"status"`
	Message string                 `json:"message" bson:"message"`
	Props   map[string]interface{} `json:"props" bson:"props"`

	User User `json:"user,omitempty" bson:"user,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
type NodeAuditInput struct {
	ID      primitive.ObjectID     `json:"_id" bson:"_id,omitempty"`
	UserID  primitive.ObjectID     `json:"userId" bson:"user_id"`
	NodeID  primitive.ObjectID     `json:"nodeId" bson:"node_id"`
	Status  int                    `json:"status" bson:"status"`
	Message string                 `json:"message" bson:"message"`
	Props   map[string]interface{} `json:"props" bson:"props"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
