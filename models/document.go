package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Document struct {
	ID                string `gorm:"type:char(36);primaryKey" json:"id"`
	OriginalName      string `gorm:"not null" json:"original_name"`
	FileFormat        string `gorm:"not null" json:"file_format"`
	Signature         string `gorm:"not null" json:"signature"`
	VerificationCount int    `gorm:"default:0" json:"verification_count"`
	IsHidden          bool   `gorm:"default:false" json:"is_hidden"`

	Hash string `gorm:"type:varchar(64);uniqueIndex"`

	SignedByUserID string `gorm:"type:uuid;not null" json:"signed_by_user_id"`
	SignedByUser   User   `gorm:"foreignKey:SignedByUserID" json:"signed_by_user"`

	SignedByTeamID *string `gorm:"type:uuid" json:"signed_by_team_id,omitempty"`
	SignedByTeam   *Team   `gorm:"foreignKey:SignedByTeamID" json:"signed_by_team,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *Document) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return
}

type CreateDocumentInput struct {
	OriginalName   string  `json:"original_name"`
	FileFormat     string  `json:"file_format"`
	Hash           string  `json:"hash"`
	Signature      string  `json:"signature"`
	SignedByUserID string  `json:"signed_by_user_id"`
	SignedByTeamID *string `json:"signed_by_team_id,omitempty"`
}

type DocumentResponse struct {
	ID                string             `json:"id"`
	OriginalName      string             `json:"original_name"`
	FileFormat        string             `json:"file_format"`
	VerificationCount int                `json:"verification_count"`
	Hash              string             `json:"hash"`
	SignedByUser      UserShortResponse  `json:"signed_by_user"`
	SignedByTeam      *TeamShortResponse `json:"signed_by_team,omitempty"`
	CreatedAt         string             `json:"created_at"`
	UpdatedAt         string             `json:"updated_at"`
}

type UploadResponse struct {
	FilePath  string `json:"file_path"`
	Signature string `json:"signature"`
	Hash      string `json:"hash"`
}

func BuildDocumentResponse(doc *Document) DocumentResponse {
	resp := DocumentResponse{
		ID:                doc.ID,
		OriginalName:      doc.OriginalName,
		FileFormat:        doc.FileFormat,
		VerificationCount: doc.VerificationCount,
		Hash:              doc.Hash,
		SignedByUser: UserShortResponse{
			ID:       doc.SignedByUser.ID,
			FullName: doc.SignedByUser.FullName,
			Email:    doc.SignedByUser.Email,
		},
		CreatedAt: doc.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: doc.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if doc.SignedByTeam != nil {
		resp.SignedByTeam = &TeamShortResponse{
			ID:   doc.SignedByTeam.ID,
			Name: doc.SignedByTeam.Name,
		}
	}

	return resp
}
