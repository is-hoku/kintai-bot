package infrastructure

import (
	"context"
	"log"

	"kintai-bot/app/interfaces/database"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TokenDBHandler struct {
	Coll *mongo.Collection
}

func NewTokenDBHandler() database.TokenDBHandler {
	clientOptions := options.Client().ApplyURI("mongodb://db")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("kintai").Collection("tokens")
	return &TokenDBHandler{collection}
}

func (handler *TokenDBHandler) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return handler.Coll.UpdateOne(context.TODO(), filter, update, opts...)
}

func (handler *TokenDBHandler) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return handler.Coll.FindOne(context.TODO(), filter, opts...)
}
