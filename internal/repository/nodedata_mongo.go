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

type NodedataMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewNodedataMongo(db *mongo.Database, i18n config.I18nConfig) *NodedataMongo {
	return &NodedataMongo{db: db, i18n: i18n}
}

func (r *NodedataMongo) GetNodedata(id string) (*model.Nodedata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result *model.Nodedata
	var pipe mongo.Pipeline

	IDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	pipe = append(pipe, bson.D{{"$match", bson.M{"_id": IDPrimitive}}})
	pipe = append(pipe, bson.D{{"$limit", 1}})

	cursor, err := r.db.Collection(TblNodedata).Aggregate(ctx, pipe)
	// fmt.Println("filter Nodedata:::", pipe)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if er := cursor.Decode(&result); er != nil {
			return result, er
		}
	}

	return result, nil
}

func (r *NodedataMongo) FindNodedata(params domain.RequestParams) (domain.Response[model.Nodedata], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Nodedata
	var response domain.Response[model.Nodedata]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.Nodedata]{}, err
	// }
	// cursor, err := r.db.Collection(TblNodedata).Find(ctx, filter, opts)
	pipe, err := CreatePipeline(params, &r.i18n)

	if err != nil {
		return response, err
	}

	cursor, err := r.db.Collection(TblNodedata).Aggregate(ctx, pipe)
	// fmt.Println("filter Nodedata:::", pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Nodedata, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblNodedata).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Nodedata]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *NodedataMongo) GetAllNodedata(params domain.RequestParams) (domain.Response[model.Nodedata], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Nodedata
	var response domain.Response[model.Nodedata]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.Nodedata]{}, err
	}

	cursor, err := r.db.Collection(TblNodedata).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Nodedata, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblNodedata).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Nodedata]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *NodedataMongo) CreateNodedata(userID string, data *model.NodedataInput) (*model.Nodedata, error) {
	var result *model.Nodedata

	collection := r.db.Collection(TblNodedata)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	tagIDPrimitive, err := primitive.ObjectIDFromHex(data.TagID)
	if err != nil {
		return nil, err
	}
	tagoptIDPrimitive := primitive.NilObjectID
	if data.TagoptID != "" {
		tagoptIDPrimitive, err = primitive.ObjectIDFromHex(data.TagoptID)
		if err != nil {
			return nil, err
		}
	}
	nodeIDPrimitive, err := primitive.ObjectIDFromHex(data.NodeID)
	if err != nil {
		return nil, err
	}

	newNodedata := model.NodedataInputMongo{
		UserID:   userIDPrimitive,
		NodeID:   nodeIDPrimitive,
		TagID:    tagIDPrimitive,
		TagoptID: tagoptIDPrimitive,
		Data:     data.Data,
		Locale:   data.Locale,
		Status:   100, //data.Status,
		// Type:      data.Type,
		Like:      0,
		Dlike:     0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newNodedata)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblNodedata).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *NodedataMongo) GqlGetNodedatas(params domain.RequestParams) ([]*model.Nodedata, error) {
	fmt.Println("GqlGetNodedatas", &r.i18n, params.Lang)
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.Nodedata
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}
	// fmt.Println(pipe)
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
		"from": "nodedata_vote",
		"as":   "votes",
		"let":  bson.D{{Key: "id", Value: "$_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$nodedata_id", "$$id"}}}}},
		},
	}}})

	cursor, err := r.db.Collection(TblNodedata).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.Nodedata, len(results))

	copy(resultSlice, results)
	return results, nil
}

func (r *NodedataMongo) UpdateNodedata(id string, userID string, data *model.Nodedata) (*model.Nodedata, error) {
	var result *model.Nodedata
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblNodedata)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if !data.TagID.IsZero() {
		newData["tag_id"] = data.TagID
	}
	if !data.TagoptID.IsZero() {
		newData["tagopt_id"] = data.TagoptID
	}
	if !data.NodeID.IsZero() {
		newData["node_id"] = data.NodeID
	}
	if data.Data != (model.NodedataData{}) {
		newData["data"] = data.Data
	}
	if data.Locale != nil {
		newData["locale"] = data.Locale
	}
	if data.Status != 0 {
		newData["status"] = data.Status
	}
	newData["like"] = data.Like
	newData["dlike"] = data.Dlike
	newData["updated_at"] = time.Now()

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *NodedataMongo) DeleteNodedata(id string) (model.Nodedata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.Nodedata{}
	collection := r.db.Collection(TblNodedata)

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

	return result, nil
}

func (r *NodedataMongo) AddAudit(userID string, data *model.NodedataAuditInput) (*model.Nodedata, error) {
	var result *model.Nodedata

	collection := r.db.Collection(TblNodedataAudit)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	nodedataIDPrimitive, err := primitive.ObjectIDFromHex(data.NodedataID)
	if err != nil {
		return nil, err
	}

	newNodedataAudit := model.NodedataAuditDB{
		UserID:     userIDPrimitive,
		NodedataID: nodedataIDPrimitive,
		Value:      data.Value,
		Props:      data.Props,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = collection.InsertOne(ctx, newNodedataAudit)
	if err != nil {
		return nil, err
	}

	result, err = r.GetNodedata(data.NodedataID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *NodedataMongo) FindAudits(params domain.RequestParams) (domain.Response[model.NodedataAudit], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.NodedataAudit
	var response domain.Response[model.NodedataAudit]

	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return response, err
	}
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

	// if params.Sort != 0 {
	// 	pipe = append(pipe, bson.D{{"$sort", params.Sort}})
	// }
	// if params.Skip != 0 {
	// 	pipe = append(pipe, bson.D{{"$skip", params.Skip}})
	// }
	// if params.Limit != 0 {
	// 	pipe = append(pipe, bson.D{{"$limit", params.Limit}})
	// }

	cursor, err := r.db.Collection(TblNodedataAudit).Aggregate(ctx, pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.NodedataAudit, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblNodedataAudit).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.NodedataAudit]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}
