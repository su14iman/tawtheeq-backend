package routes

import (
	"os"
	"tawtheeq-backend/controllers"
	"tawtheeq-backend/middlewares"
	"tawtheeq-backend/models"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Rate limiting middleware
	if os.Getenv("RATE_LIMIT_ENABLED") == "true" {
		app.Use(middlewares.RedisRateLimiter())
	}

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", controllers.Login)
	auth.Post("/forgot-password", controllers.ForgotPassword)
	auth.Post("/reset-password", controllers.ResetPassword)

	// Verify id
	api.Get("/verify/:id", controllers.VerifyFileByIdHandler)
	// Upload file
	api.Post("/upload", middlewares.RequireRoles("*"), controllers.SignFileHandler)

	// User manger
	users := api.Group("/users")
	users.Get("/", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.GetAllUsers)
	users.Post("/", middlewares.RequireRoles(string(models.SuperAdminRole), string(models.TeamLeaderRole)), controllers.CreateUser)
	users.Delete("/:id/remove", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.RemoveUser)
	users.Put("/:id/role", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.ChangeUserRole)

	// myself
	myself := api.Group("/myself")
	myself.Put("/name", middlewares.RequireRoles("*"), controllers.UpdateMyName)
	myself.Put("/password", middlewares.RequireRoles("*"), controllers.UpdateMyPassword)

	// Teams
	teams := api.Group("/teams")
	teams.Get("/", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.GetAllTeams)
	teams.Post("/", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.CreateTeam)
	teams.Delete("/:id/remove", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.RemoveTeam)
	teams.Put("/:id/name", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.UpdateTeamName)
	teams.Put("/:id/leader", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.UpdateTeamLeader)

	teams.Get("/:team_id/members", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.GetAllUsersInTeam)
	teams.Post("/members", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.AddUserToTeam)
	teams.Delete("/members/:team_id/:user_id", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.RemoveUserFromTeam)

	// my
	my := api.Group("/my")
	my.Get("/team", middlewares.RequireRoles("*"), controllers.GetMyTeam)
	my.Get("/team/members", middlewares.RequireRoles(string(models.TeamLeaderRole)), controllers.GetAllUsersInMyTeam)
	my.Post("/team/members", middlewares.RequireRoles(string(models.TeamLeaderRole)), controllers.AddUserToMyTeam)
	my.Delete("/team/members/:user_id", middlewares.RequireRoles(string(models.TeamLeaderRole)), controllers.RemoveUserFromMyTeam)

	// Documents
	documents := api.Group("/documents")
	documents.Get("/visible", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.GetAllDocumentsVisible)
	documents.Get("/hidden", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.GetAllDocumentsHidden)
	documents.Get("/team/:team_id/visible", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.GetAllDocumentsFromTeamVisible)
	documents.Get("/team/:team_id/hidden", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.GetAllDocumentsFromTeamHidden)
	documents.Get("/user/:user_id/visible", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.GetAllDocumentsFromUserVisible)
	documents.Get("/user/:user_id/hidden", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.GetAllDocumentsFromUserHidden)
	documents.Get("/user/me", middlewares.RequireRoles("*"), controllers.GetAllDocumentsFromMeVisible)
	documents.Get("/user/me/hidden", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.GetAllDocumentsFromMeHidden)
	documents.Get("/:id/hide", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.HideDocumentSuperAdmin)
	documents.Get("/:id/show", middlewares.RequireRoles(string(models.SuperAdminRole)), controllers.ShowDocumentSuperAdmin)

	documents.Get("/my", middlewares.RequireRoles("*"), controllers.GetAllDocumentsFromMeVisible)
	documents.Get("/myteam", middlewares.RequireRoles(string(models.TeamLeaderRole)), controllers.GetAllDocumentsFromMyTeam)
	documents.Get("/myteam/:id/hide", middlewares.RequireRoles(string(models.TeamLeaderRole)), controllers.HideDocumentFromMyTeam)
	documents.Get("/my/:id/hide", middlewares.RequireRoles("*"), controllers.HideDocumentFromMe)

}
