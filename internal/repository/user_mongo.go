package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/config"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"github.com/mikalai2006/geoinfo-api/internal/utils"
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

	// // add stat user tag vote.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": TblNodedataVote,
	// 		"as":   "tests",
	// 		"let":  bson.D{{Key: "userId", Value: "$_id"}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$user_id", "$$userId"}}}}},
	// 			bson.D{
	// 				{
	// 					"$group", bson.D{
	// 						{
	// 							"_id", "",
	// 						},
	// 						{"valueTagLike", bson.D{{"$sum", "$value"}}},
	// 						{"countTagLike", bson.D{{"$sum", 1}}},
	// 					},
	// 				},
	// 			},
	// 			bson.D{{Key: "$project", Value: bson.M{"_id": 0, "valueTagLike": "$valueTagLike", "countTagLike": "$countTagLike"}}},
	// 		},
	// 	},
	// }})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"test": bson.M{"$first": "$tests"}}}})

	// // add stat user node votes.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": TblNodeVote,
	// 		"as":   "tests2",
	// 		"let":  bson.D{{Key: "userId", Value: "$_id"}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$user_id", "$$userId"}}}}},
	// 			bson.D{
	// 				{
	// 					"$group", bson.D{
	// 						{
	// 							"_id", "",
	// 						},
	// 						{"valueNodeLike", bson.D{{"$sum", "$value"}}},
	// 						{"countNodeLike", bson.D{{"$sum", 1}}},
	// 					},
	// 				},
	// 			},
	// 			bson.D{{Key: "$project", Value: bson.M{"_id": 0, "valueNodeLike": "$valueNodeLike", "countNodeLike": "$countNodeLike"}}},
	// 		},
	// 	},
	// }})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"test": bson.D{{
	// 	"$mergeObjects", bson.A{
	// 		"$test",
	// 		bson.M{"$first": "$tests2"},
	// 	},
	// }},
	// }}})

	// // add stat user node votes.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": TblNode,
	// 		"as":   "countNodes",
	// 		"let":  bson.D{{Key: "userId", Value: "$_id"}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$user_id", "$$userId"}}}}},
	// 		},
	// 	},
	// }})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"test": bson.D{{
	// 	"$mergeObjects", bson.A{
	// 		"$test",
	// 		bson.M{"countNodes": bson.M{"$size": "$countNodes"}},
	// 	},
	// }},
	// }}})

	// // add stat user added nodedata.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": TblNodedata,
	// 		"as":   "countNodedatas",
	// 		"let":  bson.D{{Key: "userId", Value: "$_id"}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$user_id", "$$userId"}}}}},
	// 		},
	// 	},
	// }})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"test": bson.D{{
	// 	"$mergeObjects", bson.A{
	// 		"$test",
	// 		bson.M{"countNodedatas": bson.M{"$size": "$countNodedatas"}},
	// 	},
	// }},
	// }}})

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

func (r *UserMongo) SetStat(userID string, inputData model.UserStat) (model.User, error) {
	var result model.User
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblUsers)

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	newData := bson.M{}
	if inputData.Node != 0 {
		newData["user_stat.node"] = utils.Max(result.UserStat.Node+inputData.Node, 0)
	}
	if inputData.NodeAuthorDLike != 0 {
		newData["user_stat.nodeAuthorDLike"] = utils.Max(result.UserStat.NodeAuthorDLike+inputData.NodeAuthorDLike, 0)
	}
	if inputData.NodeAuthorLike != 0 {
		newData["user_stat.nodeAuthorLike"] = utils.Max(result.UserStat.NodeAuthorLike+inputData.NodeAuthorLike, 0)
	}
	if inputData.NodeDLike != 0 {
		newData["user_stat.nodeDLike"] = utils.Max(result.UserStat.NodeDLike+inputData.NodeDLike, 0)
	}
	if inputData.NodeLike != 0 {
		newData["user_stat.nodeLike"] = utils.Max(result.UserStat.NodeLike+inputData.NodeLike, 0)
	}
	if inputData.Nodedata != 0 {
		newData["user_stat.nodedata"] = utils.Max(result.UserStat.Nodedata+inputData.Nodedata, 0)
	}
	if inputData.NodedataAuthorDLike != 0 {
		newData["user_stat.nodedataAuthorDLike"] = utils.Max(result.UserStat.NodedataAuthorDLike+inputData.NodedataAuthorDLike, 0)
	}
	if inputData.NodedataAuthorLike != 0 {
		newData["user_stat.nodedataAuthorLike"] = utils.Max(result.UserStat.NodedataAuthorLike+inputData.NodedataAuthorLike, 0)
	}
	if inputData.NodedataDLike != 0 {
		newData["user_stat.nodedataDLike"] = utils.Max(result.UserStat.NodedataDLike+inputData.NodedataDLike, 0)
	}
	if inputData.NodedataLike != 0 {
		newData["user_stat.nodedataLike"] = utils.Max(result.UserStat.NodedataLike+inputData.NodedataLike, 0)
	}
	if inputData.Review != 0 {
		newData["user_stat.review"] = utils.Max(result.UserStat.Review+inputData.Review, 0)
	}

	// fmt.Println("newData=", newData)
	err = collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": newData}).Decode(&result)
	if err != nil {
		return result, err
	}

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

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
