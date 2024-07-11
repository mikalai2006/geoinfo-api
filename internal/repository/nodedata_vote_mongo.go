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

type NodedataVoteMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewNodedataVoteMongo(db *mongo.Database, i18n config.I18nConfig) *NodedataVoteMongo {
	return &NodedataVoteMongo{db: db, i18n: i18n}
}

func (r *NodedataVoteMongo) GetNodedataVote(id string) (*model.NodedataVote, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result *model.NodedataVote
	var pipe mongo.Pipeline

	IDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	pipe = append(pipe, bson.D{{"$match", bson.M{"_id": IDPrimitive}}})
	pipe = append(pipe, bson.D{{"$limit", 1}})

	cursor, err := r.db.Collection(TblNodedataVote).Aggregate(ctx, pipe)
	// fmt.Println("filter NodedataVote:::", pipe)
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

func (r *NodedataVoteMongo) FindNodedataVote(params domain.RequestParams) (domain.Response[model.NodedataVote], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.NodedataVote
	var response domain.Response[model.NodedataVote]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.NodedataVote]{}, err
	// }
	// cursor, err := r.db.Collection(TblNodedataVote).Find(ctx, filter, opts)
	pipe, err := CreatePipeline(params, &r.i18n)

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

	// get owner nodedata.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "users",
		"as":   "usero",
		"let":  bson.D{{Key: "nodedataUserId", Value: "$nodedata_user_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$nodedataUserId"}}}}},
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
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"owner": bson.M{"$first": "$usero"}}}})

	if err != nil {
		return response, err
	}

	cursor, err := r.db.Collection(TblNodedataVote).Aggregate(ctx, pipe)
	// fmt.Println("filter NodedataVote:::", pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.NodedataVote, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblNodedataVote).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.NodedataVote]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *NodedataVoteMongo) GetAllNodedataVote(params domain.RequestParams) (domain.Response[model.NodedataVote], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.NodedataVote
	var response domain.Response[model.NodedataVote]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.NodedataVote]{}, err
	}

	cursor, err := r.db.Collection(TblNodedataVote).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.NodedataVote, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblNodedataVote).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.NodedataVote]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *NodedataVoteMongo) CreateNodedataVote(userID string, data *model.NodedataVoteInput) (*model.NodedataVote, error) {
	var result *model.NodedataVote

	collection := r.db.Collection(TblNodedataVote)

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

	// // fmt.Println(bson.M{"nodedata_id": nodedataIDPrimitive, "user_id": userIDPrimitive})

	// var existVote model.NodedataVote
	// r.db.Collection(TblNodedataVote).FindOne(ctx, bson.M{"nodedata_id": nodedataIDPrimitive, "user_id": userIDPrimitive}).Decode(&existVote)
	// // if err != nil {
	// // 	if errors.Is(err, mongo.ErrNoDocuments) {
	// // 		return result, model.ErrAddressNotFound
	// // 	}
	// // 	return nil, err
	// // }
	// var existNodedata model.Nodedata
	// r.db.Collection(TblNodedata).FindOne(ctx, bson.M{"_id": nodedataIDPrimitive}).Decode(&existNodedata)

	// if (existVote == model.NodedataVote{}) {
	newNodedataVote := model.NodedataVoteMongo{
		UserID:         userIDPrimitive,
		NodedataID:     nodedataIDPrimitive,
		NodedataUserID: data.NodedataUserID,
		NodeID:         data.NodeID,
		Value:          data.Value,
		// Status:     100, //data.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newNodedataVote)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblNodedataVote).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	// 	// // set stat active user
	// 	// statFragment := bson.D{}
	// 	// if result.Value > 0 {
	// 	// 	statFragment = append(statFragment, bson.E{"user_stat.nodedataLike", 1})
	// 	// } else {
	// 	// 	statFragment = append(statFragment, bson.E{"user_stat.nodedataDLike", 1})
	// 	// }
	// 	// _, _ = r.db.Collection(tblUsers).UpdateOne(ctx, bson.M{"_id": result.UserID}, bson.D{
	// 	// 	{"$inc", statFragment},
	// 	// })

	// 	// // set stat author user
	// 	// if !existNodedata.UserID.IsZero() {
	// 	// 	statFragment := bson.D{}
	// 	// 	if result.Value > 0 {
	// 	// 		statFragment = append(statFragment, bson.E{"user_stat.nodedataAuthorLike", 1})
	// 	// 	} else {
	// 	// 		statFragment = append(statFragment, bson.E{"user_stat.nodedataAuthorDLike", 1})
	// 	// 	}
	// 	// 	_, _ = r.db.Collection(tblUsers).UpdateOne(ctx, bson.M{"_id": existNodedata.UserID}, bson.D{
	// 	// 		{"$inc", statFragment},
	// 	// 	})
	// 	// }
	// } else {
	// 	updateNodedataVote := &model.NodedataVote{
	// 		// UserID:     userIDPrimitive,
	// 		// NodedataID: nodedataIDPrimitive,
	// 		Value: data.Value,
	// 	}
	// 	result, err = r.UpdateNodedataVote(existVote.ID.Hex(), userID, updateNodedataVote)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	// // set stat active user
	// 	// statFragment := bson.D{}
	// 	// if result.Value > 0 {
	// 	// 	statFragment = append(statFragment, bson.E{"user_stat.nodedataLike", 1})
	// 	// 	statFragment = append(statFragment, bson.E{"user_stat.nodedataDLike", -1})
	// 	// } else {
	// 	// 	statFragment = append(statFragment, bson.E{"user_stat.nodedataLike", -1})
	// 	// 	statFragment = append(statFragment, bson.E{"user_stat.nodedataDLike", 1})
	// 	// }
	// 	// _, _ = r.db.Collection(tblUsers).UpdateOne(ctx, bson.M{"_id": result.UserID}, bson.D{
	// 	// 	{"$inc", statFragment},
	// 	// })

	// 	// // set stat author user
	// 	// if !existNodedata.UserID.IsZero() {
	// 	// 	statFragment := bson.D{}
	// 	// 	if result.Value > 0 {
	// 	// 		statFragment = append(statFragment, bson.E{"user_stat.nodedataAuthorLike", 1})
	// 	// 		statFragment = append(statFragment, bson.E{"user_stat.nodedataAuthorDLike", -1})
	// 	// 	} else {
	// 	// 		statFragment = append(statFragment, bson.E{"user_stat.nodedataAuthorLike", -1})
	// 	// 		statFragment = append(statFragment, bson.E{"user_stat.nodedataAuthorDLike", 1})
	// 	// 	}
	// 	// 	_, _ = r.db.Collection(tblUsers).UpdateOne(ctx, bson.M{"_id": existNodedata.UserID}, bson.D{
	// 	// 		{"$inc", statFragment},
	// 	// 	})
	// 	// }
	// }

	return result, nil
}

func (r *NodedataVoteMongo) GqlGetNodedataVote(params domain.RequestParams) ([]*model.NodedataVote, error) {
	fmt.Println("GqlGetNodedataVotes", &r.i18n, params.Lang)
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.NodedataVote
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}
	// fmt.Println(pipe)

	cursor, err := r.db.Collection(TblNodedataVote).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.NodedataVote, len(results))

	copy(resultSlice, results)
	return results, nil
}

func (r *NodedataVoteMongo) UpdateNodedataVote(id string, userID string, data *model.NodedataVote) (*model.NodedataVote, error) {
	var result *model.NodedataVote
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblNodedataVote)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if data.Value != 0 {
		newData["value"] = data.Value
	}
	// if data.Status != 0 {
	// 	newData["status"] = data.Status
	// }
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

func (r *NodedataVoteMongo) DeleteNodedataVote(id string) (model.NodedataVote, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.NodedataVote{}
	collection := r.db.Collection(TblNodedataVote)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	// // change user stat
	// // countValue := 0
	// // if result.Value > 0 {
	// // 	countValue = -5
	// // } else {
	// // 	countValue = 5
	// // }
	// _, _ = r.db.Collection(tblUsers).UpdateOne(ctx, bson.M{"_id": result.UserID}, bson.D{
	// 	{"$inc", bson.D{
	// 		{"user_stat.countTagLike", -1},
	// 	}},
	// })

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return result, err
	}

	return result, nil
}
