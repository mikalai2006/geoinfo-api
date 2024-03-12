package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodedataData struct {
	Value interface{} `json:"value" bson:"value"`
}

type Nodedata struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id"`

	NodeID   primitive.ObjectID `json:"nodeId" bson:"node_id"`
	TagID    primitive.ObjectID `json:"tagId" bson:"tag_id"`
	TagoptID primitive.ObjectID `json:"tagoptId" bson:"tagopt_id"`
	// Name     string                       `json:"name" bson:"name"`
	Data        NodedataData                 `json:"data" bson:"data"`
	Title       string                       `json:"title" bson:"title"`
	Description string                       `json:"description" bson:"description"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	Status      int64                        `json:"status" bson:"status"` // 1 - view, 100 - new, -1 - not view(remove)
	Like        int64                        `json:"like" bson:"like"`
	Dlike       int64                        `json:"dlike" bson:"dlike"`

	// Type string `json:"type" bson:"type"`
	User  User            `json:"user,omitempty" bson:"user,omitempty"`
	Audit []NodedataAudit `json:"audit,omitempty" bson:"audit,omitempty"`
	Votes []NodedataVote  `json:"votes,omitempty" bson:"votes,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
type NodedataInputMongo struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id"`

	NodeID      primitive.ObjectID           `json:"nodeId" bson:"node_id"`
	TagID       primitive.ObjectID           `json:"tagId" bson:"tag_id"`
	TagoptID    primitive.ObjectID           `json:"tagoptId" bson:"tagopt_id"`
	Data        NodedataData                 `json:"data" bson:"data"`
	Title       string                       `json:"title" bson:"title"`
	Description string                       `json:"description" bson:"description"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	Status      int64                        `json:"status" bson:"status"` // 1 - view, 100 - new, -1 - not view(remove)
	Like        int64                        `json:"like" bson:"like"`
	Dlike       int64                        `json:"dlike" bson:"dlike"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type NodedataInput struct {
	NodeID   string `json:"nodeId" bson:"node_id"`
	TagID    string `json:"tagId" bson:"tag_id"`
	TagoptID string `json:"tagoptId" bson:"tagopt_id"`
	// Name     string                       `json:"name" bson:"name"`
	// Value       interface{}                  `json:"value" bson:"value"`
	Data        NodedataData                 `json:"data" bson:"data"`
	Title       string                       `json:"title" bson:"title"`
	Description string                       `json:"description" bson:"description"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	Status      int64                        `json:"status" bson:"status"`
	// Type        string                       `json:"type" bson:"type"`
}

type GroupNodeData struct {
	Groups map[string]interface{} `json:"groups" bson:"groups"`
}

type NodedataAudit struct {
	ID         primitive.ObjectID     `json:"_id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID     `json:"userId" bson:"user_id"`
	NodedataID primitive.ObjectID     `json:"nodedataId" bson:"nodedata_id"`
	Value      int64                  `json:"value" bson:"value"`
	Props      map[string]interface{} `json:"props" bson:"props"`

	User User `json:"user,omitempty" bson:"user,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type NodedataAuditDB struct {
	ID         primitive.ObjectID     `json:"_id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID     `json:"userId" bson:"user_id"`
	NodedataID primitive.ObjectID     `json:"nodedataId" bson:"nodedata_id"`
	Value      int64                  `json:"value" bson:"value"`
	Props      map[string]interface{} `json:"props" bson:"props"`
	CreatedAt  time.Time              `json:"createdAt" bson:"created_at"`
	UpdatedAt  time.Time              `json:"updatedAt" bson:"updated_at"`
}

type NodedataAuditInput struct {
	ID         primitive.ObjectID     `json:"_id" bson:"_id,omitempty"`
	UserID     string                 `json:"userId" bson:"user_id"`
	NodedataID string                 `json:"nodedataId" bson:"nodedata_id"`
	Value      int64                  `json:"value" bson:"value"`
	Props      map[string]interface{} `json:"props" bson:"props"`
	CreatedAt  time.Time              `json:"createdAt" bson:"created_at"`
	UpdatedAt  time.Time              `json:"updatedAt" bson:"updated_at"`
}
