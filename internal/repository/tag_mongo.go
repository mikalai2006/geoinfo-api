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

type TagMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewTagMongo(db *mongo.Database, i18n config.I18nConfig) *TagMongo {
	return &TagMongo{db: db, i18n: i18n}
}

func (r *TagMongo) FindTag(params domain.RequestParams) (domain.Response[model.Tag], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Tag
	var response domain.Response[model.Tag]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.Tag]{}, err
	// }
	// cursor, err := r.db.Collection(TblTag).Find(ctx, filter, opts)
	pipe, err := CreatePipeline(params, &r.i18n)

	if err != nil {
		return response, err
	}

	cursor, err := r.db.Collection(TblTag).Aggregate(ctx, pipe)
	// fmt.Println("filter tag:::", pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Tag, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblTag).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Tag]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *TagMongo) GetAllTag(params domain.RequestParams) (domain.Response[model.Tag], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Tag
	var response domain.Response[model.Tag]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.Tag]{}, err
	}

	cursor, err := r.db.Collection(TblTag).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Tag, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblTag).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Tag]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *TagMongo) CreateTag(userID string, tag *model.Tag) (*model.Tag, error) {
	var result *model.Tag

	collection := r.db.Collection(TblTag)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newTag := model.Tag{
		UserID:        userIDPrimitive,
		Key:           tag.Key,
		MultiOpt:      tag.MultiOpt,
		Type:          tag.Type,
		Title:         tag.Title,
		Description:   tag.Description,
		Props:         tag.Props,
		Locale:        tag.Locale,
		Multilanguage: tag.Multilanguage,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	res, err := collection.InsertOne(ctx, newTag)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblTag).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TagMongo) GqlGetTags(params domain.RequestParams) ([]*model.Tag, error) {
	fmt.Println("GqlGetTags")
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.Tag
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}
	// fmt.Println(pipe)

	cursor, err := r.db.Collection(TblTag).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.Tag, len(results))

	copy(resultSlice, results)
	return results, nil
}

func (r *TagMongo) UpdateTag(id string, userID string, data *model.Tag) (*model.Tag, error) {
	var result *model.Tag
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

	collection := r.db.Collection(TblTag)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
		"key":           data.Key,
		"multi_opt":     data.MultiOpt,
		"type":          data.Type,
		"title":         data.Title,
		"description":   data.Description,
		"props":         data.Props,
		"locale":        data.Locale,
		"multilanguage": data.Multilanguage,
		"updated_at":    time.Now(),
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

func (r *TagMongo) DeleteTag(id string) (model.Tag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.Tag{}
	collection := r.db.Collection(TblTag)

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
