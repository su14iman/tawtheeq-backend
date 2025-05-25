package controllers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"tawtheeq-backend/config"
	"tawtheeq-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
)

func UploadFileLocal(c *fiber.Ctx, uploadDir string) (*os.File, string, string, string, string, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, "", "", "", "", utils.HandleError(err, "Error loading .env file", utils.Error)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, "", "", "", "", utils.HandleError(err, "Invalid file", utils.Error)
	}

	fileName := fileHeader.Filename
	id := uuid.New().String()

	// Create upload directory if it doesn't exist
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return nil, "", "", "", "", utils.HandleError(err, "Failed to create upload directory", utils.Error)
		}
	}

	ext := filepath.Ext(fileHeader.Filename)
	savePath := fmt.Sprintf("%s/%s%s", uploadDir, id, ext)

	if err := c.SaveFile(fileHeader, savePath); err != nil {
		return nil, "", "", "", "", utils.HandleError(err, "Failed to save file", utils.Error)
	}

	uploadedFile, err := os.Open(savePath)
	if err != nil {
		return nil, "", "", "", "", utils.HandleError(err, "Failed to open file", utils.Error)
	}

	return uploadedFile, savePath, ext, id, fileName, nil
}

func UploadFile(c *fiber.Ctx) (hash string, savePath string, ext string, id string, fileName string, err error) {
	err = godotenv.Load()
	if err != nil {
		err = utils.HandleError(err, "Error loading .env file", utils.Error)
		return
	}

	var fileHeader *multipart.FileHeader
	fileHeader, err = c.FormFile("file")
	if err != nil {
		err = utils.HandleError(err, "Invalid file", utils.Error)
		return
	}

	var file multipart.File
	file, err = fileHeader.Open()
	if err != nil {
		err = utils.HandleError(err, "Failed to open file", utils.Error)
		return
	}
	defer file.Close()

	// Prepare metadata
	fileName = fileHeader.Filename
	ext = filepath.Ext(fileName)
	id = uuid.New().String()

	// Calculate hash and buffer file
	hasher := sha256.New()
	var buf bytes.Buffer
	tee := io.TeeReader(file, &buf)
	_, err = io.Copy(hasher, tee)
	if err != nil {
		err = utils.HandleError(err, "Failed to hash file", utils.Error)
		return
	}
	hash = hex.EncodeToString(hasher.Sum(nil))
	hashedFileName := fmt.Sprintf("%s%s", hash, ext)

	// S3 upload
	if os.Getenv("S3_ENABLED") == "true" {
		bucket := os.Getenv("S3_BUCKET")
		_, s3Err := config.S3Client.StatObject(context.Background(), bucket, hashedFileName, minio.StatObjectOptions{})
		if s3Err == nil {
			savePath = hashedFileName
			return
		}
		_, s3Err = config.S3Client.PutObject(context.Background(), bucket, hashedFileName, &buf, int64(buf.Len()), minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		})
		if s3Err != nil {
			err = utils.HandleError(s3Err, "Failed to upload file to S3", utils.Error)
			return
		}
		savePath = hashedFileName
		return
	}

	// Local upload fallback
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	if _, err = os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			err = utils.HandleError(err, "Failed to create upload directory", utils.Error)
			return
		}
	}

	savePath = filepath.Join(uploadDir, hashedFileName)
	if _, err = os.Stat(savePath); os.IsNotExist(err) {
		err = os.WriteFile(savePath, buf.Bytes(), 0644)
		if err != nil {
			err = utils.HandleError(err, "Failed to save file locally", utils.Error)
			return
		}
	}

	return
}

func RemoveFile(filePath string) error {
	if os.Getenv("S3_ENABLED") == "true" {
		bucket := os.Getenv("S3_BUCKET")

		key := filePath
		if strings.Contains(filePath, "/") {
			parts := strings.Split(filePath, "/")
			key = parts[len(parts)-1]
		}

		err := config.S3Client.RemoveObject(context.Background(), bucket, key, minio.RemoveObjectOptions{})
		if err != nil {
			return utils.HandleError(err, "Failed to remove file from S3", utils.Error)
		}
		return nil
	}

	if err := os.Remove(filePath); err != nil {
		return utils.HandleError(err, "Failed to remove file locally", utils.Error)
	}
	return nil
}

func GetFileHashFromPath(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		utils.HandleError(err, "Failed to open file", utils.Error)
		return ""
	}
	defer file.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		utils.HandleError(err, "Failed to hash file", utils.Error)
		return ""
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
