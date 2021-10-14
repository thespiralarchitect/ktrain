package repository

import (
	"context"
	"fmt"
	"ktrain/cmd/api/user-api/dto"
	"ktrain/pkg/storage"

	"go.mongodb.org/mongo-driver/bson"
)

type ActivityLogRepository interface {
	CreateAction(ctx context.Context, id int64, activityLog string) (*dto.UserActivityLogMessage, error)
	GetAllLogAction(ctx context.Context, id int64) ([]*dto.ActionRequest, error)
}

type activityLogRepository struct {
	manager *storage.MongoDBManager
}

func NewActivityLogRepository(db *storage.MongoDBManager) ActivityLogRepository {
	return &activityLogRepository{
		manager: db,
	}
}
func (m *activityLogRepository) CreateAction(ctx context.Context, id int64, activityLog string) (*dto.UserActivityLogMessage, error) {
	action := dto.ActionRequest{
		ID:     id,
		Action: activityLog,
	}
	actionCollection := m.manager.Database.Collection("activityLog")
	_, err := actionCollection.InsertOne(ctx, action)
	if err != nil {
		return nil, err
	}
	action1, err := actionCollection.Find(ctx, bson.M{"user_id": id})
	if err != nil {
		return nil, err
	}
	var allAction []*dto.ActionRequest
	if err = action1.All(ctx, &allAction); err != nil {
		return nil, err
	}
	resp := dto.UserActivityLogMessage{
		ID:  id,
		Log: activityLog,
	}
	fmt.Println("ok5", allAction)
	return &resp, nil
}
func (m *activityLogRepository) GetAllLogAction(ctx context.Context, id int64) ([]*dto.ActionRequest, error) {
	actionCollection := m.manager.Database.Collection("activityLog")
	action, err := actionCollection.Find(ctx, bson.M{"user_id": id})
	if err != nil {
		return nil, err
	}
	var allAction []*dto.ActionRequest
	if err = action.All(ctx, &allAction); err != nil {
		return nil, err
	}
	return allAction, nil
}
