package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/config"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LikeMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewLikeMongo(db *mongo.Database, i18n config.I18nConfig) *LikeMongo {
	return &LikeMongo{db: db, i18n: i18n}
}

func (r *LikeMongo) FindLike(params domain.RequestParams) (domain.Response[model.Like], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Like
	var response domain.Response[model.Like]
	pipe, err := CreatePipeline(params, &r.i18n)

	if err != nil {
		return response, err
	}

	cursor, err := r.db.Collection(TblLike).Aggregate(ctx, pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Like, len(results))
	copy(resultSlice, results)

	count, err := r.db.Collection(TblLike).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Like]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *LikeMongo) CreateLike(userID string, like *model.LikeInput) (*model.Like, error) {
	var result *model.Like

	collection := r.db.Collection(TblLike)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	nodeIDPrimitive, err := primitive.ObjectIDFromHex(like.NodeID)
	if err != nil {
		return nil, err
	}

	newLike := model.Like{
		UserID:    userIDPrimitive,
		NodeID:    nodeIDPrimitive,
		Status:    like.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newLike)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblLike).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *LikeMongo) GqlGetLikes(params domain.RequestParams) ([]*model.Like, error) {
	fmt.Println("GqlGetLikes")
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.Like
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}
	// fmt.Println(pipe)

	cursor, err := r.db.Collection(TblLike).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.Like, len(results))

	copy(resultSlice, results)
	return results, nil
}

func (r *LikeMongo) GqlGetIamLike(userID string, nodeID string) (*model.Like, error) {
	fmt.Println("GqlGetIamLike")
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result *model.Like

	nodeIDPrimitive, err := primitive.ObjectIDFromHex(nodeID)
	if err != nil {
		return result, err
	}
	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return result, err
	}

	if err := r.db.Collection(TblLike).FindOne(ctx, bson.D{{"node_id", nodeIDPrimitive}, {"user_id", userIDPrimitive}}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, model.ErrLikeNotFound
		}
		return result, err
	}
	return result, nil
}

func (r *LikeMongo) UpdateLike(id string, userID string, data *model.Like) (*model.Like, error) {
	var result *model.Like
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblLike)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
		"status":     data.Status,
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

func (r *LikeMongo) DeleteLike(id string) (model.Like, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.Like{}
	collection := r.db.Collection(TblLike)

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
