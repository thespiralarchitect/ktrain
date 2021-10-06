package repository

import (
	"ktrain/cmd/model"
	"ktrain/pkg/storage"
)

type IUserRepository interface {
	GetUserByID(id int64) (*model.User, error)
	GetAuthToken(token string) (*model.AuthToken, error)
}

type userRepository struct {
	db *storage.PSQLManager
}

func NewUserRepository(db *storage.PSQLManager) IUserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUserByID(id int64) (*model.User, error) {
	user := &model.User{}
	if err := r.db.Where(&model.User{ID: id}).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetAuthToken(token string) (*model.AuthToken, error) {
	res := &model.AuthToken{}
	if err := r.db.Where(&model.AuthToken{Token: token}).First(res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
