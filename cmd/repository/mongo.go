package repository

import (
	"context"
	"fmt"
	"ktrain/cmd/api/user-api/dto"
	"ktrain/pkg/storage"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

type MongoRepository interface {
	CreateUser(user *dto.UserResquest) (*dto.UserResquest, error)
}

type mongoRepository struct {
	collection *storage.MongoDBManager
	ctx        context.Context
}

func NewMongoRepository(db *storage.MongoDBManager) MongoRepository {
	return &mongoRepository{
		collection: db,
	}
}
func (m *mongoRepository) CreateUser(user *dto.UserResquest) (*dto.UserResquest, error) {
	fmt.Println("ok2")
	quickstartDatabase := m.collection.Database("quickstart")
	podcastsCollection := quickstartDatabase.Collection("podcasts")
	_, err := podcastsCollection.InsertOne(m.ctx, user)
	if err != nil {
		log.Print("error occurred while inserting document in database :: ", err)
		return nil, err
	}
	var podcast bson.M
	if err = podcastsCollection.FindOne(m.ctx, bson.M{"fullname": "Hieu"}).Decode(&podcast); err != nil {
		log.Print("Object already exists in db:: ", err)
		return nil, err
	}
	fmt.Println(podcast)
	return user, nil
}
