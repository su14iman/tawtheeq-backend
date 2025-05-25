package config

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var S3Client *minio.Client

func InitS3() error {
	if os.Getenv("S3_ENABLED") != "true" {
		fmt.Println("✅ S3 is disabled.")
		return nil
	}

	endpoint := os.Getenv("S3_ENDPOINT")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	bucket := os.Getenv("S3_BUCKET")
	secure, _ := strconv.ParseBool(os.Getenv("S3_SECURE")) // optional, default false

	var err error
	S3Client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		return fmt.Errorf("❌ failed to init S3 client: %w", err)
	}

	exists, err := S3Client.BucketExists(context.Background(), bucket)
	if err != nil {
		return fmt.Errorf("❌ failed to check bucket: %w", err)
	}
	if !exists {
		err = S3Client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("❌ failed to create bucket: %w", err)
		}
		fmt.Println("✅ S3 bucket created:", bucket)
	} else {
		fmt.Println("✅ S3 bucket already exists:", bucket)
	}

	fmt.Println("✅ S3 initialized successfully")
	return nil
}
