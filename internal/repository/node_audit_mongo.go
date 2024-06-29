package repository

import (
	"context"
	"time"

	"github.com/mikalai2006/geoinfo-api/graph/model"
	"github.com/mikalai2006/geoinfo-api/internal/config"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NodeAuditMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewNodeAuditMongo(db *mongo.Database, i18n config.I18nConfig) *NodeAuditMongo {
	return &NodeAuditMongo{db: db, i18n: i18n}
}

func (r *NodeAuditMongo) FindNodeAudit(params domain.RequestParams) (domain.Response[model.NodeAudit], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.NodeAudit
	var response domain.Response[model.NodeAudit]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.Node]{}, err
	// }
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.NodeAudit]{}, err
	}

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "users",
		"as":   "usera",
		"let":  bson.D{{Key: "userId", Value: "$user_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
			bson.D{{"$limit", 1}},
			bson.D{{
				Key: "$lookup",
				Value: bson.M{
					"from": tblImage,
					"as":   "images",
					"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
					},
				},
			}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$usera"}}}})

	cursor, err := r.db.Collection(TblNodeAudit).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.NodeAudit, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblNode).CountDocuments(ctx, params.Filter)
	if err != nil {
		return response, err
	}

	response = domain.Response[model.NodeAudit]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *NodeAuditMongo) CreateNodeAudit(userID string, nodeAudit *model.NodeAuditInput) (*model.NodeAudit, error) {
	var result *model.NodeAudit

	collection := r.db.Collection(TblNodeAudit)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	createdAt := nodeAudit.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	newNodeAudit := model.NodeAuditInput{
		UserID:    userIDPrimitive,
		NodeID:    nodeAudit.NodeID,
		Status:    1,
		Message:   nodeAudit.Message,
		Props:     nodeAudit.Props,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newNodeAudit)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblNodeAudit).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *NodeAuditMongo) UpdateNodeAudit(id string, userID string, data *model.NodeAuditInput) (*model.NodeAudit, error) {
	var result *model.NodeAudit
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblNodeAudit)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	// idUser, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	return result, err
	// }
	filter := bson.M{"_id": idPrimitive}

	// // Find old data
	// var oldResult *model.NodeAudit
	// err = collection.FindOne(ctx, filter).Decode(&oldResult)
	// if err != nil {
	// 	return result, err
	// }
	// oldNodeAudit := model.NodeAudit{
	// 	UserID:  oldResult.UserID,
	// 	NodeID:  oldResult.NodeID,
	// 	Message: oldResult.Message,
	// 	Status:  oldResult.Status,
	// 	Props:   oldResult.Props,
	// }
	// _, err = r.db.Collection(TblNodeAudit).UpdateOne(ctx, filter, bson.M{"$set": oldNodeAudit})
	// if err != nil {
	// 	return result, err
	// }

	newData := bson.M{}
	if data.Message != "" {
		newData["message"] = data.Message
	}
	if data.Status != 0 {
		newData["status"] = data.Status
	}
	if data.Props != nil {
		newData["props"] = data.Props
	}
	// if data.Props != nil {
	// 	//newProps := make(map[string]interface{})
	// 	newProps := data.Props
	// 	if val, ok := data.Props["status"]; ok {
	// 		if val == -1.0 {
	// 			newDel := make(map[string]interface{})
	// 			newDel["user_id"] = userID
	// 			newDel["del_at"] = time.Now()
	// 			newProps["del"] = newDel
	// 		}
	// 	}
	// 	newData["props"] = newProps
	// }
	newData["updated_at"] = time.Now()
	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	// err = collection.FindOne(ctx, filter).Decode(&result)
	// if err != nil {
	// 	return result, err
	// }
	resultResponse, err := r.FindNodeAudit(domain.RequestParams{Filter: bson.D{
		{"_id", idPrimitive},
	}})
	if err != nil {
		return result, err
	}

	result = &resultResponse.Data[0]

	return result, nil
}

func (r *NodeAuditMongo) DeleteNodeAudit(id string) (model.NodeAudit, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.NodeAudit{}
	collection := r.db.Collection(TblNodeAudit)

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
