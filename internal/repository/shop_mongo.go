package repository

import (
	"context"

	"github.com/mikalai2006/geoinfo-api/internal/config"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShopMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

const (
	collectionName string = "shops"
)

func NewShopMongo(db *mongo.Database, i18n config.I18nConfig) *ShopMongo {
	return &ShopMongo{db: db, i18n: i18n}
}

func (r *ShopMongo) FindShop(params domain.RequestParams) (domain.Response[domain.Shop], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Shop
	var response domain.Response[domain.Shop]
	filter, opts, err := CreateFilterAndOptions(params)
	if err != nil {
		return domain.Response[domain.Shop]{}, err
	}

	cursor, err := r.db.Collection(collectionName).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Shop, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(collectionName).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Shop]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ShopMongo) GetAllShops(params domain.RequestParams) (domain.Response[domain.Shop], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Shop
	var response domain.Response[domain.Shop]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Shop]{}, err
	}

	cursor, err := r.db.Collection(collectionName).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Shop, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(collectionName).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Shop]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ShopMongo) CreateShop(userID string, shop *domain.Shop) (*domain.Shop, error) {
	var result *domain.Shop

	collection := r.db.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	newShop := domain.Shop{
		Title:       shop.Title,
		Description: shop.Description,
		Seo:         "",
		UserID:      userID,
	}

	res, err := collection.InsertOne(ctx, newShop)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(collectionName).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
