package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	GeoID   string                 `json:"geoId" bson:"geo_id"`
	Address map[string]interface{} `json:"address" bson:"address"`
	Lang    string                 `json:"lang" bson:"lang"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type AddressInput struct {
	GeoID   string                 `json:"geoId" bson:"geo_id"  form:"geoId"`
	Address map[string]interface{} `json:"address" bson:"address"  form:"address"`
	Lang    string                 `json:"lang" bson:"lang"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
