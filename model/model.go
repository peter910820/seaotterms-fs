package model

import (
	"time"
)

// A00_Blog

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Username   string    `gorm:"NOT NULL unique" json:"username"`
	Password   string    `gorm:"NOT NULL" json:"-"`
	Email      string    `gorm:"NOT NULL unique" json:"email"`
	Avatar     string    `gorm:"NOT NULL; default:''" json:"avatar"`
	Exp        int       `gorm:"default:0" json:"exp"`
	Management bool      `gorm:"default:false" json:"management"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	CreateName string    `gorm:"NOT NULL" json:"createName"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	UpdateName string    `json:"updateName"`
}
