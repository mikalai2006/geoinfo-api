package repository

import (
	"context"
	"time"

	"github.com/mikalai2006/geoinfo-api/internal/config"
	"github.com/mikalai2006/geoinfo-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewReviewMongo(db *mongo.Database, i18n config.I18nConfig) *ReviewMongo {
	return &ReviewMongo{db: db, i18n: i18n}
}

func (r *ReviewMongo) FindReview(params domain.RequestParams) (domain.Response[domain.Review], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Review
	var response domain.Response[domain.Review]
	filter, opts, err := CreateFilterAndOptions(params)
	if err != nil {
		return domain.Response[domain.Review]{}, err
	}

	cursor, err := r.db.Collection(tblReview).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Review, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(tblReview).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Review]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ReviewMongo) GetAllReview(params domain.RequestParams) (domain.Response[domain.Review], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Review
	var response domain.Response[domain.Review]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Review]{}, err
	}

	cursor, err := r.db.Collection(tblReview).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Review, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(tblReview).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Review]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ReviewMongo) CreateReview(userID string, review *domain.Review) (*domain.Review, error) {
	var result *domain.Review

	collection := r.db.Collection(tblReview)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newReview := domain.Review{
		Review:    review.Review,
		Rate:      review.Rate,
		GeoID:     review.GeoID,
		UserID:    userIDPrimitive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newReview)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblReview).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
