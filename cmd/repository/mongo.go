package repository

import (
	"context"
	"ktrain/cmd/api/user-api/dto"
	"ktrain/pkg/storage"

	"go.mongodb.org/mongo-driver/bson"
)

type ActivityLogRepository interface {
	CreateAction(ctx context.Context, id int64, action string) (string, error)
	GetAllLogAction(ctx context.Context, id int64) ([]*dto.ActionRequest, error)
}

type activityLogRepository struct {
	collection *storage.MongoDBManager
}

func NewActivityLogRepository(db *storage.MongoDBManager) ActivityLogRepository {
	return &activityLogRepository{
		collection: db,
	}
}
func (m *activityLogRepository) CreateAction(ctx context.Context, id int64, activityLog string) (string, error) {
	action := dto.ActionRequest{
		ID:     id,
		Action: activityLog,
	}
	_, err := m.collection.InsertOne(ctx, action)
	if err != nil {
		return "", err
	}
	return "Inserting document successfully", nil
}
func (m *activityLogRepository) GetAllLogAction(ctx context.Context, id int64) ([]*dto.ActionRequest, error) {
	action, err := m.collection.Find(ctx, bson.M{"user_id": id})
	if err != nil {
		return nil, err
	}
	var allAction []*dto.ActionRequest
	if err = action.All(ctx, &allAction); err != nil {
		return nil, err
	}
	return allAction, nil
}
