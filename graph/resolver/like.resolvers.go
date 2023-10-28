package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"

	"github.com/mikalai2006/geoinfo-api/graph/generated"
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/middleware"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
)

// ID is the resolver for the _id field.
func (r *likeResolver) ID(ctx context.Context, obj *model.Like) (string, error) {
	return obj.ID.Hex(), nil
}

// UserID is the resolver for the userId field.
func (r *likeResolver) UserID(ctx context.Context, obj *model.Like) (string, error) {
	return obj.UserID.Hex(), nil
}

// NodeID is the resolver for the nodeId field.
func (r *likeResolver) NodeID(ctx context.Context, obj *model.Like) (string, error) {
	return obj.NodeID.Hex(), nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *likeResolver) CreatedAt(ctx context.Context, obj *model.Like) (string, error) {
	return obj.CreatedAt.String(), nil
}

// UpdatedAt is the resolver for the updatedAt field.
func (r *likeResolver) UpdatedAt(ctx context.Context, obj *model.Like) (string, error) {
	return obj.UpdatedAt.String(), nil
}

// Like is the resolver for the like field.
func (r *queryResolver) Like(ctx context.Context, nodeID string) (*model.Like, error) {
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	userID, err := middleware.GetUID(gc)
	if err != nil {
		return nil, err
	}
	ilike, err := r.Repo.Like.GqlGetIamLike(userID, nodeID)
	if err != nil {
		return nil, err
	}
	return ilike, nil
}

// Like returns generated.LikeResolver implementation.
func (r *Resolver) Like() generated.LikeResolver { return &likeResolver{r} }

type likeResolver struct{ *Resolver }
