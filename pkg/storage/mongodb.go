package storage

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBManager struct {
	*mongo.Client
}

func NewMongoDBManager(ctx context.Context) (*MongoDBManager, error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf(
		"mongodb:%s ",
		viper.GetString("mongodb.uri"),
	))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	return &MongoDBManager{
		Client: client,
	}, nil
}
