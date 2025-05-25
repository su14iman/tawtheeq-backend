package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamMember struct {
	ID        string    `gorm:"type:char(36);primaryKey"`
	TeamID    string    `gorm:"type:char(36);not null"`
	Team      Team      `gorm:"foreignKey:TeamID;references:ID"`
	UserID    string    `gorm:"type:char(36);not null"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (m *TeamMember) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

type AddTeamMemberInput struct {
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`
}

type AddTeamMemberToMyInput struct {
	UserID string `json:"user_id"`
}

type TeamMemberResponse struct {
	Team      TeamShortResponse `json:"team"`
	User      UserShortResponse `json:"user"`
	CreatedAt time.Time         `json:"created_at"`
}

func BuildTeamMemberResponse(teamMember *TeamMember) TeamMemberResponse {
	resp := TeamMemberResponse{
		CreatedAt: teamMember.CreatedAt,
	}

	resp.Team = TeamShortResponse{
		ID:   teamMember.Team.ID,
		Name: teamMember.Team.Name,
	}

	resp.User = UserShortResponse{
		ID:       teamMember.User.ID,
		FullName: teamMember.User.FullName,
		Email:    teamMember.User.Email,
	}

	return resp
}
