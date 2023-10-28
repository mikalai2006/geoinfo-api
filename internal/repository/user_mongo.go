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

type UserMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewUserMongo(db *mongo.Database, i18n config.I18nConfig) *UserMongo {
	return &UserMongo{db: db, i18n: i18n}
}

func (r *UserMongo) Iam(userID string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result model.User
	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return model.User{}, err
	}

	params := domain.RequestParams{}
	params.Filter = bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(tblUsers).FindOne(ctx, params.Filter).Decode(&result)
	if err != nil {
		return model.User{}, err
	}

	pipe, err := CreatePipeline(params, &r.i18n) // mongo.Pipeline{bson.D{{"_id", userIDPrimitive}}} //
	if err != nil {
		return result, err
	}

	// add populate.
	pipe = append(pipe, bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from": tblImage,
			"as":   "images",
			// "localField":   "_id",
			// "foreignField": "service_id",
			"let": bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
			},
		},
	}})

	cursor, err := r.db.Collection(tblUsers).Aggregate(ctx, pipe) // .FindOne(ctx, filter).Decode(&result)
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

func (r *UserMongo) GetUser(id string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result model.User

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.User{}, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(tblUsers).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return model.User{}, err
	}

	pipe, err := CreatePipeline(domain.RequestParams{
		Filter: filter,
	}, &r.i18n)
	if err != nil {
		return result, err
	}

	// add populate.
	pipe = append(pipe, bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from": tblImage,
			"as":   "images",
			// "localField":   "_id",
			// "foreignField": "service_id",
			"let": bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
			},
		},
	}})

	cursor, err := r.db.Collection(tblUsers).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
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

func (r *UserMongo) FindUser(params domain.RequestParams) (domain.Response[model.User], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.User
	var response domain.Response[model.User]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.User]{}, err
	}
	fmt.Println("params:::", params)

	// add populate.
	pipe = append(pipe, bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from": tblImage,
			"as":   "images",
			// "localField":   "_id",
			// "foreignField": "service_id",
			"let": bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
			},
		},
	}})

	cursor, err := r.db.Collection(tblUsers).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.User, len(results))
	copy(resultSlice, results)

	count, err := r.db.Collection(tblUsers).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.User]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *UserMongo) CreateUser(userID string, user *model.User) (*model.User, error) {
	var result *model.User

	collection := r.db.Collection(tblUsers)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newUser := model.User{
		Avatar:    user.Avatar,
		Name:      user.Name,
		UserID:    userIDPrimitive,
		Login:     user.Login,
		Lang:      user.Lang,
		Currency:  user.Currency,
		Online:    user.Online,
		Verify:    user.Verify,
		LastTime:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblUsers).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *UserMongo) DeleteUser(id string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.User{}
	collection := r.db.Collection(tblUsers)

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

func (r *UserMongo) UpdateUser(id string, user *model.User) (model.User, error) {
	var result model.User
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblUsers)

	// data, err := utils.GetBodyToData(user)
	// if err != nil {
	// 	return result, err
	// }

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if user.Lang != "" {
		newData["lang"] = user.Lang
	}
	if user.Name != "" {
		newData["name"] = user.Name
	}
	if user.Login != "" {
		newData["login"] = user.Login
	}

	// fmt.Println("data=", user)
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

func (r *UserMongo) GqlGetUsers(params domain.RequestParams) ([]*model.User, error) {
	fmt.Println("GqlGetUsers")
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.User
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}

	// add populate.
	pipe = append(pipe, bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from": tblImage,
			"as":   "images",
			// "localField":   "_id",
			// "foreignField": "service_id",
			"let": bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
			},
		},
	}})

	// fmt.Println(pipe)

	cursor, err := r.db.Collection(tblUsers).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.User, len(results))

	copy(resultSlice, results)
	return results, nil
}
