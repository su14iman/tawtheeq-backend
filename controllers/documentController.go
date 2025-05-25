package controllers

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"tawtheeq-backend/config"
	"tawtheeq-backend/models"
	"tawtheeq-backend/repositories"
	"tawtheeq-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

// HideDocumentFromMe godoc
// @Summary Hide document
// @Description Hide document
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} models.DocumentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /documents/my/{id}/hide [get]
// @Security Bearer
func HideDocumentFromMe(c *fiber.Ctx) error {
	docRepo := repositories.NewDocumentRepository(config.DB)
	id := c.Params("id")
	userId := c.Locals("userID").(string)

	err := docRepo.HideFromUser(id, userId)
	if err != nil {
		utils.HandleError(err, fmt.Sprintf("Failed to hide document %s", id), utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hide document"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Document hidden successfully"})
}

// HideDocumentFromMyTeam godoc
// @Summary Hide document
// @Description Hide document
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} models.DocumentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /documents/myteam/{id}/hide [get]
// @Security Bearer
func HideDocumentFromMyTeam(c *fiber.Ctx) error {
	docRepo := repositories.NewDocumentRepository(config.DB)
	id := c.Params("id")
	teamID := c.Locals("team_id").(string)

	err := docRepo.HideFromTeam(id, teamID)
	if err != nil {
		utils.HandleError(err, fmt.Sprintf("Failed to hide document %s", id), utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hide document"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Document hidden successfully"})
}

// HideDocumentSuperAdmin godoc
// @Summary Hide document
// @Description Hide document
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} models.DocumentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /documents/{id}/hide [get]
// @Security Bearer
func HideDocumentSuperAdmin(c *fiber.Ctx) error {
	docRepo := repositories.NewDocumentRepository(config.DB)
	id := c.Params("id")

	err := docRepo.Hide(id)
	if err != nil {
		utils.HandleError(err, fmt.Sprintf("Failed to hide document %s", id), utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hide document"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Document hidden successfully"})
}

// ShowDocumentSuperAdmin	godoc
// @Summary Show document
// @Description Show document
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} models.DocumentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /documents/{id}/show [get]
// @Security Bearer
func ShowDocumentSuperAdmin(c *fiber.Ctx) error {
	docRepo := repositories.NewDocumentRepository(config.DB)
	id := c.Params("id")

	err := docRepo.Show(id)
	if err != nil {
		utils.HandleError(err, fmt.Sprintf("Failed to show document %s", id), utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to show document"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Document shown successfully"})
}

// GetAllDocumentsFromMyTeam godoc
// @Summary Get all documents from my team
// @Description Get all documents from my team
// @Tags documents
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} models.DocumentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /documents/myteam [get]
func GetAllDocumentsFromMyTeam(c *fiber.Ctx) error {
	docRepo := repositories.NewDocumentRepository(config.DB)
	teamID := c.Locals("team_id").(string)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	docs, err := docRepo.FindByTeamVisible(teamID, limit, offset)
	if err != nil {
		utils.HandleError(err, "Failed to fetch documents", utils.Warning)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch documents"})
	}

	return c.JSON(fiber.Map{
		"documents": docs,
		"meta": fiber.Map{
			"page":   page,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetAllDocumentsFromMeHidden godoc
// @Summary Get all documents from me hidden
// @Description Get all documents from me hidden
// @Tags documents
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} models.DocumentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /documents/me/hidden [get]
// @Security Bearer
func GetAllDocumentsFromMeHidden(c *fiber.Ctx) error {
	docRepo := repositories.NewDocumentRepository(config.DB)
	userId := c.Locals("userID").(string)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	docs, err := docRepo.FindByUserHidden(userId, limit, offset)
	if err != nil {
		utils.HandleError(err, "Failed to fetch documents", utils.Warning)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch documents"})
	}

	return c.JSON(fiber.Map{
		"documents": docs,
		"meta": fiber.Map{
			"page":   page,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetAllDocumentsFromMeVisible godoc
// @Summary Get all documents from me visible
// @Description Get all documents from me visible
// @Tags documents
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} models.DocumentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /documents/user/me [get]
// @Security Bearer
func GetAllDocumentsFromMeVisible(c *fiber.Ctx) error {
	docRepo := repositories.NewDocumentRepository(config.DB)
	userId := c.Locals("userID").(string)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	docs, err := docRepo.FindByUserVisible(userId, limit, offset)
	if err != nil {
		utils.HandleError(err, "Failed to fetch documents", utils.Warning)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch documents"})
	}

	return c.JSON(fiber.Map{
		"documents": docs,
		"meta": fiber.Map{
			"page":   page,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetAllDocumentsFromUserHidden godoc
// @Summary Get all documents from user hidden
// @Description Get all documents from user hidden
// @Tags documents
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} models.DocumentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /documents/user/{user_id}/hidden [get]
// @Security Bearer
func GetAllDocumentsFromUserHidden(c *fiber.Ctx) error {
	docRepo := repositories.NewDocumentRepository(config.DB)
	userID := c.Params("user_id")

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	docs, err := docRepo.FindByUserHidden(userID, limit, offset)
	if err != nil {
		utils.HandleError(err, "Failed to fetch documents", utils.Warning)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch documents"})
	}

	return c.JSON(fiber.Map{
		"documents": docs,
		"meta": fiber.Map{
			"page":   page,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetAllDocumentsFromUserVisible godoc
// @Summary Get all documents from user visible
// @Description Get all documents from user visible
// @Tags documents
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} models.DocumentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /documents/user/{user_id}/visible [get]
// @Security Bearer
func GetAllDocumentsFromUserVisible(c *fiber.Ctx) error {
	docRepo := repositories.NewDocumentRepository(config.DB)
	userID := c.Params("user_id")

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	docs, err := docRepo.FindByUserVisible(userID, limit, offset)
	if err != nil {
		utils.HandleError(err, "Failed to fetch documents", utils.Warning)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch documents"})
	}

	return c.JSON(fiber.Map{
		"documents": docs,
		"meta": fiber.Map{
			"page":   page,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetAllDocumentsFromTeamHidden godoc
// @Summary Get all documents from team hidden
// @Description Get all documents from team hidden
// @Tags documents
// @Accept json
// @Produce json
// @Param team_id path string true "Team ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} models.DocumentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /documents/team/{team_id}/hidden [get]
// @Security Bearer
func GetAllDocumentsFromTeamHidden(c *fiber.Ctx) error {
	docRepo := repositories.NewDocumentRepository(config.DB)
	teamID := c.Params("team_id")

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	docs, err := docRepo.FindByTeamHidden(teamID, limit, offset)
	if err != nil {
		utils.HandleError(err, "Failed to fetch documents", utils.Warning)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch documents"})
	}

	return c.JSON(fiber.Map{
		"documents": docs,
		"meta": fiber.Map{
			"page":   page,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetAllDocumentsFromTeamVisible godoc
// @Summary Get all documents from team visible
// @Description Get all documents from team visible
// @Tags documents
// @Accept json
// @Produce json
// @Param team_id path string true "Team ID"
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} models.DocumentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /documents/team/{team_id}/visible [get]
// @Security Bearer
func GetAllDocumentsFromTeamVisible(c *fiber.Ctx) error {
	docRepo := repositories.NewDocumentRepository(config.DB)
	teamID := c.Params("team_id")

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	docs, err := docRepo.FindByTeamVisible(teamID, limit, offset)
	if err != nil {
		utils.HandleError(err, "Failed to fetch documents", utils.Warning)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch documents"})
	}

	return c.JSON(fiber.Map{
		"documents": docs,
		"meta": fiber.Map{
			"page":   page,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetAllDocumentsHidden godoc
// @Summary Get all documents hidden
// @Description Get all documents hidden
// @Tags documents
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} models.DocumentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /documents/hidden [get]
// @Security Bearer
func GetAllDocumentsHidden(c *fiber.Ctx) error {
	docRepo := repositories.NewDocumentRepository(config.DB)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	docs, err := docRepo.FindAllHidden(limit, offset)
	if err != nil {
		utils.HandleError(err, "Failed to fetch documents", utils.Warning)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch documents"})
	}
	return c.JSON(fiber.Map{
		"documents": docs,
		"meta": fiber.Map{
			"page":   page,
			"limit":  limit,
			"offset": offset,
		},
	})

}

// GetAllDocumentsVisible godoc
// @Summary Get all documents visible
// @Description Get all documents visible
// @Tags documents
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param page query int false "Page"
// @Success 200 {object} models.DocumentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /documents/visible [get]
// @Security Bearer
func GetAllDocumentsVisible(c *fiber.Ctx) error {
	docRepo := repositories.NewDocumentRepository(config.DB)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	docs, err := docRepo.FindAllVisible(limit, offset)
	if err != nil {
		utils.HandleError(err, "Failed to fetch documents", utils.Warning)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch documents"})
	}
	return c.JSON(fiber.Map{
		"documents": docs,
		"meta": fiber.Map{
			"page":   page,
			"limit":  limit,
			"offset": offset,
		},
	})

}

// SignFileHandlerNew godoc
// @Summary Upload file
// @Description Upload file
// @Tags documents
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} models.UploadResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /upload [post]
// @Security Bearer
func SignFileHandler(c *fiber.Ctx) error {

	uploadDir := ""
	if os.Getenv("S3_ENABLED") == "true" {
		uploadDir = os.Getenv("TEMP_DIR")
		if uploadDir == "" {
			utils.HandleError(
				fmt.Errorf("TEMP_DIR not set"),
				"TEMP_DIR not set, using default ./temp",
				utils.Error,
			)
			uploadDir = "./temp"
		}
	} else {
		uploadDir = os.Getenv("LOCALLY_UPLOAD_DIR")
		if uploadDir == "" {
			utils.HandleError(
				fmt.Errorf("LOCALLY_UPLOAD_DIR not set"),
				"LOCALLY_UPLOAD_DIR not set, using default ./uploads",
				utils.Error,
			)
			uploadDir = "./uploads"
		}
	}

	// upload the file temporarily
	localFile, localPath, ext, id, localFileName, err := UploadFileLocal(c, uploadDir)
	if err != nil {
		return utils.HandleError(err, "Failed to upload file", utils.Error)
	}
	defer localFile.Close()

	// get the hash of the file
	hash, err := utils.CalculateFileHash(localPath)
	if err != nil {
		return utils.HandleError(err, "Failed to calculate file hash", utils.Error)
	}

	// check if the file already exists
	repoDocument := repositories.NewDocumentRepository(config.DB)
	existingDoc, _ := repoDocument.FindByHash(hash)
	if existingDoc != nil {
		// remove the temporary file
		if err := os.Remove(localPath); err != nil {
			utils.HandleError(err, "Failed to remove temporary file", utils.Warning)
		}

		// return the existing document
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error":    "File already exists",
				"createAt": existingDoc.CreatedAt,
				"document": existingDoc,
			},
		)
	}

	// generate a signature for the file
	signature, err := SignFile(localPath, id)
	if err != nil {
		return utils.HandleError(err, "Failed to sign and embed", utils.Error)
	}

	// update image or pdf with the signature
	if strings.HasSuffix(strings.ToLower(ext), ".jpg") || strings.HasSuffix(strings.ToLower(ext), ".jpeg") || strings.HasSuffix(strings.ToLower(ext), ".png") {
		if err := utils.AddIDToImage(localPath, id, signature); err != nil {
			return utils.HandleError(err, "Failed to add ID to image", utils.Error)
		}
	}

	if strings.HasSuffix(strings.ToLower(ext), ".pdf") {
		if err := utils.AddIDToPDF(localPath, id, signature); err != nil {
			return utils.HandleError(err, "Failed to add ID to PDF", utils.Error)
		}
	}

	// upload the file to S3
	newRandomHash, err := utils.CalculateFileHash(localPath)
	if err != nil {
		return utils.HandleError(err, "Failed to calculate file hash", utils.Error)
	}

	hashedFileName := ""
	if os.Getenv("S3_ENABLED") == "true" {
		hashedFileName = fmt.Sprintf("%s%s", newRandomHash, ext)
		bucket := os.Getenv("S3_BUCKET")

		// check if file already exists in S3
		_, statErr := config.S3Client.StatObject(context.Background(), bucket, hashedFileName, minio.StatObjectOptions{})
		if statErr != nil {
			fileBytes, err := os.ReadFile(localPath)
			if err != nil {
				return utils.HandleError(err, "Failed to read file for S3 upload", utils.Error)
			}

			_, err = config.S3Client.PutObject(context.Background(), bucket, hashedFileName, bytes.NewReader(fileBytes), int64(len(fileBytes)), minio.PutObjectOptions{
				ContentType: "application/octet-stream",
			})
			if err != nil {
				return utils.HandleError(err, "Failed to upload file to S3", utils.Error)
			}
		}

		// remove tmp file
		if err := os.Remove(localPath); err != nil {
			utils.HandleError(err, "Failed to remove temporary file", utils.Warning)
		}

	} else {
		hashedFileName = localFileName
	}

	// save in database
	userVal := c.Locals("userID")
	userId := ""
	if userVal != nil {
		userId, _ = userVal.(string)
	}
	teamVal := c.Locals("teamId")
	teamId := ""
	if teamVal != nil {
		teamId, _ = teamVal.(string)
	}

	if userId == "" {
		utils.HandleError(fmt.Errorf("userID not found"), "User ID not found", utils.Warning)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found"})
	}

	docRepo := repositories.NewDocumentRepository(config.DB)
	doc := &models.Document{
		OriginalName:   hashedFileName,
		FileFormat:     ext,
		Hash:           hash,
		Signature:      signature,
		SignedByUserID: userId,
	}

	if teamId != "" {
		doc.SignedByTeamID = &teamId
	}

	if err := docRepo.Create(doc); err != nil {
		utils.HandleError(err, "Failed to create document", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create document"})
	}

	return c.JSON(fiber.Map{
		"message":   "File signed and uploaded successfully",
		"file":      hashedFileName,
		"signature": signature,
		"document":  doc,
	})
}

// VerifyFileByIdHandler godoc
// @Summary Verify file by ID
// @Description Verify file by ID
// @Tags documents
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} models.DocumentResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /verify/{id} [get]
func VerifyFileByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	DocumentRepo := repositories.NewDocumentRepository(config.DB)
	doc, err := DocumentRepo.FindWithRelations(id)
	if err != nil {
		utils.HandleError(err, fmt.Sprintf("Document not found: %s", id), utils.Warning)
		return c.Status(404).JSON(models.ErrorResponse{Error: "Document not found", CreateAt: time.Now()})
	}
	return c.Status(200).JSON(doc)
}
