package models

import (
	"time"
)

type PasswordResetToken struct {
	ID        string    `gorm:"type:char(36);primaryKey"`
	UserID    string    `gorm:"type:char(36);not null"`
	Token     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" example:"user@example.com"`
}

type ResetPasswordInput struct {
	Token       string `json:"token" example:"123e4567-e89b-12d3-a456-426614174000"`
	NewPassword string `json:"new_password" example:"MyNewPassword123"`
}
