package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStat struct {
	Node            int64 `json:"node" bson:"node"`
	NodeLike        int64 `json:"nodeLike" bson:"nodeLike"`
	NodeDLike       int64 `json:"nodeDLike" bson:"nodeDLike"`
	NodeAuthorLike  int64 `json:"nodeAuthorLike" bson:"nodeAuthorLike"`
	NodeAuthorDLike int64 `json:"nodeAuthorDLike" bson:"nodeAuthorDLike"`

	Nodedata            int64 `json:"nodedata" bson:"nodedata"`
	NodedataLike        int64 `json:"nodedataLike" bson:"nodedataLike"`
	NodedataDLike       int64 `json:"nodedataDLike" bson:"nodedataDLike"`
	NodedataAuthorLike  int64 `json:"nodedataAuthorLike" bson:"nodedataAuthorLike"`
	NodedataAuthorDLike int64 `json:"nodedataAuthorDLike" bson:"nodedataAuthorDLike"`

	Review int64 `json:"review" bson:"review"`
}

type User struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" primitive:"true"`
	UserID primitive.ObjectID `json:"userId,omitempty" bson:"user_id,omitempty" primitive:"true"`

	Name     string `json:"name" bson:"name" form:"name"`
	Login    string `json:"login" bson:"login" form:"login"`
	Currency string `json:"currency" bson:"currency" form:"currency"`
	Lang     string `json:"lang" bson:"lang" form:"lang"`
	Avatar   string `json:"avatar" bson:"avatar"`
	Online   bool   `json:"online" bson:"online" form:"online"`
	Verify   bool   `json:"verify" bson:"verify"`

	UserStat UserStat `json:"userStat" bson:"user_stat"`

	Roles  []string `json:"roles" bson:"-"`
	Md     int      `json:"md" bson:"-"`
	Images []Image  `json:"images,omitempty" bson:"images,omitempty"`

	LastTime  time.Time `json:"lastTime" bson:"last_time"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type UserInput struct {
	ID       string `json:"id" bson:"_id" form:"id" primitive:"true"`
	UserID   string `json:"userId" bson:"user_id" form:"userId" primitive:"true"`
	Name     string `json:"name" bson:"name" form:"name"`
	Login    string `json:"login" bson:"login" form:"login"`
	Currency string `json:"currency" bson:"currency" form:"currency"`
	Lang     string `json:"lang" bson:"lang" form:"lang"`
	Avatar   string `json:"avatar" bson:"avatar" form:"avatar"`
}
