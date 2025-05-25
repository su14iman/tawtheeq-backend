package repositories

import (
	"tawtheeq-backend/models"

	"gorm.io/gorm"
)

type TeamMemberRepository struct {
	db *gorm.DB
}

func NewTeamMemberRepository(db *gorm.DB) *TeamMemberRepository {
	return &TeamMemberRepository{db}
}

func (r *TeamMemberRepository) Add(member *models.TeamMember) error {
	return r.db.Create(member).Error
}

func (r *TeamMemberRepository) Remove(teamID, userID string) error {
	return r.db.Delete(&models.TeamMember{}, "team_id = ? AND user_id = ?", teamID, userID).Error
}

func (r *TeamMemberRepository) GetMembers(teamID string) ([]models.TeamMember, error) {
	var members []models.TeamMember
	err := r.db.Where("team_id = ?", teamID).Find(&members).Error
	return members, err
}

func (r *TeamMemberRepository) Update(member *models.TeamMember) error {
	return r.db.Save(member).Error
}

func (r *TeamMemberRepository) GetMembersPaginated(teamID string, limit int, offset int) ([]models.TeamMember, error) {
	var members []models.TeamMember
	err := r.db.Where("team_id = ?", teamID).Limit(limit).Offset(offset).Find(&members).Error
	return members, err
}

func (r *TeamMemberRepository) CountMembers(teamID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.TeamMember{}).Where("team_id = ?", teamID).Count(&count).Error
	return count, err
}
