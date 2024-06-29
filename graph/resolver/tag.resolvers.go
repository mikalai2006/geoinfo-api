package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.48

import (
	"context"

	"github.com/mikalai2006/geoinfo-api/graph/generated"
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tags is the resolver for the tags field.
func (r *queryResolver) Tags(ctx context.Context, limit *int, skip *int, input *model.ParamsTag) (*model.PaginationTag, error) {
	var results *model.PaginationTag

	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return results, err
	}
	lang := gc.MustGet("i18nLocale").(string)

	allItems, err := r.Repo.Tag.GqlGetTags(domain.RequestParams{
		Options: domain.Options{Limit: int64(*limit), Skip: int64(*skip), Sort: bson.D{bson.E{"sort_order", 1}}},
		Filter:  bson.D{},
		Lang:    lang,
	})
	if err != nil {
		return results, err
	}

	items := make([]*model.Tag, len(allItems))
	for i, _ := range allItems {
		items[i] = allItems[i]
	}

	count, err := r.DB.Collection(repository.TblTag).CountDocuments(ctx, bson.M{})
	if err != nil {
		return results, err
	}
	countInt := int(count)

	results = &model.PaginationTag{
		Data:  items,
		Total: &countInt,
		Limit: limit,
		Skip:  skip,
	}

	return results, nil
}

// Tag is the resolver for the tag field.
func (r *queryResolver) Tag(ctx context.Context, input *model.ParamsTag) (*model.Tag, error) {
	var result *model.Tag
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return result, err
	}
	lang := gc.MustGet("i18nLocale").(string)

	filter := bson.D{}
	if input.ID != nil {
		userIDPrimitive, err := primitive.ObjectIDFromHex(*input.ID)
		if err != nil {
			return result, err
		}

		filter = append(filter, bson.E{"_id", userIDPrimitive})
	}

	// if err := r.DB.Collection(repository.TblTag).FindOne(ctx, filter).Decode(&result); err != nil {
	// 	if errors.Is(err, mongo.ErrNoDocuments) {
	// 		return result, model.ErrTagNotFound
	// 	}
	// 	return result, err
	// }

	allItems, err := r.Repo.Tag.GqlGetTags(domain.RequestParams{
		Options: domain.Options{Limit: 1, Skip: 0},
		Filter:  filter,
		Lang:    lang,
	})
	if err != nil {
		return result, err
	}

	if len(allItems) > 0 {
		result = allItems[0]
	}

	return result, nil
}

// ID is the resolver for the _id field.
func (r *tagResolver) ID(ctx context.Context, obj *model.Tag) (string, error) {
	return obj.ID.Hex(), nil
}

// UserID is the resolver for the userId field.
func (r *tagResolver) UserID(ctx context.Context, obj *model.Tag) (string, error) {
	return obj.UserID.Hex(), nil
}

// Props is the resolver for the props field.
func (r *tagResolver) Props(ctx context.Context, obj *model.Tag) (any, error) {
	return obj.Props, nil
}

// Tag returns generated.TagResolver implementation.
func (r *Resolver) Tag() generated.TagResolver { return &tagResolver{r} }

type tagResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *tagResolver) CreatedAt(ctx context.Context, obj *model.Tag) (string, error) {
	return obj.CreatedAt.String(), nil
}
func (r *tagResolver) UpdatedAt(ctx context.Context, obj *model.Tag) (string, error) {
	return obj.UpdatedAt.String(), nil
}
func (r *tagResolver) Options(ctx context.Context, obj *model.Tag) ([]*model.Tagopt, error) {
	gc, err := utils.GinContextFromContext(ctx)
	lang := gc.MustGet("i18nLocale").(string)

	var result []*model.Tagopt

	allItems, err := r.Repo.Tagopt.GqlGetTagopts(domain.RequestParams{
		// Options: domain.Options{Limit: int64(*limit), Skip: int64(*skip)},
		Filter: bson.D{{"tag_id", obj.ID}},
		Lang:   lang,
	})
	if err != nil {
		return result, err
	}

	items := make([]*model.Tagopt, len(allItems))
	for i, _ := range allItems {
		items[i] = allItems[i]
	}

	result = items

	// var osmId interface{}
	// var hg *model.Node
	// var ok bool

	// for p := graphql.GetFieldContext(ctx).Parent; ; p = p.Parent {
	// 	hg, ok = p.Result.(*model.Node)
	// 	if ok {
	// 		break
	// 	}
	// }

	// if hg == nil {
	// 	panic("failing to get parent")
	// }
	// osmId = hg.OsmID
	// // IDs := []string{}
	// key := fmt.Sprintf("%v.%v", osmId, obj.ID.Hex())
	// // IDs = append(IDs, key)
	// // fmt.Println("osm_id success get parent", osmId, key, IDs)

	// // rezCtx := graphql.GetFieldContext(ctx)
	// // // fmt.Println(rezCtx.Path())
	// // for k, v := range rezCtx.Parent.Parent.Parent.Args {
	// // 	if k == "osmId" {
	// // 		osmId = v
	// // 		// fmt.Println(k, "=", v)
	// // 	}
	// // }

	// // result, err := r.Repo.Tagopt.GqlGetTagopts(domain.RequestParams{
	// // 	Options: domain.Options{Limit: 10},
	// // 	Filter:  bson.D{{"tag_id", obj.ID}, {"osm_id", osmId}},
	// // })
	// // if err != nil {
	// // 	return result, err
	// // }
	// result, _ := loaders.GetTagopt(ctx, key)
	// // if len(errs) > 0 {
	// // 	return result, errs[0]
	// // }
	// // fmt.Println("result=", len(result))

	return result, nil
}
