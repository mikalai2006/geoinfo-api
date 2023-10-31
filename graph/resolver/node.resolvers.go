package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"fmt"

	"github.com/mikalai2006/geoinfo-api/graph/generated"
	"github.com/mikalai2006/geoinfo-api/graph/loaders"
	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/middleware"
	"github.com/mikalai2006/geoinfo-api/internal/repository"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateNode is the resolver for the createNode field.
func (r *mutationResolver) CreateNode(ctx context.Context, input model.NewNode) (*model.Node, error) {
	return &model.Node{ID: primitive.NewObjectID(), Lat: input.Lat, Lon: input.Lon, OsmID: input.OsmID}, nil
}

// ID is the resolver for the _id field.
func (r *nodeResolver) ID(ctx context.Context, obj *model.Node) (string, error) {
	return obj.ID.Hex(), nil
}

// UserID is the resolver for the userId field.
func (r *nodeResolver) UserID(ctx context.Context, obj *model.Node) (string, error) {
	return obj.UserID.Hex(), nil
}

// TagsData is the resolver for the tagsData field.
func (r *nodeResolver) TagsData(ctx context.Context, obj *model.Node) ([]*model.Tag, error) {
	gc, err := utils.GinContextFromContext(ctx)
	lang := gc.MustGet("i18nLocale").(string)
	if err != nil {
		return nil, err
	}

	listIDs := []primitive.ObjectID{}

	// TODO
	// for i := range obj.Tags {
	// 	uIDPrimitive, err := primitive.ObjectIDFromHex(obj.Tags[i])
	// 	if err != nil {
	// 		return []*model.Tag{}, err
	// 	}
	// 	listIDs = append(listIDs, uIDPrimitive)
	// }

	result, err := r.Repo.Tag.GqlGetTags(domain.RequestParams{
		// Options: domain.Options{Limit: 10},
		Filter: bson.M{"_id": bson.M{"$in": listIDs}},
		Lang:   lang,
	})
	if err != nil {
		return result, err
	}
	// result, errs := loaders.GetTags(ctx, obj.Tags)
	// if len(errs) > 0 {
	// 	fmt.Println("Error:", errs)
	// }
	return result, nil
}

// AmenityID is the resolver for the amenityId field.
func (r *nodeResolver) AmenityID(ctx context.Context, obj *model.Node) (string, error) {
	return obj.ID.Hex(), nil
}

// Props is the resolver for the props field.
func (r *nodeResolver) Props(ctx context.Context, obj *model.Node) (interface{}, error) {
	return obj.Props, nil
}

// User is the resolver for the user field.
func (r *nodeResolver) User(ctx context.Context, obj *model.Node) (*model.User, error) {
	var result *model.User
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return result, err
	}
	lang := gc.MustGet("i18nLocale").(string)

	filter := bson.D{}
	filter = append(filter, bson.E{"_id", obj.UserID})

	allItems, err := r.Repo.User.GqlGetUsers(domain.RequestParams{
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

// Like is the resolver for the like field.
func (r *nodeResolver) Like(ctx context.Context, obj *model.Node) (*model.NodeLike, error) {
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	userID, err := middleware.GetUID(gc)
	// if err != nil {
	// 	return nil, err
	// }
	ilikes, err := r.Repo.Like.GqlGetLikes(domain.RequestParams{
		Options: domain.Options{Limit: 10},
		Filter:  bson.D{{"node_id", obj.ID}},
	})
	if err != nil {
		return nil, err
	}

	l := 0
	dl := 0
	iamlike := model.Like{}
	for _, like := range ilikes {
		if like.Status == 1 {
			l = l + 1
		} else if like.Status == -1 {
			dl = dl + 1
		}

		if like.UserID.Hex() == userID {
			iamlike = *like
		}
	}

	result := model.NodeLike{
		Like:  l,
		Dlike: dl,
		Ilike: iamlike,
	}

	return &result, nil
}

// Reviews is the resolver for the reviews field.
func (r *nodeResolver) Reviews(ctx context.Context, obj *model.Node) ([]*model.Review, error) {
	result, err := r.Repo.Review.GqlGetReviews(domain.RequestParams{
		Options: domain.Options{Limit: 10},
		Filter:  bson.D{{"osm_id", obj.OsmID}},
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

// Address is the resolver for the address field.
func (r *nodeResolver) Address(ctx context.Context, obj *model.Node) (*model.Address, error) {
	// result, err := r.Repo.Address.GqlGetAdresses(domain.RequestParams{
	// 	Options: domain.Options{Limit: 1},
	// 	Filter:  bson.D{{"osm_id", obj.OsmID}},
	// })
	result, err := loaders.GetAddress(ctx, obj.OsmID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ReviewsInfo is the resolver for the reviewsInfo field.
func (r *nodeResolver) ReviewsInfo(ctx context.Context, obj *model.Node) (*model.ReviewInfo, error) {
	var result model.ReviewInfo

	count, err := r.Repo.Review.GqlGetCountReviews(domain.RequestParams{Filter: bson.D{{"osm_id", obj.OsmID}}})
	if err != nil {
		return &result, err
	}
	result = *count

	return &result, nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *nodeResolver) CreatedAt(ctx context.Context, obj *model.Node) (string, error) {
	return obj.CreatedAt.String(), nil
}

// UpdatedAt is the resolver for the updatedAt field.
func (r *nodeResolver) UpdatedAt(ctx context.Context, obj *model.Node) (string, error) {
	return obj.UpdatedAt.String(), nil
}

// Nodes is the resolver for the nodes field.
func (r *queryResolver) Nodes(ctx context.Context, first *int, after *string, limit *int, skip *int, input *model.ParamsNode) (*model.PaginationNode, error) {
	var results *model.PaginationNode

	options := options.Find()
	//options.SetSort(bson.D{{"createdAt", 1}})
	if limit != nil {
		options.SetLimit(int64(*limit))
	}
	options.SetSkip(int64(*skip))
	q := bson.D{}
	if after != nil {
		// idPrimitive, err := primitive.ObjectIDFromHex(*after)
		// if err != nil {
		// 	return results, err
		// }

		// q = append(q, bson.E{"_id", bson.D{{"$lt", idPrimitive}}})
		// q = append(q, bson.E{"createdAt", bson.D{{"$lt", &after}}})
	}

	if input.LatA != nil {
		q = append(q, bson.E{"lat", bson.D{{"$gt", input.LatA}}})
	}
	if input.LatB != nil {
		q = append(q, bson.E{"lat", bson.D{{"$lt", input.LatB}}})
	}
	if input.LonA != nil {
		q = append(q, bson.E{"lon", bson.D{{"$gt", input.LonA}}})
	}
	if input.LonB != nil {
		q = append(q, bson.E{"lon", bson.D{{"$lt", input.LonB}}})
	}
	// if input.Type != nil && len(input.Type) > 0 {
	// 	types := make([]string, len(input.Type))
	// 	for i := range input.Type {
	// 		types[i] = *input.Type[i]
	// 	}
	// 	q = append(q, bson.E{"type", bson.D{{"$in", types}}})
	// }

	// Filter by substring name
	if input.Query != nil && *input.Query != "" {
		strName := primitive.Regex{Pattern: fmt.Sprintf("%v", *input.Query), Options: "i"}
		q = append(q, bson.E{"name", bson.D{{"$regex", strName}}})
		// fmt.Println("q:", q)
	}

	// inputData := []model.NodeFilterTag{}
	// err := json.Unmarshal([]byte(*input.Filter), &inputData)
	// if err != nil {
	// 	return results, err
	// }

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "nodedata",
		// "let":  bson.D{{Key: "nodeId", Value: bson.D{{"$toString", "$_id"}}}},
		// "pipeline": mongo.Pipeline{
		// 	bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$node_id", "$$nodeId"}}}}},
		// },
		"localField":   "_id",
		"foreignField": "node_id",
		"as":           "data",
	}}})

	inputData := input.Filter
	if len(inputData) > 0 {
		// filterOutput := bson.M{}
		// filter := []bson.M{}
		// for i := range inputData {
		// 	// typeFilter := ""
		// 	filterTypeOptions := bson.A{}
		// 	if input.Filter[i].Type != "" {
		// 		// typeFilter = input.Filter[i].Type
		// 		filterTypeOptions = append(filterTypeOptions, bson.M{"$eq": bson.A{"$type", input.Filter[i].Type}})
		// 	} else {
		// 		continue
		// 	}

		// 	if len(inputData[i].Options) > 0 {
		// 		filterOptions := bson.A{}
		// 		for j := range inputData[i].Options {
		// 			tID, _ := primitive.ObjectIDFromHex(inputData[i].Options[j].TagID)

		// 			arrValue := bson.A{}
		// 			for v := range inputData[i].Options[j].Value {
		// 				arrValue = append(arrValue, inputData[i].Options[j].Value[v])
		// 			}

		// 			filterOptions = append(filterOptions, bson.M{
		// 				"$and": bson.A{
		// 					bson.M{
		// 						"$eq": bson.A{"$$item.tag_id", tID},
		// 					},
		// 					bson.M{
		// 						"$in": bson.A{"$$item.data.value", arrValue},
		// 					},
		// 				},
		// 				// "data.tag_id": tID,
		// 				// "data.data.value": bson.D{
		// 				// 	{"$in", arrValue},
		// 				// },
		// 			})
		// 		}
		// 		filterTypeOptions = append(filterTypeOptions, bson.M{"$and": filterOptions})
		// 		// filterTypeOptions["$and"] = filterOptions
		// 	}

		// 	filter = append(filter, bson.M{"$and": filterTypeOptions})
		// }
		// filterOutput["$or"] = filter
		// // filterOutput = append(filterOutput, bson.E{"$or", filter})

		// pipe = append(pipe, bson.D{
		// 	{"$addFields", bson.M{
		// 		"dataCount": bson.M{
		// 			"$size": bson.M{
		// 				"$filter": bson.M{
		// 					"input": "$data",
		// 					"as":    "item",
		// 					"cond":  filterOutput,
		// 				},
		// 			},
		// 		},
		// 	}},
		// })

		// pipe = append(pipe, bson.D{
		// 	{"$match", bson.M{
		// 		"dataCount": bson.M{
		// 			"$gt": 0,
		// 		},
		// 	}}})
		// fmt.Println("pipe=", pipe)

		//////////////////// 2
		filterNodeData := mongo.Pipeline{}
		// filterNodeData = append(filterNodeData, bson.D{{Key: "$lookup", Value: bson.M{
		// 	"from": "node",
		// 	// "let":  bson.D{{Key: "nodeId", Value: bson.D{{"$toString", "$_id"}}}},
		// 	// "pipeline": mongo.Pipeline{
		// 	// 	bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$node_id", "$$nodeId"}}}}},
		// 	// },
		// 	"localField":   "node_id",
		// 	"foreignField": "_id",
		// 	"as":           "node",
		// }}})
		// filterNodeData = append(filterNodeData, bson.D{{"$unwind", bson.D{{"path", "$node"}}}})

		filter := []bson.M{}
		for i := range inputData {
			// typeFilter := ""
			filterTypeOptions := bson.A{}
			if input.Filter[i].Type != "" {
				filterTypeOptions = append(filterTypeOptions, bson.M{"type": input.Filter[i].Type})
			} else {
				continue
			}

			// if input.LatA != nil {
			// 	filterTypeOptions = append(filterTypeOptions, bson.M{"lat": bson.D{{"$gt", *input.LatA}}})
			// }
			// if input.LatB != nil {
			// 	filterTypeOptions = append(filterTypeOptions, bson.M{"lat": bson.D{{"$lt", *input.LatB}}})
			// }
			// if input.LonA != nil {
			// 	filterTypeOptions = append(filterTypeOptions, bson.M{"lon": bson.D{{"$gt", *input.LonA}}})
			// }
			// if input.LonB != nil {
			// 	filterTypeOptions = append(filterTypeOptions, bson.M{"lon": bson.D{{"$lt", *input.LonB}}})
			// }

			if len(inputData[i].Options) > 0 {
				filterOptions := bson.A{}
				for j := range inputData[i].Options {
					tID, _ := primitive.ObjectIDFromHex(inputData[i].Options[j].TagID)

					arrValue := bson.A{}
					for v := range inputData[i].Options[j].Value {
						arrValue = append(arrValue, inputData[i].Options[j].Value[v])
					}

					filterOptions = append(filterOptions, bson.M{
						"tag_id": tID,
						"data.value": bson.D{
							{"$in", arrValue},
						},
					})
				}
				filterTypeOptions = append(filterTypeOptions, bson.M{"$and": filterOptions})
				// filterTypeOptions["$and"] = filterOptions
			}

			filter = append(filter, bson.M{"$and": filterTypeOptions})
		}
		filterNodeData = append(filterNodeData, bson.D{{"$match", bson.D{{"$or", filter}}}})

		fmt.Println("filterOutput=", filterNodeData)

		var allAllowOpts []model.Nodedata
		// if limit != nil {
		// 	filterOutput = append(filterOutput, bson.M{"$limit": limit})
		// }
		cur, err := r.DB.Collection(repository.TblNodedata).Aggregate(ctx, filterNodeData)
		if err != nil {
			return results, err
		}
		if er := cur.All(ctx, &allAllowOpts); er != nil {
			return results, er
		}
		fmt.Println("len=", len(allAllowOpts))
		IDs := []primitive.ObjectID{}
		for e := range allAllowOpts {
			IDs = append(IDs, allAllowOpts[e].NodeID)
		}
		fmt.Println("IDs len=", len(IDs))
		fmt.Println("filterOutput <<<<<=================")

		pipe = append(pipe, bson.D{{"$match", bson.D{{"_id", bson.D{{"$in", IDs}}}}}})
		// fmt.Println("pipe=", pipe)
	}

	// if input.Name != nil && *input.Name != "" {
	// 	var allAllowOpts []model.Nodedata
	// 	fmt.Println("Filter by name=================>>>>>")
	// 	aggregateQuery := []bson.M{}
	// 	strName := primitive.Regex{Pattern: fmt.Sprintf("^%v", *input.Name), Options: "i"}
	// 	fmt.Println("strName: ", strName)
	// 	tID, _ := primitive.ObjectIDFromHex("650712e1d1f4b9ad0b718092")
	// 	aggregateQuery = append(aggregateQuery, bson.M{"$match": bson.D{
	// 		{"$and", bson.A{
	// 			bson.D{{"tag_id", tID}},
	// 			bson.D{{"data.value", bson.D{{"$regex", strName}}}},
	// 		}},
	// 	}})
	// 	fmt.Println("aggregateQuery=", aggregateQuery)
	// 	// aggregateQuery = append(aggregateQuery,
	// 	// 	bson.M{"$lookup": bson.M{
	// 	// 		"from":         "tag",
	// 	// 		"as":           "tagg",
	// 	// 		"foreignField": "tag_id",
	// 	// 		"localField":   "_id",
	// 	// 	},
	// 	// 	})
	// 	if limit != nil {
	// 		aggregateQuery = append(aggregateQuery, bson.M{"$limit": limit})
	// 	}
	// 	cur, err := r.DB.Collection(repository.TblNodedata).Aggregate(ctx, aggregateQuery)
	// 	if err != nil {
	// 		return results, err
	// 	}
	// 	if er := cur.All(ctx, &allAllowOpts); er != nil {
	// 		return results, er
	// 	}
	// 	fmt.Println("len=", len(allAllowOpts))
	// 	IDs := []primitive.ObjectID{}
	// 	for e := range allAllowOpts {
	// 		IDs = append(IDs, allAllowOpts[e].NodeID)
	// 	}
	// 	fmt.Println("IDs len=", len(IDs))
	// 	fmt.Println("Filter by name<<<<<=================")
	// 	q = append(q, bson.E{"_id", bson.D{{"$in", IDs}}})
	// 	// fmt.Println("q=", q)
	// }

	if skip != nil {
		pipe = append(pipe, bson.D{{"$skip", skip}})
	}
	if limit != nil {
		pipe = append(pipe, bson.D{{"$limit", limit}})
	}
	var allItems []model.Node
	cursor, err := r.DB.Collection(repository.TblNode).Aggregate(ctx, pipe) //Aggregate(ctx, pipe) // Find(ctx, q, options)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)
	if er := cursor.All(ctx, &allItems); er != nil {
		return results, er
	}
	// fmt.Println("allItems len=", len(allItems))

	if len(allItems) == 0 {
		return results, nil
	}

	data := make([]*model.Node, len(allItems))
	for i, _ := range allItems {

		data[i] = &allItems[i]
	}

	total := 0
	results = &model.PaginationNode{
		Total: &total,
		Data:  data,
		Limit: limit,
		Skip:  skip,
	}

	return results, nil
}

// Osm is the resolver for the osm field.
func (r *queryResolver) Osm(ctx context.Context, first *int, after *string, limit *int, skip *int, lonA *float64, latA *float64, lonB *float64, latB *float64) ([]*model.Node, error) {
	var results []*model.Node

	options := options.Find()
	options.SetSort(bson.D{{"createdAt", 1}})
	if limit != nil {
		options.SetLimit(int64(*limit))
	}
	options.SetSkip(int64(*skip))
	q := bson.D{}
	if after != nil {
		// idPrimitive, err := primitive.ObjectIDFromHex(*after)
		// if err != nil {
		// 	return results, err
		// }

		// q = append(q, bson.E{"_id", bson.D{{"$lt", idPrimitive}}})
		// q = append(q, bson.E{"createdAt", bson.D{{"$lt", &after}}})
	}

	if latA != nil {
		q = append(q, bson.E{"lat", bson.D{{"$gt", latA}}})
	}
	if latB != nil {
		q = append(q, bson.E{"lat", bson.D{{"$lt", latB}}})
	}
	if lonA != nil {
		q = append(q, bson.E{"lon", bson.D{{"$gt", lonA}}})
	}
	if lonB != nil {
		q = append(q, bson.E{"lon", bson.D{{"$lt", lonB}}})
	}

	var allItems []model.Node
	cursor, err := r.DB.Collection(repository.TblNode).Find(ctx, q, options)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &allItems); er != nil {
		return results, er
	}

	if len(allItems) == 0 {
		return results, nil
	}

	for i, _ := range allItems {
		results = append(results, &allItems[i])
	}

	return results, nil
}

// Node is the resolver for the node field.
func (r *queryResolver) Node(ctx context.Context, id *string, osmID *string) (*model.Node, error) {
	var result *model.Node

	filter := bson.D{}
	if id != nil {
		userIDPrimitive, err := primitive.ObjectIDFromHex(*id)
		if err != nil {
			return result, err
		}

		filter = append(filter, bson.E{"_id", userIDPrimitive})
	} else if osmID != nil {
		filter = append(filter, bson.E{"osm_id", osmID})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", filter}})
	pipe = append(pipe, bson.D{{"$limit", 1}})

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "nodedata",
		// "let":  bson.D{{Key: "nodeId", Value: bson.D{{"$toString", "$_id"}}}},
		// "pipeline": mongo.Pipeline{
		// 	bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$node_id", "$$nodeId"}}}}},
		// },
		"localField":   "_id",
		"foreignField": "node_id",
		"as":           "data",
	}}})

	// if err := r.DB.Collection(repository.TblNode).FindOne(ctx, filter).Decode(&result); err != nil {
	// 	if errors.Is(err, mongo.ErrNoDocuments) {
	// 		return result, model.ErrNodeNotFound
	// 	}
	// 	return result, err
	// }
	var allItems []model.Node
	cursor, err := r.DB.Collection(repository.TblNode).Aggregate(ctx, pipe) //Aggregate(ctx, pipe) // Find(ctx, q, options)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)
	if er := cursor.All(ctx, &allItems); er != nil {
		return result, er
	}

	result = &allItems[0]

	return result, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Node returns generated.NodeResolver implementation.
func (r *Resolver) Node() generated.NodeResolver { return &nodeResolver{r} }

type mutationResolver struct{ *Resolver }
type nodeResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *nodeResolver) Data(ctx context.Context, obj *model.Node) ([]*model.Nodedata, error) {
	gc, err := utils.GinContextFromContext(ctx)
	lang := gc.MustGet("i18nLocale").(string)
	if err != nil {
		return nil, err
	}

	result, err := r.Repo.Nodedata.GqlGetNodedatas(domain.RequestParams{
		// Options: domain.Options{Limit: 10},
		Filter: bson.M{"node_id": obj.ID, "status": bson.M{"$gt": 0}},
		Lang:   lang,
	})
	if err != nil {
		return result, err
	}
	// result, errs := loaders.GetTags(ctx, obj.Tags)
	// if len(errs) > 0 {
	// 	fmt.Println("Error:", errs)
	// }
	return result, nil
}
