// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type FetchAction struct {
	ID        *string `json:"id,omitempty"`
	Service   *string `json:"service,omitempty"`
	ServiceID *string `json:"serviceId,omitempty"`
	UserID    *string `json:"userId,omitempty"`
	Type      *int    `json:"type,omitempty"`
	Status    *int    `json:"status,omitempty"`
}

type FetchAddress struct {
	ID     *string `json:"id,omitempty"`
	UserID *string `json:"userId,omitempty"`
	OsmID  *string `json:"osmId,omitempty"`
}

type FetchAmenity struct {
	ID     *string `json:"id,omitempty"`
	UserID *string `json:"userId,omitempty"`
	Key    *string `json:"key,omitempty"`
	Type   *string `json:"type,omitempty"`
	Status *int    `json:"status,omitempty"`
}

type FetchImage struct {
	ID        *string `json:"id,omitempty"`
	ServiceID *string `json:"serviceId,omitempty"`
	Service   *string `json:"service,omitempty"`
	UserID    *string `json:"userId,omitempty"`
}

type FetchNodedata struct {
	ID       *string `json:"id,omitempty"`
	UserID   *string `json:"userId,omitempty"`
	NodeID   *string `json:"nodeId,omitempty"`
	TagID    *string `json:"tagId,omitempty"`
	TagoptID *string `json:"tagoptId,omitempty"`
	Value    *string `json:"value,omitempty"`
	Type     *string `json:"type,omitempty"`
}

type FetchReview struct {
	ID     *string `json:"id,omitempty"`
	UserID *string `json:"userId,omitempty"`
	OsmID  *string `json:"osmId,omitempty"`
}

type FetchTag struct {
	ID     *string `json:"id,omitempty"`
	UserID *string `json:"userId,omitempty"`
}

type FetchTagopt struct {
	ID     *string `json:"id,omitempty"`
	UserID *string `json:"userId,omitempty"`
	Value  *string `json:"value,omitempty"`
}

type Hello struct {
	Hello string `json:"hello"`
}

type NewNode struct {
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Type  string  `json:"type"`
	OsmID string  `json:"osmId"`
}

type PageInfo struct {
	StartCursor string `json:"startCursor"`
	EndCursor   string `json:"endCursor"`
	HasNextPage *bool  `json:"hasNextPage,omitempty"`
}

type PaginationAction struct {
	Total *int      `json:"total,omitempty"`
	Limit *int      `json:"limit,omitempty"`
	Skip  *int      `json:"skip,omitempty"`
	Data  []*Action `json:"data,omitempty"`
}

type PaginationAddress struct {
	Total *int       `json:"total,omitempty"`
	Limit *int       `json:"limit,omitempty"`
	Skip  *int       `json:"skip,omitempty"`
	Data  []*Address `json:"data,omitempty"`
}

type PaginationAmenity struct {
	Total *int       `json:"total,omitempty"`
	Limit *int       `json:"limit,omitempty"`
	Skip  *int       `json:"skip,omitempty"`
	Data  []*Amenity `json:"data,omitempty"`
}

type PaginationImage struct {
	Total *int     `json:"total,omitempty"`
	Limit *int     `json:"limit,omitempty"`
	Skip  *int     `json:"skip,omitempty"`
	Data  []*Image `json:"data,omitempty"`
}

type PaginationNode struct {
	Total *int    `json:"total,omitempty"`
	Limit *int    `json:"limit,omitempty"`
	Skip  *int    `json:"skip,omitempty"`
	Data  []*Node `json:"data,omitempty"`
}

type PaginationNodedata struct {
	Total *int        `json:"total,omitempty"`
	Limit *int        `json:"limit,omitempty"`
	Skip  *int        `json:"skip,omitempty"`
	Data  []*Nodedata `json:"data,omitempty"`
}

type PaginationTag struct {
	Total *int   `json:"total,omitempty"`
	Limit *int   `json:"limit,omitempty"`
	Skip  *int   `json:"skip,omitempty"`
	Data  []*Tag `json:"data,omitempty"`
}

type PaginationTagopt struct {
	Total *int      `json:"total,omitempty"`
	Limit *int      `json:"limit,omitempty"`
	Skip  *int      `json:"skip,omitempty"`
	Data  []*Tagopt `json:"data,omitempty"`
}

type ParamsNode struct {
	ID     *string          `json:"id,omitempty"`
	LonA   *float64         `json:"lonA,omitempty"`
	LatA   *float64         `json:"latA,omitempty"`
	LonB   *float64         `json:"lonB,omitempty"`
	LatB   *float64         `json:"latB,omitempty"`
	Query  *string          `json:"query,omitempty"`
	Filter []*NodeFilterTag `json:"filter,omitempty"`
}

type ReviewEdge struct {
	Cursor string  `json:"cursor"`
	Node   *Review `json:"node,omitempty"`
}

type ReviewInfo struct {
	Count   *int        `json:"count,omitempty"`
	Value   *int        `json:"value,omitempty"`
	Ratings interface{} `json:"ratings,omitempty"`
}

type ReviewsConnection struct {
	Edges    []*ReviewEdge `json:"edges"`
	PageInfo *PageInfo     `json:"pageInfo"`
}
