package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBHandler interface {
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult
}

type TokenDBHandler interface {
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	FindOne(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult
}
