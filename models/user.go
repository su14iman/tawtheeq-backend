package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	SuperAdminRole Role = "super_admin"
	TeamLeaderRole Role = "team_leader"
	TeamMemberRole Role = "team_member"
)

type User struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	FullName  string    `gorm:"type:varchar(255)" json:"full_name"`
	Email     string    `gorm:"type:varchar(191);uniqueIndex" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      Role      `gorm:"type:varchar(50)" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}

type CreateUserInput struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	// Role     Role   `json:"role"`
}

type UserResponse struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Role      Role   `json:"role"`
	CreatedAt string `json:"created_at"`
}

type UserShortResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type ChangeUserRoleInput struct {
	Role string `json:"role" example:"team_leader"`
}
