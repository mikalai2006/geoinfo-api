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

	cursor, err := r.db.Collection(TblReview).Find(ctx, filter, opts)
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

	count, err := r.db.Collection(TblReview).CountDocuments(ctx, bson.M{})
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

	cursor, err := r.db.Collection(TblReview).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
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

	count, err := r.db.Collection(TblReview).CountDocuments(ctx, bson.M{})
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

func (r *ReviewMongo) GqlGetReviews(params domain.RequestParams) ([]*model.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.Review
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}

	cursor, err := r.db.Collection(TblReview).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.Review, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	// count, err := r.db.Collection(TblReview).CountDocuments(ctx, bson.M{})
	// if err != nil {
	// 	return results, err
	// }

	// results = []*domain.Review{
	// 	Total: int(count),
	// 	Skip:  int(params.Options.Skip),
	// 	Limit: int(params.Options.Limit),
	// 	Data:  resultSlice,
	// }
	return results, nil
}

func (r *ReviewMongo) GqlGetCountReviews(params domain.RequestParams) (*model.ReviewInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results model.ReviewInfo
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return &results, err
	}

	var allItems []*model.Review
	cursor, err := r.db.Collection(TblReview).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return &results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &allItems); er != nil {
		return &results, er
	}

	count := len(allItems)
	results.Count = &count

	var sum = int(0)
	for _, t := range allItems {
		sum += t.Rate
	}
	results.Value = &sum

	pipe = append(pipe,
		bson.D{
			{"$group", bson.D{
				{"_id", "$rate"},
				// {"average_price", bson.D{{"$avg", "$price"}}},
				{"count", bson.D{{"$sum", 1}}},
			}}})

	var allItemsGroup []*map[string]int
	cursorGroup, err := r.db.Collection(TblReview).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return &results, err
	}
	defer cursorGroup.Close(ctx)

	if er := cursorGroup.All(ctx, &allItemsGroup); er != nil {
		return &results, er
	}

	results.Ratings = allItemsGroup

	// if typeC, ok := f["type"]; ok {
	// 	typeContent = typeC.(string)
	// }
	// keys := make([]int, 0, len(allItemsGroup))
	// counts := make([]int, 0, len(allItemsGroup))
	// for k := range allItemsGroup {
	// 	f := *allItemsGroup[k]
	// 	counts = append(counts, f["count"])
	// 	keys = append(keys, f["_id"])
	// }
	// sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	// for _, v := range keys {
	// 	f := *allItemsGroup[counts[v-1]]
	// 	fmt.Println("allItemsGroup:::", v, f["_id"], "-", f["count"])
	// }

	return &results, nil
}

func (r *ReviewMongo) CreateReview(userID string, review *domain.Review) (*domain.Review, error) {
	var result *domain.Review

	collection := r.db.Collection(TblReview)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newReview := domain.Review{
		Review:    review.Review,
		Rate:      review.Rate,
		OsmID:     review.OsmID,
		UserID:    userIDPrimitive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newReview)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblReview).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
