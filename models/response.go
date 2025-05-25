package models

import (
	"time"
)

type ErrorResponse struct {
	Error    string    `json:"error"`
	CreateAt time.Time `json:"created_at"`
}

type SuccessResponse struct {
	Message  string    `json:"message"`
	CreateAt time.Time `json:"created_at"`
}

type JWTResponse struct {
	Token    string    `json:"token"`
	CreateAt time.Time `json:"created_at"`
}

type DocumentInfo struct {
	FilePath          string
	FileType          string
	FileSignature     string
	CountVerification int
	VerfiyUserID      UserShortResponse
	CreatedAt         time.Time
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"user"`
}
