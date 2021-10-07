package storage

import (
	"fmt"

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
			"host=%s user=%v password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			viper.GetString("postgres.host"),
			viper.GetString("postgres.username"),
			viper.GetString("postgres.password"),
			viper.GetString("postgres.database"),
			viper.GetInt("postgres.port"),
			viper.GetString("postgres.ssl_mode"),
			viper.GetString("postgres.timezone"),
		)))

	if err != nil {
		return nil, err
	}
	return &PSQLManager{db.Debug()}, nil
}
