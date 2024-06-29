package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.48

import (
	"context"
	"fmt"
	"math"

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
	return obj.AmenityID.Hex(), nil
}

// Props is the resolver for the props field.
func (r *nodeResolver) Props(ctx context.Context, obj *model.Node) (any, error) {
	return obj.Props, nil
}

// My is the resolver for the my field.
func (r *nodeResolver) My(ctx context.Context, obj *model.Node) (*bool, error) {
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	userID, err := middleware.GetUID(gc)
	// if err != nil {
	// 	return nil, err
	// }
	my := obj.UserID.Hex() == userID
	return &my, nil
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
		Filter:  bson.D{{"node_id", obj.ID}},
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

	count, err := r.Repo.Review.GqlGetCountReviews(domain.RequestParams{Filter: bson.D{{"node_id", obj.ID}}})
	if err != nil {
		return &result, err
	}
	result = *count

	return &result, nil
}

// Audits is the resolver for the audits field.
func (r *nodeResolver) Audits(ctx context.Context, obj *model.Node) ([]*model.NodeAudit, error) {
	var result []*model.NodeAudit
	audits, err := r.Repo.NodeAudit.FindNodeAudit(domain.RequestParams{
		Options: domain.Options{Limit: 10},
		Filter:  bson.D{{"node_id", obj.ID}},
	})
	if err != nil {
		return nil, err
	}

	for i := range audits.Data {
		result = append(result, &audits.Data[i])
	}

	return result, nil
}

// Nodes is the resolver for the nodes field.
func (r *queryResolver) Nodes(ctx context.Context, limit *int, skip *int, input *model.ParamsNode) (*model.PaginationNode, error) {
	var results *model.PaginationNode

	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	md, err := middleware.GetMaxDistance(gc)
	if err != nil {
		return nil, err
	}
	// fmt.Println("max distance=", md)

	// options := options.Find()
	// //options.SetSort(bson.D{{"createdAt", 1}})
	// if limit != nil {
	// 	options.SetLimit(int64(*limit))
	// }
	// options.SetSkip(int64(*skip))
	q := bson.D{}

	// if input.LatA != nil {
	// 	q = append(q, bson.E{"lat", bson.D{{"$gt", input.LatA}}})
	// }
	// if input.LatB != nil {
	// 	q = append(q, bson.E{"lat", bson.D{{"$lt", input.LatB}}})
	// }
	// if input.LonA != nil {
	// 	q = append(q, bson.E{"lon", bson.D{{"$gt", input.LonA}}})
	// }
	// if input.LonB != nil {
	// 	q = append(q, bson.E{"lon", bson.D{{"$lt", input.LonB}}})
	// }

	if input.Center != nil && len(input.Center) == 2 {

		lat := *input.Center[0]
		lon := *input.Center[1]
		latAccuracy := float64(float64(180*md*1000) / 40075017)
		lngAccuracy := float64(latAccuracy) / math.Cos(float64(math.Pi/180)*lat)

		latA := lat - float64(latAccuracy)
		if input.LatA != nil && latA < *input.LatA {
			latA = *input.LatA
		}
		lonA := lon - lngAccuracy
		if input.LonA != nil && lonA < *input.LonA {
			lonA = *input.LonA
		}
		latB := lat + float64(latAccuracy)
		if input.LatB != nil && latB > *input.LatB {
			latB = *input.LatB
		}
		lonB := lon + lngAccuracy
		if input.LonB != nil && lonB > *input.LonB {
			lonB = *input.LonB
		}
		// fmt.Println("center: ", lat, lon)
		// fmt.Println("latAccuracy: ", latAccuracy)
		// fmt.Println("lngAccuracy: ", lngAccuracy)
		// fmt.Println("bounds: ", latA, lonA, ":", latB, lonB)
		q = append(q, bson.E{"lat", bson.D{{"$gt", latA}}})
		q = append(q, bson.E{"lon", bson.D{{"$gt", lonA}}})
		q = append(q, bson.E{"lat", bson.D{{"$lt", latB}}})
		q = append(q, bson.E{"lon", bson.D{{"$lt", lonB}}})
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

	// Filter by country code
	if input.C != nil && len(input.C) > 0 {
		// strName := primitive.Regex{Pattern: fmt.Sprintf("%v", *input.Query), Options: "i"}
		q = append(q, bson.E{"ccode", bson.D{{"$in", input.C}}})
	}
	// fmt.Println("q:", q, input.C)

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

	type ItemFilter struct {
		Type      string
		CountCond int
	}

	inputData := input.Filter
	if len(inputData) > 0 {
		// -> 1
		itemFilter := []ItemFilter{}
		filterOutput := bson.M{}
		filter := []bson.M{}
		for i := range inputData {
			newItemFilter := ItemFilter{
				Type:      input.Filter[i].Type,
				CountCond: len(inputData[i].Options),
			}
			// typeFilter := ""
			filterTypeOptions := bson.A{}
			if input.Filter[i].Type != "" {
				// typeFilter = input.Filter[i].Type
				filterTypeOptions = append(filterTypeOptions, bson.M{"$eq": bson.A{"$type", input.Filter[i].Type}})
			} else {
				continue
			}

			if len(inputData[i].Options) > 0 {
				filterOptions := bson.A{}
				for j := range inputData[i].Options {
					tID, _ := primitive.ObjectIDFromHex(inputData[i].Options[j].TagID)

					arrValue := bson.A{}
					for v := range inputData[i].Options[j].Value {
						arrValue = append(arrValue, inputData[i].Options[j].Value[v])
					}

					filterOptions = append(filterOptions, bson.M{
						"$and": bson.A{
							bson.M{
								"$eq": bson.A{"$$item.tag_id", tID},
							},
							bson.M{
								"$in": bson.A{"$$item.data.value", arrValue},
							},
						},
						// "data.tag_id": tID,
						// "data.data.value": bson.D{
						// 	{"$in", arrValue},
						// },
					})
				}
				filterTypeOptions = append(filterTypeOptions, bson.M{"$or": filterOptions})
				// filterTypeOptions["$and"] = filterOptions
			}

			filter = append(filter, bson.M{"$and": filterTypeOptions})
			itemFilter = append(itemFilter, newItemFilter)
		}
		filterOutput["$or"] = filter
		// filterOutput = append(filterOutput, bson.E{"$or", filter})

		pipeBranches := bson.A{}
		for i := range itemFilter {
			pipeBranches = append(pipeBranches, bson.M{
				"case": bson.M{
					"$eq": bson.A{"$type", itemFilter[i].Type},
				},
				"then": itemFilter[i].CountCond,
			})
		}
		// fmt.Println(pipeBranches)
		pipe = append(pipe, bson.D{
			{"$addFields", bson.M{
				"countCond": bson.M{
					"$switch": bson.M{
						"branches": pipeBranches,
						"default":  -1,
					},
				},
			}},
		})

		pipe = append(pipe, bson.D{
			{"$addFields", bson.M{
				"dataCount": bson.M{
					"$size": bson.M{
						"$filter": bson.M{
							"input": "$data",
							"as":    "item",
							"cond":  filterOutput,
						},
					},
				},
			}},
		})
		// pipe = append(pipe, bson.D{
		// 	{"$addFields", bson.M{
		// 		"data": bson.M{
		// 			"count": "$dataCount",
		// 		},
		// 	}},
		// })
		pipe = append(pipe, bson.D{
			{"$match", bson.M{
				// "dataCount": bson.M{
				// 	"$gt": 0,
				// },
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$gte": bson.A{"$dataCount", "$countCond"}},
						bson.M{"$ne": bson.A{"$countCond", -1}},
					},
				},
			}}})
		// fmt.Println("pipe=", pipe)

		//////////////////// 2
		// filterNodeData := mongo.Pipeline{}
		// // filterNodeData = append(filterNodeData, bson.D{{Key: "$lookup", Value: bson.M{
		// // 	"from": "node",
		// // 	// "let":  bson.D{{Key: "nodeId", Value: bson.D{{"$toString", "$_id"}}}},
		// // 	// "pipeline": mongo.Pipeline{
		// // 	// 	bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$node_id", "$$nodeId"}}}}},
		// // 	// },
		// // 	"localField":   "node_id",
		// // 	"foreignField": "_id",
		// // 	"as":           "node",
		// // }}})
		// // filterNodeData = append(filterNodeData, bson.D{{"$unwind", bson.D{{"path", "$node"}}}})

		// filter := []bson.M{}
		// for i := range inputData {
		// 	// typeFilter := ""
		// 	filterTypeOptions := bson.A{}
		// 	if input.Filter[i].Type != "" {
		// 		filterTypeOptions = append(filterTypeOptions, bson.M{"type": input.Filter[i].Type})
		// 	} else {
		// 		continue
		// 	}

		// 	// if input.LatA != nil {
		// 	// 	filterTypeOptions = append(filterTypeOptions, bson.M{"lat": bson.D{{"$gt", *input.LatA}}})
		// 	// }
		// 	// if input.LatB != nil {
		// 	// 	filterTypeOptions = append(filterTypeOptions, bson.M{"lat": bson.D{{"$lt", *input.LatB}}})
		// 	// }
		// 	// if input.LonA != nil {
		// 	// 	filterTypeOptions = append(filterTypeOptions, bson.M{"lon": bson.D{{"$gt", *input.LonA}}})
		// 	// }
		// 	// if input.LonB != nil {
		// 	// 	filterTypeOptions = append(filterTypeOptions, bson.M{"lon": bson.D{{"$lt", *input.LonB}}})
		// 	// }

		// 	if len(inputData[i].Options) > 0 {
		// 		filterOptions := bson.A{}
		// 		for j := range inputData[i].Options {
		// 			tID, _ := primitive.ObjectIDFromHex(inputData[i].Options[j].TagID)

		// 			arrValue := bson.A{}
		// 			for v := range inputData[i].Options[j].Value {
		// 				arrValue = append(arrValue, inputData[i].Options[j].Value[v])
		// 			}

		// 			filterOptions = append(filterOptions, bson.M{
		// 				"tag_id": tID,
		// 				"data.value": bson.D{
		// 					{"$in", arrValue},
		// 				},
		// 			})
		// 		}
		// 		filterTypeOptions = append(filterTypeOptions, bson.M{"$and": filterOptions})
		// 		// filterTypeOptions["$and"] = filterOptions
		// 	}

		// 	filter = append(filter, bson.M{"$and": filterTypeOptions})
		// }
		// filterNodeData = append(filterNodeData, bson.D{{"$match", bson.D{{"$or", filter}}}})

		// fmt.Println("filterOutput=", filterNodeData)

		// var allAllowOpts []model.Nodedata
		// // if limit != nil {
		// // 	filterOutput = append(filterOutput, bson.M{"$limit": limit})
		// // }
		// cur, err := r.DB.Collection(repository.TblNodedata).Aggregate(ctx, filterNodeData)
		// if err != nil {
		// 	return results, err
		// }
		// if er := cur.All(ctx, &allAllowOpts); er != nil {
		// 	return results, er
		// }
		// fmt.Println("len=", len(allAllowOpts))
		// IDs := []primitive.ObjectID{}
		// for e := range allAllowOpts {
		// 	IDs = append(IDs, allAllowOpts[e].NodeID)
		// }
		// fmt.Println("IDs len=", len(IDs))
		// fmt.Println("filterOutput <<<<<=================")

		// pipe = append(pipe, bson.D{{"$match", bson.D{{"_id", bson.D{{"$in", IDs}}}}}})
		// // fmt.Println("pipe=", pipe)
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
	// pipe = append(pipe, bson.D{{"$sort", bson.D{{"lat", 1}}}})
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
	// fmt.Println("Find total: ", len(data))
	total := len(data)
	results = &model.PaginationNode{
		Total: &total,
		Data:  data,
		Limit: limit,
		Skip:  skip,
	}

	return results, nil
}

// Osm is the resolver for the osm field.
func (r *queryResolver) Osm(ctx context.Context, limit *int, skip *int, lonA *float64, latA *float64, lonB *float64, latB *float64) ([]*model.Node, error) {
	var results []*model.Node

	options := options.Find()
	options.SetSort(bson.D{{"createdAt", 1}})
	if limit != nil {
		options.SetLimit(int64(*limit))
	}
	options.SetSkip(int64(*skip))
	q := bson.D{}

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
func (r *queryResolver) Node(ctx context.Context, id *string, osmID *string, lat *float64, lon *float64) (*model.Node, error) {
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
	} else if lat != nil && lon != nil {
		filter = append(filter, bson.E{"lat", lat})
		filter = append(filter, bson.E{"lon", lon})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", filter}})
	pipe = append(pipe, bson.D{{"$limit", 1}})

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "users",
		"as":   "userb",
		"let":  bson.D{{Key: "userId", Value: "$user_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
			bson.D{{"$limit", 1}},
			bson.D{{
				Key: "$lookup",
				Value: bson.M{
					"from": "image",
					"as":   "images",
					"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
					},
				},
			}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$userb"}}}})
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"as":   "data",
		"from": "nodedata",
		"let":  bson.D{{Key: "nodeId", Value: "$_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$node_id", "$$nodeId"}}}}},
			bson.D{{Key: "$lookup", Value: bson.M{
				"from": "users",
				"as":   "usera",
				"let":  bson.D{{Key: "userId", Value: "$user_id"}},
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
					bson.D{{"$limit", 1}},
					bson.D{{
						Key: "$lookup",
						Value: bson.M{
							"as":   "images",
							"from": "image",
							"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
							"pipeline": mongo.Pipeline{
								bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
							},
						},
					}},
				},
			}}},
			bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$usera"}}}},

			// tagopt
			bson.D{{Key: "$lookup", Value: bson.M{
				"from": "tagopt",
				"as":   "tagopts",
				"let":  bson.D{{Key: "tagoptId", Value: "$tagopt_id"}},
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$tagoptId"}}}}},
					bson.D{{"$limit", 1}},
					// bson.D{{
					// 	Key: "$lookup",
					// 	Value: bson.M{
					// 		"as":   "images",
					// 		"from": "image",
					// 		"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
					// 		"pipeline": mongo.Pipeline{
					// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
					// 		},
					// 	},
					// }},
				},
			}}},
			bson.D{{Key: "$set", Value: bson.M{"tagopt": bson.M{"$first": "$tagopts"}}}},

			// tag
			bson.D{{Key: "$lookup", Value: bson.M{
				"from": repository.TblTag,
				"as":   "tags",
				"let":  bson.D{{Key: "tagId", Value: "$tag_id"}},
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$tagId"}}}}},
					bson.D{{"$limit", 1}},
					// bson.D{{
					// 	Key: "$lookup",
					// 	Value: bson.M{
					// 		"as":   "images",
					// 		"from": "image",
					// 		"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
					// 		"pipeline": mongo.Pipeline{
					// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
					// 		},
					// 	},
					// }},
				},
			}}},
			bson.D{{Key: "$set", Value: bson.M{"tag": bson.M{"$first": "$tags"}}}},

			// audit section
			bson.D{{Key: "$lookup", Value: bson.M{
				"as":   "audit",
				"from": "nodedata_audit",
				"let":  bson.D{{Key: "nodedataId", Value: "$_id"}},
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$nodedata_id", "$$nodedataId"}}}}},
					bson.D{{Key: "$lookup", Value: bson.M{
						"from": "users",
						"as":   "userc",
						"let":  bson.D{{Key: "userId", Value: "$user_id"}},
						"pipeline": mongo.Pipeline{
							bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
							bson.D{{"$limit", 1}},
							bson.D{{
								Key: "$lookup",
								Value: bson.M{
									"as":   "images",
									"from": "image",
									"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
									"pipeline": mongo.Pipeline{
										bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
									},
								},
							}},
						},
					}}},
					bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$userc"}}}},
				},
				// "localField":   "_id",
				// "foreignField": "node_id",
			}}},
		},
		// "localField":   "_id",
		// "foreignField": "node_id",
	}}})

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"as":   "images",
		"from": "image",
		"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$and": bson.A{
				bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}},
				bson.M{"$expr": bson.M{"$eq": [2]string{"$service", "node"}}},
			}},
			}},
			bson.D{{Key: "$lookup", Value: bson.M{
				"from": "users",
				"as":   "usera",
				"let":  bson.D{{Key: "userId", Value: "$user_id"}},
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
					bson.D{{"$limit", 1}},
					bson.D{{
						Key: "$lookup",
						Value: bson.M{
							"as":   "images",
							"from": "image",
							"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
							"pipeline": mongo.Pipeline{
								bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
							},
						},
					}},
				},
			}}},
			bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$usera"}}}},
		},

		// "localField":   "_id",
		// "foreignField": "node_id",
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

	if len(allItems) > 0 {
		result = &allItems[0]
	}

	return result, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Node returns generated.NodeResolver implementation.
func (r *Resolver) Node() generated.NodeResolver { return &nodeResolver{r} }

type mutationResolver struct{ *Resolver }
type nodeResolver struct{ *Resolver }
