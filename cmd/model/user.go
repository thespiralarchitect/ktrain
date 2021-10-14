package model

import (
	"time"
)

type User struct {
	IsAdmin    bool
	ID         int64       `gorm:"type:integer"`
	Fullname   string      `gorm:"type:character varying(255)"`
	Username   string      `gorm:"type:character varying(255)"`
	Gender     string      `gorm:"type:character varying(10)"`
	Birthday   time.Time   `gorm:"type:timestamp"`
	AuthTokens []AuthToken `gorm:"foreignKey:UserID;references:ID"` //constraint:OnUpdate:CASCADE,OnDelete:CASCADE
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type AuthToken struct {
	ID        int64  `gorm:"type:integer"`
	UserID    int64  `gorm:"type:integer"`
	Token     string `gorm:"type:character varying(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
