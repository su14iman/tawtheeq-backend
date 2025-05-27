package config

import (
	"fmt"
	"os"
	"time"

	"tawtheeq-backend/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var db *gorm.DB
	var err error

	maxAttempts := 10
	for i := 1; i <= maxAttempts; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		utils.HandleError(err, fmt.Sprintf("MySQL not ready, retry %d/%d", i, maxAttempts), utils.Warning)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		panic("❌ Failed to connect to DB: " + err.Error())
	}

	DB = db
	fmt.Println("✅ Connected to database")
}
