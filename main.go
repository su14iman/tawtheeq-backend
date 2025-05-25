// @title Tawtheeq API
// @version 1.0
// @description RESTful API for signing and verifying files with digital signatures

// @host localhost:4000
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"

	"github.com/joho/godotenv"

	"tawtheeq-backend/config"
	"tawtheeq-backend/models"
	"tawtheeq-backend/routes"
	"tawtheeq-backend/utils"

	_ "tawtheeq-backend/docs"
)

func main() {

	utils.HandleError(fmt.Errorf("log test"), "This is a test log entry", utils.Warning)

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		utils.HandleError(err, "Failed to load .env file", utils.Error)
	}

	if err := utils.InitLogging(); err != nil {
		log.Fatal("Failed to init logging:", err)
	}
	defer utils.CloseLogging()

	port := os.Getenv("APP_PORT")
	if port == "" {
		utils.HandleError(
			fmt.Errorf("APP_PORT not set"),
			"APP_PORT not set, using default port 3000",
			utils.Error,
		)
		port = "3000"
	}

	S3_ENABLED := os.Getenv("S3_ENABLED")
	if S3_ENABLED == "true" {
		err := config.InitS3()
		if err != nil {
			utils.HandleError(err, "Failed to initialize S3", utils.Error)
		}
	}

	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		utils.HandleError(
			fmt.Errorf("FRONTEND_ORIGIN not set"),
			"FRONTEND_ORIGIN not set, using default http://localhost:3000",
			utils.Error,
		)
		frontendOrigin = "http://localhost:3000"
	}

	// Start Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: utils.MainErrorHandler,
	})

	// Middlewares
	app.Use(cors.New(cors.Config{
		AllowOrigins:     frontendOrigin,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// connect to database
	config.ConnectDatabase()
	// Migrate models
	config.DB.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.TeamMember{},
		&models.Document{},
		models.PasswordResetToken{},
	)
	// Create super admin if not exists
	config.CreateSuperAdminIfNotExists()

	if os.Getenv("ENABLE_SWAGGER") == "true" {
		// Register Swagger docs handler
		app.Get("/swagger/*", swagger.HandlerDefault)
		log.Println("📚 Swagger enabled at /swagger/index.html")
	}

	// Initialize Redis rate limiting
	config.InitRateLimiting()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":" + port))
}
