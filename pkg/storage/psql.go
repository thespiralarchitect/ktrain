package storage

import (
	"fmt"
	"ktrain/cmd/model"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PSQLManager struct {
	*gorm.DB
}

func NewPSQLManager() (*PSQLManager, error) {
	db, err := gorm.Open(postgres.Open(
		fmt.Sprintf(
			"host=%s user=%v password=%d dbname=%s port=%d sslmode=%s TimeZone=%s",
			viper.GetString("postgres.host"),
			viper.GetString("postgres.username"),
			viper.GetInt("postgres.password"),
			viper.GetString("postgres.database"),
			viper.GetInt("postgres.port"),
			viper.GetString("postgres.ssl_mode"),
			viper.GetString("postgres.timezone"),
		)))
	//db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=1234 dbname=Hieu port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"))
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&model.User{}, &model.AuthToken{})
	if err != nil {
		log.Printf("error create dattabase", err)
		return nil, err
	}
	return &PSQLManager{db.Debug()}, nil
}
