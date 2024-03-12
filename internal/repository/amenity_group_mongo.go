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

type AmenityGroupMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewAmenityGroupMongo(db *mongo.Database, i18n config.I18nConfig) *AmenityGroupMongo {
	return &AmenityGroupMongo{db: db, i18n: i18n}
}

func (r *AmenityGroupMongo) FindAmenityGroup(params domain.RequestParams) (domain.Response[model.AmenityGroup], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.AmenityGroup
	var response domain.Response[model.AmenityGroup]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.AmenityGroup]{}, err
	// }
	// cursor, err := r.db.Collection(TblAmenityGroup).Find(ctx, filter, opts)
	pipe, err := CreatePipeline(params, &r.i18n)

	if err != nil {
		return response, err
	}

	cursor, err := r.db.Collection(TblAmenityGroup).Aggregate(ctx, pipe)
	// fmt.Println("filter AmenityGroup:::", pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.AmenityGroup, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblAmenityGroup).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.AmenityGroup]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *AmenityGroupMongo) GetAllAmenityGroup(params domain.RequestParams) (domain.Response[model.AmenityGroup], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.AmenityGroup
	var response domain.Response[model.AmenityGroup]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.AmenityGroup]{}, err
	}

	cursor, err := r.db.Collection(TblAmenityGroup).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.AmenityGroup, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblAmenityGroup).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.AmenityGroup]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *AmenityGroupMongo) CreateAmenityGroup(userID string, AmenityGroup *model.AmenityGroup) (*model.AmenityGroup, error) {
	var result *model.AmenityGroup
	collection := r.db.Collection(TblAmenityGroup)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newAmenityGroup := model.AmenityGroup{
		UserID:      userIDPrimitive,
		Title:       AmenityGroup.Title,
		Description: AmenityGroup.Description,
		Props:       AmenityGroup.Props,
		Locale:      AmenityGroup.Locale,
		Status:      AmenityGroup.Status,
		SortOrder:   AmenityGroup.SortOrder,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	res, err := collection.InsertOne(ctx, newAmenityGroup)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblAmenityGroup).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *AmenityGroupMongo) GqlGetAmenityGroups(params domain.RequestParams) ([]*model.AmenityGroup, error) {
	fmt.Println("GqlGetAmenityGroups")
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.AmenityGroup
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}
	// fmt.Println(pipe)

	cursor, err := r.db.Collection(TblAmenityGroup).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.AmenityGroup, len(results))

	copy(resultSlice, results)
	return results, nil
}

func (r *AmenityGroupMongo) UpdateAmenityGroup(id string, userID string, data *model.AmenityGroup) (*model.AmenityGroup, error) {
	var result *model.AmenityGroup
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	fmt.Println("hello UpdateAmenityGroup")
	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	return nil, err
	// }

	// newAmenityGroup := model.AmenityGroup{
	// 	UserID:      userIDPrimitive,
	// 	Key:         data.Key,
	// 	Title:       data.Title,
	// 	Description: data.Description,
	// 	Props:       data.Props,
	// 	Locale:      data.Locale,
	// 	UpdatedAt:   time.Now(),
	// }
	// obj := data.(map[string]interface{})
	// obj["user_id"] = userIDPrimitive
	// data = obj

	collection := r.db.Collection(TblAmenityGroup)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if data.Title != "" {
		newData["title"] = data.Title
	}
	if data.Description != "" {
		newData["description"] = data.Description
	}
	if data.Props != nil || len(data.Props) > 0 {
		newData["props"] = data.Props
	}
	if data.Locale != nil && len(data.Locale) > 0 {
		newData["locale"] = data.Locale
	}
	if data.Status != 0 {
		newData["status"] = data.Status
	}
	if data.SortOrder != 0 {
		newData["sort_order"] = data.SortOrder
	}
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

func (r *AmenityGroupMongo) DeleteAmenityGroup(id string) (model.AmenityGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.AmenityGroup{}
	collection := r.db.Collection(TblAmenityGroup)

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
