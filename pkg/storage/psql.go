package storage

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PSQLManager struct {
	*gorm.DB
}

func NewPSQLManager() (*PSQLManager, error) {
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=changeme dbname=samples port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"))
	if err != nil {
		return nil, err
	}

	return &PSQLManager{db.Debug()}, nil
}
