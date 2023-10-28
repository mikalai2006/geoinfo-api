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

type TagoptMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewTagoptMongo(db *mongo.Database, i18n config.I18nConfig) *TagoptMongo {
	return &TagoptMongo{db: db, i18n: i18n}
}

func (r *TagoptMongo) FindTagopt(params domain.RequestParams) (domain.Response[model.Tagopt], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Tagopt
	var response domain.Response[model.Tagopt]

	// filter := params.Filter.(map[string]interface{})
	// if filter["tag_id"] != nil {
	// 	tagIDPrimitive, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", filter["tag_id"]))
	// 	if err != nil {
	// 		return response, err
	// 	}

	// 	filter["tag_id"] = tagIDPrimitive
	// }
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.Tagopt]{}, err
	// }

	// cursor, err := r.db.Collection(TblTagopt).Find(ctx, filter, opts)
	// if err != nil {
	// 	return response, err
	// }
	// defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &results); er != nil {
	// 	return response, er
	// }

	// if params.Filter["tag_id"] {

	// }

	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return response, err
	}
	cursor, err := r.db.Collection(TblTagopt).Aggregate(ctx, pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Tagopt, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblTagopt).CountDocuments(ctx, params.Filter)
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Tagopt]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *TagoptMongo) GqlGetTagopts(params domain.RequestParams) ([]*model.Tagopt, error) {
	fmt.Println("GqlGetTagopts: ")
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.Tagopt
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}

	cursor, err := r.db.Collection(TblTagopt).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.Tagopt, len(results))

	copy(resultSlice, results)
	return results, nil
}

func (r *TagoptMongo) GetAllTagopt(params domain.RequestParams) (domain.Response[model.Tagopt], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Tagopt
	var response domain.Response[model.Tagopt]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.Tagopt]{}, err
	}

	cursor, err := r.db.Collection(TblTagopt).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Tagopt, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblTagopt).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Tagopt]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *TagoptMongo) CreateTagopt(userID string, tag *model.TagoptInput) (*model.Tagopt, error) {
	var result *model.Tagopt

	collection := r.db.Collection(TblTagopt)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	tagIDPrimitive, err := primitive.ObjectIDFromHex(tag.TagID)
	if err != nil {
		return nil, err
	}

	newTag := model.Tagopt{
		UserID:      userIDPrimitive,
		TagID:       tagIDPrimitive,
		Value:       tag.Value,
		Title:       tag.Title,
		Description: tag.Description,
		Locale:      tag.Locale,
		Props:       tag.Props,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	res, err := collection.InsertOne(ctx, newTag)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblTagopt).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TagoptMongo) UpdateTagopt(id string, userID string, data *model.TagoptInput) (*model.Tagopt, error) {
	var result *model.Tagopt
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	return nil, err
	// }

	// newTag := model.Tag{
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

	collection := r.db.Collection(TblTagopt)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if data.TagID != "" {
		tagIDPrimitive, err := primitive.ObjectIDFromHex(data.TagID)
		if err != nil {
			return result, err
		}
		newData["tag_id"] = tagIDPrimitive
	}
	if data.Value != "" {
		newData["value"] = data.Value
	}
	if data.Title != "" {
		newData["title"] = data.Title
	}
	if data.Description != "" {
		newData["description"] = data.Description
	}
	if data.Props != nil {
		newData["props"] = data.Props
	}
	if data.Locale != nil {
		newData["locale"] = data.Locale
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

func (r *TagoptMongo) DeleteTagopt(id string) (model.Tagopt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.Tagopt{}
	collection := r.db.Collection(TblTagopt)

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
