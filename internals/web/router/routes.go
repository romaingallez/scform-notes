package router

import (
	"scrapping/internals/web/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App) {
	// Create handlers
	gradeHandler := handlers.NewGradeHandler()

	// Static routes
	app.Static("/static", "./static")
	app.Static("/assets", "./assets/dist")

	// API routes
	app.Get("/", gradeHandler.HandleIndex)
	app.Post("/grades", gradeHandler.HandleGrades)
	app.Get("/search", gradeHandler.HandleSearch)
	app.Get("/print", gradeHandler.HandlePrint)
	app.Get("/export", gradeHandler.HandleExport)
	app.Get("/export/excel", gradeHandler.HandleExcelExport)
}
