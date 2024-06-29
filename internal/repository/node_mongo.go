package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/config"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NodeMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewNodeMongo(db *mongo.Database, i18n config.I18nConfig) *NodeMongo {
	return &NodeMongo{db: db, i18n: i18n}
}

func (r *NodeMongo) FindNode(params domain.RequestParams) (domain.Response[model.Node], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Node
	var response domain.Response[model.Node]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.Node]{}, err
	// }
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.Node]{}, err
	}
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

	cursor, err := r.db.Collection(TblNode).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Node, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblNode).CountDocuments(ctx, params.Filter)
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Node]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *NodeMongo) GetAllNode(params domain.RequestParams) (domain.Response[model.Node], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Node
	var response domain.Response[model.Node]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.Node]{}, err
	}

	cursor, err := r.db.Collection(TblNode).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Node, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblNode).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Node]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *NodeMongo) CreateNode(userID string, node *model.Node) (*model.Node, error) {
	var result *model.Node

	collection := r.db.Collection(TblNode)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	createdAt := node.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	newNode := model.NodeInput{
		UserID:    userIDPrimitive,
		OsmID:     node.OsmID,
		Lon:       node.Lon,
		Lat:       node.Lat,
		Type:      node.Type,
		Props:     node.Props,
		AmenityID: node.AmenityID,
		Name:      node.Name,
		CCode:     node.CCode,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newNode)
	if err != nil {
		return nil, err
	}

	// change user stat
	// var operations []mongo.WriteModel
	// operationA := mongo.NewUpdateOneModel()
	// operationA.SetFilter(bson.M{"_id": userIDPrimitive})
	// operationA.SetUpdate(bson.D{
	// 	{"$inc", bson.D{
	// 		{"user_stat.node", 1},
	// 	}},
	// })
	// operations = append(operations, operationA)
	// _, err = r.db.Collection(TblNode).BulkWrite(ctx, operations)
	// _, _ = r.db.Collection(tblUsers).UpdateOne(ctx, bson.M{"_id": userIDPrimitive}, bson.D{
	// 	{"$inc", bson.D{
	// 		{"user_stat.node", 1},
	// 	}},
	// })

	err = r.db.Collection(TblNode).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *NodeMongo) UpdateNode(id string, userID string, data *model.Node) (*model.Node, error) {
	var result *model.Node
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblNode)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	// idUser, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	return result, err
	// }
	filter := bson.M{"_id": idPrimitive}

	// Find old data
	var oldResult *model.Node
	err = collection.FindOne(ctx, filter).Decode(&oldResult)
	if err != nil {
		return result, err
	}
	oldNodeAudit := model.NodeInput{
		UserID:    oldResult.UserID,
		Lon:       oldResult.Lon,
		Lat:       oldResult.Lat,
		Type:      oldResult.Type,
		Name:      oldResult.Name,
		OsmID:     oldResult.OsmID,
		AmenityID: oldResult.AmenityID,
		Props:     oldResult.Props,
	}
	_, err = r.db.Collection(TblNodeHistory).InsertOne(ctx, oldNodeAudit)
	if err != nil {
		return result, err
	}

	newData := bson.M{}
	if data.Lon != 0 {
		newData["lon"] = data.Lon
	}
	if data.Lat != 0 {
		newData["lat"] = data.Lat
	}
	if data.Type != "" {
		newData["type"] = data.Type
	}
	if data.OsmID != "" {
		newData["osm_id"] = data.OsmID
	}
	if !data.AmenityID.IsZero() {
		newData["amenity_id"] = data.AmenityID
	}
	if data.Props != nil {
		//newProps := make(map[string]interface{})
		newProps := data.Props
		if val, ok := data.Props["status"]; ok {
			if val == -1.0 {
				newDel := make(map[string]interface{})
				newDel["user_id"] = userID
				newDel["del_at"] = time.Now()
				newProps["del"] = newDel
			}
		}
		newData["props"] = newProps
	}
	if data.Name != "" {
		newData["name"] = data.Name
	}
	if data.CCode != "" {
		newData["ccode"] = data.CCode
	}
	// if data.Status != 0 {
	// 	newData["status"] = data.Status
	// }
	newData["updated_at"] = time.Now()
	// bson.M{
	// 	"lon":        data.Lon,
	// 	"lat":        data.Lat,
	// 	"type":       data.Type,
	// 	"osm_id":     data.OsmID,
	// 	"amenity_id": data.AmenityID,
	// 	"props":      data.Props,
	// 	"name":       data.Name,
	// 	"updated_at": time.Now(),
	// }
	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	// update type for nodedata collection
	// if data.Type != "" {
	// _, err = r.db.Collection(TblNodedata).UpdateMany(
	// 	ctx,
	// 	bson.M{"node_id": result.ID},
	// 	bson.M{"$set": bson.M{"type": result.Type, "lat": result.Lat, "lon": result.Lon}},
	// )
	// if err != nil {
	// 	return result, err
	// }
	// }

	return result, nil
}

func (r *NodeMongo) DeleteNode(id string) (model.Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.Node{}
	collection := r.db.Collection(TblNode)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return result, err
	}

	// change user stat
	// _, _ = r.db.Collection(tblUsers).UpdateOne(ctx, bson.M{"_id": result.UserID}, bson.D{
	// 	{"$inc", bson.D{
	// 		{"user_stat.node", -1},
	// 	}},
	// })

	// // remove nodedata.
	// _, err = r.db.Collection(TblNodedata).DeleteMany(ctx, bson.M{"node_id": idPrimitive})
	// if err != nil {
	// 	return result, err
	// }

	// // remove reviews.
	// _, err = r.db.Collection(TblReview).DeleteMany(ctx, bson.M{"node_id": idPrimitive})
	// if err != nil {
	// 	return result, err
	// }

	// // remove audits.
	// _, err = r.db.Collection(TblNodeAudit).DeleteMany(ctx, bson.M{"node_id": idPrimitive})
	// if err != nil {
	// 	return result, err
	// }

	return result, nil
}

func (r *NodeMongo) FindForKml(params domain.RequestParams) (domain.Response[domain.Kml], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Kml
	var response domain.Response[domain.Kml]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.Node]{}, err
	// }
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Kml]{}, err
	}
	fmt.Println("params lang=", params.Lang)
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from": "nodedata",
	// 	// "let":  bson.D{{Key: "nodeId", Value: bson.D{{"$toString", "$_id"}}}},
	// 	// "pipeline": mongo.Pipeline{
	// 	// 	bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$node_id", "$$nodeId"}}}}},
	// 	// },
	// 	"localField":   "_id",
	// 	"foreignField": "node_id",
	// 	"as":           "data",
	// }}})
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "nodedata",
		"as":   "data",
		"let":  bson.D{{Key: "nodeId", Value: "$_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$node_id", "$$nodeId"}}}}},
			// bson.D{{"$limit", 1}},
			bson.D{{
				Key: "$lookup",
				Value: bson.M{
					"from": "tag",
					"as":   "tagx",
					"let":  bson.D{{Key: "tagId", Value: "$tag_id"}},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$tagId"}}}}},
						bson.D{{"$limit", 1}},
						bson.D{{
							Key: "$replaceWith", Value: bson.M{
								"$mergeObjects": bson.A{
									"$$ROOT",
									bson.D{{
										Key: "$ifNull", Value: bson.A{
											fmt.Sprintf("$locale.%s", params.Lang),
											fmt.Sprintf("$locale.%s", r.i18n.Default),
										},
									}},
								},
							},
						}},
					},
				},
			}},
			bson.D{{Key: "$set", Value: bson.M{"tag": bson.M{"$first": "$tagx"}}}},
			bson.D{{
				Key: "$lookup",
				Value: bson.M{
					"from": "users",
					"as":   "userx",
					"let":  bson.D{{Key: "userId", Value: "$user_id"}},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
						bson.D{{"$limit", 1}},
					},
				},
			}},
			bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$userx"}}}},
		},
	}}})

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "users",
		"as":   "userb",
		"let":  bson.D{{Key: "userId", Value: "$user_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
			bson.D{{"$limit", 1}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$userb"}}}})

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "amenity",
		"as":   "amenityx",
		"let":  bson.D{{Key: "amenityId", Value: "$amenity_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$amenityId"}}}}},
			bson.D{{"$limit", 1}},
			bson.D{{
				Key: "$replaceWith", Value: bson.M{
					"$mergeObjects": bson.A{
						"$$ROOT",
						bson.D{{
							Key: "$ifNull", Value: bson.A{
								fmt.Sprintf("$locale.%s", params.Lang),
								fmt.Sprintf("$locale.%s", r.i18n.Default),
							},
						}},
					},
				},
			}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"amenity": bson.M{"$first": "$amenityx"}}}})

	cursor, err := r.db.Collection(TblNode).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Kml, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblNode).CountDocuments(ctx, params.Filter)
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Kml]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *NodeMongo) CreateFile(params domain.RequestParams) (domain.Response[domain.NodeFileItem], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results domain.Response[domain.NodeFileItem]

	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}
	// fmt.Println("params lang=", params.Lang)

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

	var allItems []domain.NodeFileItem
	cursor, err := r.db.Collection(TblNode).Aggregate(ctx, pipe) //Aggregate(ctx, pipe) // Find(ctx, q, options)
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

	data := make([]domain.NodeFileItem, len(allItems))
	for i, _ := range allItems {

		data[i] = allItems[i]
	}
	// fmt.Println("Find total: ", len(data))
	total := len(data)

	results = domain.Response[domain.NodeFileItem]{
		Total: total,
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  data,
	}
	return results, nil
}
