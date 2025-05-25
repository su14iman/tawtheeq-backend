package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID        string       `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string       `gorm:"not null"`
	LeaderID  string       `gorm:"type:char(36);not null"`
	Leader    User         `gorm:"foreignKey:LeaderID;references:ID"`
	Members   []TeamMember `gorm:"foreignKey:TeamID" json:"members"`
	CreatedAt time.Time    `gorm:"autoCreateTime"`
}

func (t *Team) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return
}

type CreateTeamInput struct {
	Name     string `json:"name"`
	LeaderID string `json:"leader_id"`
}

type TeamResponse struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Leader    UserShortResponse `json:"leader"`
	CreatedAt time.Time         `json:"created_at"`
}

type TeamShortResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func BuildTeamResponse(team *Team) TeamResponse {
	resp := TeamResponse{
		ID:        team.ID,
		Name:      team.Name,
		CreatedAt: team.CreatedAt,
	}

	resp.Leader = UserShortResponse{
		ID:       team.Leader.ID,
		FullName: team.Leader.FullName,
		Email:    team.Leader.Email,
	}

	return resp
}

type ChangeTeamNameInput struct {
	Name string `json:"name"`
}
