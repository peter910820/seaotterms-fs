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
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt"`
}

type FileResponse struct {
	Files       []string `json:"files"`
	Directories []string `json:"directories"`
}

func GenerateResponse(message string, data any) Response {
	return Response{
		Message:   message,
		TimeStamp: time.Now(),
		Data:      data,
	}
}
