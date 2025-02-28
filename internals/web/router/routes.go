package router

import (
	"scrapping/internals/web/handlers"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App) {
	// Create handlers
	gradeHandler := handlers.NewGradeHandler()

	// WebSocket middleware
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket route
	app.Get("/ws", websocket.New(handlers.WebSocketHandler))

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
