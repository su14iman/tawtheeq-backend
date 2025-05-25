package repositories

import (
	"tawtheeq-backend/models"

	"gorm.io/gorm"
)

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db}
}

func (r *DocumentRepository) Create(doc *models.Document) error {
	return r.db.Create(doc).Error
}

func (r *DocumentRepository) FindByID(id string) (*models.Document, error) {
	var doc models.Document
	err := r.db.First(&doc, "id = ?", id).Error

	if err == nil {
		doc.VerificationCount++
		err = r.db.Save(&doc).Error
		if err != nil {
			return nil, err
		}
	}

	return &doc, err
}

func (r *DocumentRepository) FindByHash(hash string) (*models.Document, error) {
	var doc models.Document
	err := r.db.Where("hash = ?", hash).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *DocumentRepository) FindByIDHidden(id string) (*models.Document, error) {
	var doc models.Document
	err := r.db.First(&doc, "id = ? AND is_hidden = ?", id, true).Error
	return &doc, err
}

func (r *DocumentRepository) FindByIDVisible(id string) (*models.Document, error) {
	var doc models.Document
	err := r.db.First(&doc, "id = ? AND is_hidden = ?", id, false).Error
	return &doc, err
}

func (r *DocumentRepository) FindAll(limit int, offset int) ([]models.Document, error) {
	var docs []models.Document
	err := r.db.Limit(limit).Offset(offset).Find(&docs).Error
	return docs, err
}

func (r *DocumentRepository) FindAllHidden(limit int, offset int) ([]models.Document, error) {
	var docs []models.Document
	err := r.db.Where("is_hidden = ?", true).Limit(limit).Offset(offset).Find(&docs).Error
	return docs, err
}

func (r *DocumentRepository) FindAllVisible(limit int, offset int) ([]models.Document, error) {
	var docs []models.Document
	err := r.db.Where("is_hidden = ?", false).Limit(limit).Offset(offset).Find(&docs).Error
	return docs, err
}

func (r *DocumentRepository) FindByUser(userID string, limit int, offset int) ([]models.Document, error) {
	var docs []models.Document
	err := r.db.Where("signed_by_user_id = ?", userID).Limit(limit).Offset(offset).Find(&docs).Error
	return docs, err
}

func (r *DocumentRepository) FindByUserHidden(userID string, limit int, offset int) ([]models.Document, error) {
	var docs []models.Document
	err := r.db.Where("signed_by_user_id = ? AND is_hidden = ?", userID, true).Limit(limit).Offset(offset).Find(&docs).Error
	return docs, err
}

func (r *DocumentRepository) FindByUserVisible(userID string, limit int, offset int) ([]models.Document, error) {
	var docs []models.Document
	err := r.db.Where("signed_by_user_id = ? AND is_hidden = ?", userID, false).Limit(limit).Offset(offset).Find(&docs).Error
	return docs, err
}

func (r *DocumentRepository) FindByTeam(teamID string, limit int, offset int) ([]models.Document, error) {
	var docs []models.Document
	err := r.db.Where("signed_by_team_id = ?", teamID).Limit(limit).Offset(offset).Find(&docs).Error
	return docs, err
}

func (r *DocumentRepository) FindByTeamHidden(teamID string, limit int, offset int) ([]models.Document, error) {
	var docs []models.Document
	err := r.db.Where("signed_by_team_id = ? AND is_hidden = ?", teamID, true).Limit(limit).Offset(offset).Find(&docs).Error
	return docs, err
}

func (r *DocumentRepository) FindByTeamVisible(teamID string, limit int, offset int) ([]models.Document, error) {
	var docs []models.Document
	err := r.db.Where("signed_by_team_id = ? AND is_hidden = ?", teamID, false).Limit(limit).Offset(offset).Find(&docs).Error
	return docs, err
}

func (r *DocumentRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Document{}).Count(&count).Error
	return count, err
}

func (r *DocumentRepository) CountAllHidden() (int64, error) {
	var count int64
	err := r.db.Model(&models.Document{}).Where("is_hidden = ?", true).Count(&count).Error
	return count, err
}

func (r *DocumentRepository) CountAllVisible() (int64, error) {
	var count int64
	err := r.db.Model(&models.Document{}).Where("is_hidden = ?", false).Count(&count).Error
	return count, err
}

func (r *DocumentRepository) Update(doc *models.Document) error {
	return r.db.Save(doc).Error
}

func (r *DocumentRepository) Delete(id string) error {
	return r.db.Delete(&models.Document{}, "id = ?", id).Error
}

func (r *DocumentRepository) Hide(id string) error {
	return r.db.Model(&models.Document{}).Where("id = ?", id).Update("is_hidden", true).Error
}

func (r *DocumentRepository) HideFromTeam(id string, teamID string) error {
	return r.db.Model(&models.Document{}).Where("id = ? AND signed_by_team_id = ?", id, teamID).Update("is_hidden", true).Error
}

func (r *DocumentRepository) HideFromUser(id string, userID string) error {
	return r.db.Model(&models.Document{}).Where("id = ? AND signed_by_user_id = ?", id, userID).Update("is_hidden", true).Error
}

func (r *DocumentRepository) Show(id string) error {
	return r.db.Model(&models.Document{}).Where("id = ?", id).Update("is_hidden", false).Error
}

func (r *DocumentRepository) FindWithRelations(id string) (*models.Document, error) {
	var doc models.Document
	err := r.db.
		Preload("SignedByUser").
		Preload("SignedByTeam").
		First(&doc, "id = ?", id).Error
	return &doc, err
}

func (r *DocumentRepository) FindWithRelationsHidden(id string) (*models.Document, error) {
	var doc models.Document
	err := r.db.
		Preload("SignedByUser").
		Preload("SignedByTeam").
		First(&doc, "id = ? AND is_hidden = ?", id, true).Error
	return &doc, err
}

func (r *DocumentRepository) FindWithRelationsVisible(id string) (*models.Document, error) {
	var doc models.Document
	err := r.db.
		Preload("SignedByUser").
		Preload("SignedByTeam").
		First(&doc, "id = ? AND is_hidden = ?", id, false).Error
	return &doc, err
}
