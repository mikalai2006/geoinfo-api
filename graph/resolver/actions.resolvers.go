package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"fmt"

	"github.com/mikalai2006/geoinfo-api/graph/generated"
	"github.com/mikalai2006/geoinfo-api/graph/model"
)

// ID is the resolver for the id field.
func (r *actionResolver) ID(ctx context.Context, obj *model.Action) (string, error) {
	panic(fmt.Errorf("not implemented: ID - id"))
}

// UserID is the resolver for the userId field.
func (r *actionResolver) UserID(ctx context.Context, obj *model.Action) (string, error) {
	panic(fmt.Errorf("not implemented: UserID - userId"))
}

// ServiceID is the resolver for the serviceId field.
func (r *actionResolver) ServiceID(ctx context.Context, obj *model.Action) (string, error) {
	panic(fmt.Errorf("not implemented: ServiceID - serviceId"))
}

// Props is the resolver for the props field.
func (r *actionResolver) Props(ctx context.Context, obj *model.Action) (interface{}, error) {
	panic(fmt.Errorf("not implemented: Props - props"))
}

// CreatedAt is the resolver for the createdAt field.
func (r *actionResolver) CreatedAt(ctx context.Context, obj *model.Action) (string, error) {
	panic(fmt.Errorf("not implemented: CreatedAt - createdAt"))
}

// UpdatedAt is the resolver for the updatedAt field.
func (r *actionResolver) UpdatedAt(ctx context.Context, obj *model.Action) (string, error) {
	panic(fmt.Errorf("not implemented: UpdatedAt - updatedAt"))
}

// Actions is the resolver for the actions field.
func (r *queryResolver) Actions(ctx context.Context, limit *int, skip *int, input *model.FetchAction) (*model.PaginationAction, error) {
	panic(fmt.Errorf("not implemented: Actions - actions"))
}

// Action is the resolver for the action field.
func (r *queryResolver) Action(ctx context.Context, input *model.FetchImage) (*model.Action, error) {
	panic(fmt.Errorf("not implemented: Action - action"))
}

// Action returns generated.ActionResolver implementation.
func (r *Resolver) Action() generated.ActionResolver { return &actionResolver{r} }

type actionResolver struct{ *Resolver }