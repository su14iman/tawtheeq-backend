package config

import (
	"fmt"
	"os"

	"tawtheeq-backend/models"
	"tawtheeq-backend/repositories"

	"golang.org/x/crypto/bcrypt"
)

func CreateSuperAdminIfNotExists() {
	email := os.Getenv("SUPERADMIN_EMAIL")
	password := os.Getenv("SUPERADMIN_PASSWORD")

	if email == "" || password == "" {
		fmt.Println("⚠️  SUPERADMIN credentials not found in .env")
		return
	}

	repo := repositories.NewUserRepository(DB)

	existing, err := repo.FindByEmail(email)
	if err == nil {
		if existing.Role != models.SuperAdminRole {
			existing.Role = models.SuperAdminRole
			repo.Update(existing)
			fmt.Println("✅ Updated existing user to SuperAdmin role")
		} else {
			fmt.Println("✅ SuperAdmin already exists")
		}
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("❌ Error hashing password:", err)
		return
	}

	newSuper := &models.User{
		FullName: "Super Admin",
		Email:    email,
		Password: string(hashedPassword),
		Role:     models.SuperAdminRole,
	}

	if _, err := repo.Create(newSuper); err != nil {
		fmt.Println("❌ Error creating SuperAdmin:", err)
		return
	}

	fmt.Println("✅ SuperAdmin created from env")
}
