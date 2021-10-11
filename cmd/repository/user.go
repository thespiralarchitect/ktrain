package repository

import (
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
	if err := r.db.Where(&model.User{ID: user.ID}).Updates(&model.User{Fullname: user.Fullname,
		Birthday: user.Birthday,
		Gender:   user.Gender}).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) DeleteUser(id int64) error {
	user, err := r.GetUserByID(id)
	if err != nil {
		return err
	}
	return r.db.Unscoped().Delete(&user).Error
}

func (r *userRepository) GetListUser(ids []int64) ([]*model.User, error) {
	users := []*model.User{}
	if ids == nil {
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
