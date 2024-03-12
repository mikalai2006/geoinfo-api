package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"

	"github.com/mikalai2006/geoinfo-api/graph/generated"
	"github.com/mikalai2006/geoinfo-api/graph/model"
)

// ID is the resolver for the id field.
func (r *nodedataVoteResolver) ID(ctx context.Context, obj *model.NodedataVote) (string, error) {
	return obj.ID.Hex(), nil
}

// UserID is the resolver for the userId field.
func (r *nodedataVoteResolver) UserID(ctx context.Context, obj *model.NodedataVote) (string, error) {
	return obj.UserID.Hex(), nil
}

// NodedataID is the resolver for the nodedataId field.
func (r *nodedataVoteResolver) NodedataID(ctx context.Context, obj *model.NodedataVote) (string, error) {
	return obj.NodedataID.Hex(), nil
}

// Value is the resolver for the value field.
func (r *nodedataVoteResolver) Value(ctx context.Context, obj *model.NodedataVote) (interface{}, error) {
	return obj.Value, nil
}

// NodedataVote returns generated.NodedataVoteResolver implementation.
func (r *Resolver) NodedataVote() generated.NodedataVoteResolver { return &nodedataVoteResolver{r} }

type nodedataVoteResolver struct{ *Resolver }