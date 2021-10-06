package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id          int64  `gorm:"type:integer;not null:primarykey"`
	Fullname    string `gorm :"type:varchar(255) ;not null"`
	Username    string `gorm :"type:varchar(255) ;not null"`
	Gender      string `gorm :"type:varchar(10) ;not null"`
	Birthday    string `gorm :"type:varchar(10) ;not null"`
	CreatedAt   time.Time
	UpdateAt    time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Auth_tokens []*Auth_token  `gorm:"foreignKey:UserID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Auth_token struct {
	Id        int64  `gorm:"type:integer;not null:primarykey"`
	UserID    int64  `gorm:"type:integer;not null"`
	Tocken    string `gorm :"type:varchar(255) ;not null"`
	CreatedAt time.Time
	UpdateAt  time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
