package repository

import (
	"ktrain/cmd/model"
	"ktrain/pkg/storage"
	"ktrain/pkg/tokens"
)

type IUserRepository interface {
	GetUserByID(id int64) (*model.User, error)
	GetAuthToken(token string) (*model.AuthToken, error)
	GetListUser() ([]*model.User, error)
	CreateUser(newUser *model.User) (*model.User, error)
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

func (r *userRepository) GetListUser() ([]*model.User, error) {
	users := []*model.User{}
	if err := r.db.Where("id > ?", 0).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (r *userRepository) CreateUser(newUser *model.User) (*model.User, error) {
	if err := r.db.Create(newUser).Error; err != nil {
		return nil, err
	}
	auth := &model.AuthToken{
		UserID: newUser.ID,
		Token:  tokens.CreateToken(newUser.ID, newUser.Username, newUser.Birthday, newUser.CreatedAt),
	}
	if err := r.db.Create(auth).Error; err != nil {
		return nil, err
	}
	return newUser, nil
}
