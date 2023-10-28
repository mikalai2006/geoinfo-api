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

type TicketMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewTicketMongo(db *mongo.Database, i18n config.I18nConfig) *TicketMongo {
	return &TicketMongo{db: db, i18n: i18n}
}

func (r *TicketMongo) FindTicket(params domain.RequestParams) (domain.Response[model.Ticket], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Ticket
	var response domain.Response[model.Ticket]
	filter, opts, err := CreateFilterAndOptions(params)
	if err != nil {
		return domain.Response[model.Ticket]{}, err
	}

	cursor, err := r.db.Collection(TblTicket).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Ticket, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblTicket).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Ticket]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *TicketMongo) GetAllTicket(params domain.RequestParams) (domain.Response[model.Ticket], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Ticket
	var response domain.Response[model.Ticket]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.Ticket]{}, err
	}

	cursor, err := r.db.Collection(TblTicket).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Ticket, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblTicket).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Ticket]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *TicketMongo) CreateTicket(userID string, Ticket *model.Ticket) (*model.Ticket, error) {
	var result *model.Ticket

	collection := r.db.Collection(TblTicket)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newTicket := model.Ticket{
		UserID:      userIDPrimitive,
		Title:       Ticket.Title,
		Status:      Ticket.Status,
		Progress:    Ticket.Progress,
		Description: Ticket.Description,
		Props:       Ticket.Props,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	res, err := collection.InsertOne(ctx, newTicket)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblTicket).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TicketMongo) GqlGetTickets(params domain.RequestParams) ([]*model.Ticket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.Ticket
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}

	cursor, err := r.db.Collection(TblTicket).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.Ticket, len(results))

	copy(resultSlice, results)
	return results, nil
}

func (r *TicketMongo) DeleteTicket(id string) (model.Ticket, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.Ticket{}
	collection := r.db.Collection(TblTicket)

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
