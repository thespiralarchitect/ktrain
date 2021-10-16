package repository

import (
	"errors"
	"ktrain/cmd/model"
	"ktrain/pkg/storage"
	"ktrain/pkg/tokens"
)

type IUserRepository interface {
	GetUserByID(id int64) (*model.User, error)
	GetAuthToken(token string) (*model.AuthToken, error)
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(id int64) error
	GetListUser(ids []int64) ([]*model.User, error)
	CreateUser(newUser *model.User) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
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
	res := model.AuthToken{}
	if err := r.db.Where(&model.AuthToken{Token: token}).First(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}
func (r *userRepository) UpdateUser(user *model.User) (*model.User, error) {
	q := r.db.Where(&model.User{ID: user.ID}).Updates(&model.User{Fullname: user.Fullname,
		Birthday: user.Birthday,
		Gender:   user.Gender})
	if q.Error != nil {
		return nil, q.Error
	}
	if q.RowsAffected == 0 {
		return nil, errors.New("no field update value")
	}
	return user, nil
}

func (r *userRepository) DeleteUser(id int64) error {
	if err := r.db.Where(&model.AuthToken{UserID: id}).Delete(&model.AuthToken{UserID: id}).Error; err != nil {
		return err
	}
	if err := r.db.Where(&model.User{ID: id}).Delete(&model.User{ID: id}).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetListUser(ids []int64) ([]*model.User, error) {
	users := []*model.User{}
	if len(ids) == 0 {
		if err := r.db.Where("id > ?", 0).Find(&users).Error; err != nil {
			return nil, err
		}
	} else {
		for _, id := range ids {
			user := &model.User{}
			if err := r.db.Where(&model.User{ID: id}).First(user).Error; err != nil {
				return nil, err
			}
			users = append(users, user)
		}
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
func (r *userRepository) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	if err := r.db.Where(&model.User{Username: username}).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
