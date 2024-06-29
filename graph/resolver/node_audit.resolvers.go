package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.48

import (
	"context"

	"github.com/mikalai2006/geoinfo-api/graph/generated"
	"github.com/mikalai2006/geoinfo-api/graph/model"
)

// ID is the resolver for the id field.
func (r *nodeAuditResolver) ID(ctx context.Context, obj *model.NodeAudit) (string, error) {
	return obj.ID.Hex(), nil
}

// UserID is the resolver for the userId field.
func (r *nodeAuditResolver) UserID(ctx context.Context, obj *model.NodeAudit) (string, error) {
	return obj.UserID.Hex(), nil
}

// NodeID is the resolver for the nodeId field.
func (r *nodeAuditResolver) NodeID(ctx context.Context, obj *model.NodeAudit) (string, error) {
	return obj.NodeID.Hex(), nil
}

// Props is the resolver for the props field.
func (r *nodeAuditResolver) Props(ctx context.Context, obj *model.NodeAudit) (any, error) {
	return obj.Props, nil
}

// NodeAudit returns generated.NodeAuditResolver implementation.
func (r *Resolver) NodeAudit() generated.NodeAuditResolver { return &nodeAuditResolver{r} }

type nodeAuditResolver struct{ *Resolver }
