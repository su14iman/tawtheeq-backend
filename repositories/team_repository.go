package repositories

import (
	"tawtheeq-backend/models"

	"gorm.io/gorm"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db}
}

func (r *TeamRepository) Create(team *models.Team) error {
	return r.db.Create(team).Error
}

func (r *TeamRepository) FindByID(id string) (*models.Team, error) {
	var team models.Team
	err := r.db.Preload("Members").First(&team, "id = ?", id).Error
	return &team, err
}

func (r *TeamRepository) FindTeamByLeaderId(leaderId string) (*models.Team, error) {
	var team models.Team
	err := r.db.Preload("Members.User").First(&team, "leader_id = ?", leaderId).Error
	return &team, err
}

func (r *TeamRepository) FindTeamByMemberId(userID string) (*models.Team, error) {
	var team models.Team
	err := r.db.
		Joins("JOIN team_members ON team_members.team_id = teams.id").
		Where("team_members.user_id = ?", userID).
		Preload("Members.User").
		First(&team).Error

	return &team, err
}

func (r *TeamRepository) FindAll() ([]models.Team, error) {
	var teams []models.Team
	err := r.db.Find(&teams).Error
	return teams, err
}

func (r *TeamRepository) Update(team *models.Team) error {
	return r.db.Save(team).Error
}

func (r *TeamRepository) Delete(id string) error {
	return r.db.Delete(&models.Team{}, "id = ?", id).Error
}

func (r *TeamRepository) FindAllPaginated(limit int, offset int) ([]models.Team, error) {
	var teams []models.Team
	err := r.db.Limit(limit).Offset(offset).Find(&teams).Error
	return teams, err
}

func (r *TeamRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Team{}).Count(&count).Error
	return count, err
}
