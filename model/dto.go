package model

import "time"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Message   string    `json:"message"`
	TimeStamp time.Time `json:"timeStamp"`
	Data      any       `json:"data"`
}

type LoginResponse struct {
	Username   string    `gorm:"NOT NULL unique" json:"username"`
	Email      string    `gorm:"NOT NULL unique" json:"email"`
	Avatar     string    `gorm:"NOT NULL; default:''" json:"avatar"`
	Management bool      `gorm:"default:false" json:"management"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	CreateName string    `gorm:"NOT NULL" json:"createName"`
}

func GenerateResponse(message string, data any) Response {
	return Response{
		Message:   message,
		TimeStamp: time.Now(),
		Data:      data,
	}
}
