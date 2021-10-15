package storage

import (
	"context"
	"fmt"
	"ktrain/pkg/logger"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBManager struct {
	*mongo.Client
	*mongo.Database
}

func NewMongoDBManager(ctx context.Context) (*MongoDBManager, error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf(
		"mongodb:%s",
		viper.GetString("mongodb.uri"),
	))
	// if viper.GetBool("mongodb.hasAuth") {
	// 	clientOptions.SetAuth(options.Credential{
	// 		Username: viper.GetString("mongodb.username"),
	// 		Password: viper.GetString("mongodb.password"),
	// 	})
	// }

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	actionDatabase := client.Database(viper.GetString("mongodb.database"))
	return &MongoDBManager{
		Client:   client,
		Database: actionDatabase,
	}, nil
}
func (m *MongoDBManager) Close(ctx context.Context) {
	err := m.Client.Disconnect(ctx)
	if err != nil {
		logger.Log().Fatalf("Could not close storage, err: %v", err)
	}
}
