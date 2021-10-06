package repository

import (
	"fmt"
	"ktrain/cmd/model"
	"ktrain/pkg/storage"
)

type IUserRepository interface {
	GetUserByID(id int64) (*model.User, error)
	GetAuthToken(token string) (*model.AuthToken, error)
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(id int64) error
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
	fmt.Println(token)
	//s := "kPHAsagKV8ANrP7FtU6x5fhWeCXqVfuLS9799s2wDc9FBzQa3dBzag3ks3pyKFB2021-10-04 20:50:00"
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
	user := model.User{}
	if err := r.db.Where(&model.User{ID: id}).First(&user).Error; err != nil {
		return err

	}
	res := model.AuthToken{}
	if err := r.db.Where(&model.AuthToken{UserID: user.ID}).First(&res).Error; err != nil {
		return err
	}
	if err := r.db.Delete(&user).Error; err != nil {
		return err
	}
	if err := r.db.Delete(&res).Error; err != nil {
		return err
	}
	return nil
	//return r.db.Delete(&user).Error
	//return r.db.Unscoped().Delete(&user).Error
}
