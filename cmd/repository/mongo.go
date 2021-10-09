package repository

import (
	"context"
	"ktrain/cmd/api/user-api/dto"
	"ktrain/pkg/storage"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type MongoRepository interface {
	CreateAction(ctx context.Context, id int64, action string) (string, error)
	GetAllLogAction(ctx context.Context, id int64) ([]*dto.ActionResquest, error)
}

type mongoRepository struct {
	collection *storage.MongoDBManager
}

func NewMongoRepository(db *storage.MongoDBManager) MongoRepository {
	return &mongoRepository{
		collection: db,
	}
}
func (m *mongoRepository) CreateAction(ctx context.Context, id int64, activity_log string) (string, error) {
	actionDatabase := m.collection.Database("Action")
	actionCollection := actionDatabase.Collection("activity_logs")
	action := dto.ActionResquest{
		ID:     id,
		Action: activity_log,
	}
	_, err := actionCollection.InsertOne(ctx, action)
	if err != nil {
		log.Print("error occurred while inserting document in database :: ", err)
		return "", err
	}
	return "Inserting document successfully", nil
}
func (m *mongoRepository) GetAllLogAction(ctx context.Context, id int64) ([]*dto.ActionResquest, error) {
	actionDatabase := m.collection.Database("Action")
	actionCollection := actionDatabase.Collection("activity_logs")
	action, err := actionCollection.Find(ctx, bson.M{"user_id": id})
	if err != nil {
		return nil, err
	}
	var allAction []*dto.ActionResquest
	if err = action.All(ctx, &allAction); err != nil {
		return nil, err
	}
	return allAction, nil
}
