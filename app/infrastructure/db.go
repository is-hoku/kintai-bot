package infrastructure

import (
	"context"
	"log"

	"kintai-bot/app/interfaces/database"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBHandler struct {
	Coll *mongo.Collection
}

func NewDBHandler() database.DBHandler {
	clientOptions := options.Client().ApplyURI("mongodb://db")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("kintai").Collection("users")
	return &DBHandler{collection}
}

func (handler *DBHandler) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return handler.Coll.InsertOne(context.TODO(), document, opts...)
}

func (handler *DBHandler) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return handler.Coll.FindOne(context.TODO(), filter, opts...)
}
