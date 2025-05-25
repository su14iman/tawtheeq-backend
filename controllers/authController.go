// AuthController
package controllers

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"tawtheeq-backend/config"
	"tawtheeq-backend/models"
	"tawtheeq-backend/repositories"
	"tawtheeq-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary Login
// @Description Login user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.LoginInput true "Login Input"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/login [post]
// @Security Bearer
func Login(c *fiber.Ctx) error {
	expHours, _ := strconv.Atoi(os.Getenv("JWT_EXP_HOURS"))
	if expHours <= 0 {
		expHours = 72
	}

	var input models.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input", "created_at": time.Now()})
	}

	repo := repositories.NewUserRepository(config.DB)
	user, err := repo.FindByEmail(input.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password", "created_at": time.Now()})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password", "created_at": time.Now()})
	}

	// teamId
	teamId, err := getMyTeamID(user.ID)
	if err != nil {
		teamId = ""
	}

	claims := jwt.MapClaims{
		"id":     user.ID,
		"email":  user.Email,
		"role":   user.Role,
		"teamId": teamId,
		"exp":    time.Now().Add(time.Hour * time.Duration(expHours)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token", "created_at": time.Now()})
	}

	return c.JSON(fiber.Map{
		"token": signedToken,
		"user": fiber.Map{
			"id":     user.ID,
			"name":   user.FullName,
			"email":  user.Email,
			"role":   user.Role,
			"teamId": teamId,
		},
	})
}

// getMyTeamID retrieves the team ID associated with the given user ID.
func getMyTeamID(userId string) (string, error) {
	repo := repositories.NewTeamRepository(config.DB)
	team, err := repo.FindTeamByLeaderId(userId)
	if err == nil && team != nil {
		return team.ID, nil
	}
	team, err = repo.FindTeamByMemberId(userId)
	if err != nil || team == nil {
		return "", err
	}
	return team.ID, nil
}

// ForgotPassword godoc
// @Summary Send password reset email
// @Description Sends a reset password link to the user's email if it exists
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.ForgotPasswordInput true "User email"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/forgot-password [post]
func ForgotPassword(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	userRepo := repositories.NewUserRepository(config.DB)
	user, err := userRepo.FindByEmail(req.Email)
	if err != nil {
		return c.Status(200).JSON(fiber.Map{"message": "If email exists, reset link sent", "created_at": time.Now()})

	}

	token := utils.GenerateResetToken()
	resetLink := fmt.Sprintf("%s?token=%s", os.Getenv("FRONTEND_RESET_PASSWORD"), token)

	config.DB.Create(&models.PasswordResetToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})

	body := fmt.Sprintf("Click this link to reset your password:\n\n%s", resetLink)
	utils.SendEmail(user.Email, "Reset your password", body)

	return c.JSON(fiber.Map{"message": "If email exists, reset link sent", "created_at": time.Now()})
}

// ResetPassword godoc
// @Summary Reset password using token
// @Description Resets the user's password if the token is valid
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.ResetPasswordInput true "Token and new password"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/reset-password [post]
func ResetPassword(c *fiber.Ctx) error {
	type Request struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var reset models.PasswordResetToken
	err := config.DB.Where("token = ?", req.Token).First(&reset).Error
	if err != nil || time.Now().After(reset.ExpiresAt) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	userRepo := repositories.NewUserRepository(config.DB)
	user, err := userRepo.FindByID(reset.UserID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "User not found"})
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	user.Password = string(hashed)
	config.DB.Save(&user)

	config.DB.Delete(&reset)

	return c.JSON(fiber.Map{"message": "Password updated successfully"})
}
