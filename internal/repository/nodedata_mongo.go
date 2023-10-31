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

	newNodedata := model.Nodedata{
		UserID:    userIDPrimitive,
		NodeID:    nodeIDPrimitive,
		TagID:     tagIDPrimitive,
		TagoptID:  tagoptIDPrimitive,
		Data:      data.Data,
		Locale:    data.Locale,
		Status:    100, //data.Status,
		Type:      data.Type,
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

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
		"tag_id":     data.TagID,
		"tagopt_id":  data.TagoptID,
		"node_id":    data.NodeID,
		"data":       data.Data,
		"locale":     data.Locale,
		"status":     data.Status,
		"type":       data.Type,
		"updated_at": time.Now(),
	}})
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
