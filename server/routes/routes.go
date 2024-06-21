package routes

import (
	"facial-scan/controllers"
	"facial-scan/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/health", controllers.HealthCheck)
	api.Post("/register", middleware.AuthRequired(), controllers.Register)
	api.Post("/upload", middleware.AuthRequired(), controllers.UploadSaveAndRespond)
	api.Patch("/create-account", middleware.AuthRequired(), controllers.AccountCreation)
	api.Get("/has-username", middleware.AuthRequired(), controllers.HasUsername)
	api.Get("/dashboard", middleware.AuthRequired(), controllers.DashBoardLoad)
	api.Get("/get-score", middleware.AuthRequired(), controllers.GetScore)
	api.Get("/leaderboard-male", middleware.AuthRequired(), controllers.GetLeaderboardMale)
	api.Get("/leaderboard-female", middleware.AuthRequired(), controllers.GetLeaderboardFemale)
}
