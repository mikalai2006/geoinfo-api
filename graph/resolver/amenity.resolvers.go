package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.48

import (
	"context"
	"errors"
	"fmt"

	"github.com/mikalai2006/geoinfo-api/graph/generated"
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ID is the resolver for the _id field.
func (r *amenityResolver) ID(ctx context.Context, obj *model.Amenity) (string, error) {
	return obj.ID.Hex(), nil
}

// UserID is the resolver for the userId field.
func (r *amenityResolver) UserID(ctx context.Context, obj *model.Amenity) (string, error) {
	return obj.UserID.Hex(), nil
}

// Props is the resolver for the props field.
func (r *amenityResolver) Props(ctx context.Context, obj *model.Amenity) (any, error) {
	return obj.Props, nil
}

// Tags is the resolver for the tags field.
func (r *amenityResolver) Tags(ctx context.Context, obj *model.Amenity) ([]*string, error) {
	result := []*string{}
	for i, _ := range obj.Tags {
		s := fmt.Sprintf("%v", obj.Tags[i])
		result = append(result, &s)
	}
	return result, nil
}

// Amenities is the resolver for the amenities field.
func (r *queryResolver) Amenities(ctx context.Context, limit *int, skip *int, input *model.FetchAmenity) (*model.PaginationAmenity, error) {
	var results *model.PaginationAmenity
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return results, err
	}
	lang := gc.MustGet("i18nLocale").(string)

	allItems, err := r.Repo.Amenity.GqlGetAmenitys(domain.RequestParams{
		Options: domain.Options{Limit: int64(*limit)},
		Filter:  bson.D{},
		Lang:    lang,
	})
	if err != nil {
		return results, err
	}

	items := make([]*model.Amenity, len(allItems))
	for i, _ := range allItems {

		items[i] = allItems[i]
	}

	count, err := r.DB.Collection(repository.TblAmenity).CountDocuments(ctx, bson.M{})
	if err != nil {
		return results, err
	}
	countInt := int(count)

	results = &model.PaginationAmenity{
		Data:  items,
		Total: &countInt,
	}
	if skip != nil {
		results.Skip = skip
	}
	if limit != nil {
		results.Limit = limit
	}

	return results, nil
}

// Amenity is the resolver for the amenity field.
func (r *queryResolver) Amenity(ctx context.Context, id *string) (*model.Amenity, error) {
	var result *model.Amenity
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return result, err
	}
	lang := gc.MustGet("i18nLocale").(string)

	filter := bson.D{}
	if id != nil {
		userIDPrimitive, err := primitive.ObjectIDFromHex(*id)
		if err != nil {
			return result, err
		}

		filter = append(filter, bson.E{"_id", userIDPrimitive})
	}

	// if err := r.DB.Collection(repository.TblAmenity).FindOne(ctx, filter).Decode(&result); err != nil {
	// 	if errors.Is(err, mongo.ErrNoDocuments) {
	// 		return result, model.ErrAmenityNotFound
	// 	}
	// 	return result, err
	// }
	allItems, err := r.Repo.Amenity.GqlGetAmenitys(domain.RequestParams{
		Options: domain.Options{Limit: 1},
		Filter:  bson.D{},
		Lang:    lang,
	})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return result, model.ErrAmenityNotFound
	}

	result = allItems[0]

	return result, nil
}

// Amenity returns generated.AmenityResolver implementation.
func (r *Resolver) Amenity() generated.AmenityResolver { return &amenityResolver{r} }

type amenityResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *amenityResolver) CreatedAt(ctx context.Context, obj *model.Amenity) (string, error) {
	return obj.CreatedAt.String(), nil
}
func (r *amenityResolver) UpdatedAt(ctx context.Context, obj *model.Amenity) (string, error) {
	return obj.UpdatedAt.String(), nil
}
