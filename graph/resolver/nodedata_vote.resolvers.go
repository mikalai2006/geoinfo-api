package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.48

import (
	"context"

	"github.com/mikalai2006/geoinfo-api/graph/generated"
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// NodeID is the resolver for the nodeId field.
func (r *nodedataVoteResolver) NodeID(ctx context.Context, obj *model.NodedataVote) (string, error) {
	return obj.NodeID.Hex(), nil
}

// NodedataUserID is the resolver for the nodedataUserId field.
func (r *nodedataVoteResolver) NodedataUserID(ctx context.Context, obj *model.NodedataVote) (string, error) {
	return obj.NodedataUserID.Hex(), nil
}

// Value is the resolver for the value field.
func (r *nodedataVoteResolver) Value(ctx context.Context, obj *model.NodedataVote) (any, error) {
	return obj.Value, nil
}

// Nodedatavotes is the resolver for the nodedatavotes field.
func (r *queryResolver) Nodedatavotes(ctx context.Context, limit *int, skip *int, input *model.FetchNodedataVote) (*model.PaginationNodedataVote, error) {
	var results *model.PaginationNodedataVote

	filter := bson.D{}
	if input.NodedataID != nil {
		nID, _ := primitive.ObjectIDFromHex(*input.NodedataID)
		filter = append(filter, bson.E{"nodedata_id", nID})
	}
	if input.UserID != nil {
		tID, _ := primitive.ObjectIDFromHex(*input.UserID)
		filter = append(filter, bson.E{"user_id", tID})
	}
	if input.NodedataUserID != nil {
		nuID, _ := primitive.ObjectIDFromHex(*input.NodedataUserID)
		filter = append(filter, bson.E{"nodedata_user_id", nuID})
	}

	allItems, err := r.Repo.NodedataVote.FindNodedataVote(domain.RequestParams{
		Options: domain.Options{Limit: int64(*limit), Sort: bson.D{{"updated_at", -1}}},
		Filter:  filter,
	})
	if err != nil {
		return results, err
	}

	data := make([]*model.NodedataVote, len(allItems.Data))
	for i, _ := range allItems.Data {
		data[i] = &allItems.Data[i]
	}

	total := len(data)

	results = &model.PaginationNodedataVote{
		Data:  data,
		Total: &total,
		Limit: limit,
		Skip:  skip,
	}

	return results, nil
}

// NodedataVote returns generated.NodedataVoteResolver implementation.
func (r *Resolver) NodedataVote() generated.NodedataVoteResolver { return &nodedataVoteResolver{r} }

type nodedataVoteResolver struct{ *Resolver }
