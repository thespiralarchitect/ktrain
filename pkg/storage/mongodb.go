package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBManager struct {
	*mongo.Client
	*mongo.Database
	*mongo.Collection
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
	actionDatabase := client.Database(viper.GetString("mongodb.database"))
	actionCollection := actionDatabase.Collection(viper.GetString("mongodb.collection"))
	return &MongoDBManager{
		Client:     client,
		Database:   actionDatabase,
		Collection: actionCollection,
	}, nil
}
func (m *MongoDBManager) Close(ctx context.Context) {
	err := m.Client.Disconnect(ctx)
	if err != nil {
		log.Fatalf("Could not close storage, err: %v", err)
	}
}
