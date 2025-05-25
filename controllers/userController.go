package controllers

import (
	"strconv"
	"time"

	"tawtheeq-backend/config"
	"tawtheeq-backend/models"
	"tawtheeq-backend/repositories"
	"tawtheeq-backend/utils"

	"github.com/gofiber/fiber/v2"
)

// CreateUser godoc
// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.CreateUserInput true "Create User Input"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users [post]
// @Security Bearer
func CreateUser(c *fiber.Ctx) error {

	var input models.CreateUserInput
	if err := c.BodyParser(&input); err != nil {
		utils.HandleError(err, "Failed to parse request body", utils.Error)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid input",
			"created_at": time.Now(),
		})
	}

	repo := repositories.NewUserRepository(config.DB)

	password, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.HandleError(err, "Failed to hash password", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Failed to hash password",
			"created_at": time.Now(),
		})
	}

	newUser := models.User{
		// ID:       uuid.New().String(),
		FullName: input.FullName,
		Email:    input.Email,
		Password: password,
		Role:     models.Role(models.TeamMemberRole), // Default role
		// Role:     models.Role(c.FormValue("role")),
	}

	user, err := repo.Create(&newUser)

	if err != nil {
		utils.HandleError(err, "Failed to create user", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Failed to create user",
			"created_at": time.Now(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user":    user,
	})
}

// RemoveUser godoc
// @Summary Remove user
// @Description Remove user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.UserShortResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [delete]
// @Security Bearer
func RemoveUser(c *fiber.Ctx) error {
	repo := repositories.NewUserRepository(config.DB)
	id := c.Params("id")
	err := repo.Delete(id)
	if err != nil {
		utils.HandleError(err, "Failed to delete user", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
		"userID":  id,
	})
}

// ChangeUserRole godoc
// @Summary Change user role
// @Description Allows Super Admin to change a user's role
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param input body models.ChangeUserRoleInput true "New role"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id}/role [put]
// @Security Bearer
func ChangeUserRole(c *fiber.Ctx) error {

	id := c.Params("id")

	var input models.ChangeUserRoleInput
	if err := c.BodyParser(&input); err != nil {
		utils.HandleError(err, "Failed to parse request body", utils.Error)
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid input", CreateAt: time.Now()})
	}

	userRepo := repositories.NewUserRepository(config.DB)
	user, err := userRepo.FindByID(id)
	if err != nil {
		utils.HandleError(err, "User not found", utils.Error)
		return c.Status(404).JSON(models.ErrorResponse{Error: "User not found", CreateAt: time.Now()})
	}

	user.Role = models.Role(input.Role)
	if err := userRepo.Update(user); err != nil {
		utils.HandleError(err, "Failed to update user role", utils.Error)
		return c.Status(500).JSON(models.ErrorResponse{Error: "Failed to update user role", CreateAt: time.Now()})
	}

	return c.JSON(fiber.Map{
		"message":   "User role updated successfully",
		"user_id":   user.ID,
		"new_role":  user.Role,
		"updatedAt": time.Now(),
	})
}

// UpdateMyName godoc
// @Summary Change user name
// @Description Allows user to change their own name
// @Tags myself
// @Accept json
// @Produce json
// @Param full_name query string true "New full name"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /myself/name [put]
// @Security Bearer
// @Security JWT
// @Security ApiKeyAuth
func UpdateMyName(c *fiber.Ctx) error {
	repo := repositories.NewUserRepository(config.DB)

	userIDRaw := c.Locals("userID")
	userID, ok1 := userIDRaw.(string)
	if !ok1 {
		utils.HandleError(nil, "Invalid or missing user context", utils.Error)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid or missing user context",
			"created_at": time.Now(),
		})
	}

	user, err := repo.FindByID(userID)
	if err != nil {
		utils.HandleError(err, "User not found", utils.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	user.FullName = c.FormValue("full_name")

	err = repo.Update(user)
	if err != nil {
		utils.HandleError(err, "Failed to update user name", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user name",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User name updated successfully",
		"user":    user,
	})
}

// UpdateMyPassword godoc
// @Summary Change user password
// @Description Allows user to change their own password
// @Tags myself
// @Accept json
// @Produce json
// @Param old_password query string true "Old password"
// @Param password query string true "New password"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /myself/password [put]
// @Security Bearer
// @Security JWT
// @Security ApiKeyAuth
func UpdateMyPassword(c *fiber.Ctx) error {
	repo := repositories.NewUserRepository(config.DB)

	userIDRaw := c.Locals("userID")
	userID, ok1 := userIDRaw.(string)
	if !ok1 {
		utils.HandleError(nil, "Invalid or missing user context", utils.Error)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Invalid or missing user context",
			"created_at": time.Now(),
		})
	}

	user, err := repo.FindByID(userID)
	if err != nil {
		utils.HandleError(err, "User not found", utils.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":      "User not found",
			"created_at": time.Now(),
		})
	}

	passwordOld, err := utils.HashPassword(c.FormValue("old_password"))
	if err != nil {
		utils.HandleError(err, "Failed to hash old password", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Failed to hash password",
			"created_at": time.Now(),
		})
	}
	password, err := utils.HashPassword(c.FormValue("password"))
	if err != nil {
		utils.HandleError(err, "Failed to hash new password", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Failed to hash password",
			"created_at": time.Now(),
		})
	}

	if user.Password != passwordOld {
		utils.HandleError(nil, "Old password is incorrect", utils.Error)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Old password is incorrect",
			"created_at": time.Now(),
		})
	}

	user.Password = password
	err = repo.Update(user)
	if err != nil {
		utils.HandleError(err, "Failed to update user password", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":      "Failed to update user password",
			"created_at": time.Now(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "User password updated successfully",
		"created_at": time.Now(),
	})
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get all users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param page query int false "Page" default(1)
// @Success 200 {array} models.UserShortResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users [get]
// @Security Bearer
// @Security JWT
// @Security ApiKeyAuth
func GetAllUsers(c *fiber.Ctx) error {
	repo := repositories.NewUserRepository(config.DB)

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	users, err := repo.FindAll(limit, offset)
	if err != nil {
		utils.HandleError(err, "Failed to retrieve users", utils.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}
	return c.Status(fiber.StatusOK).JSON(users)
}
