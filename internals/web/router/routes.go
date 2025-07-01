package router

import (
	"scrapping/internals/web/handlers"
	"scrapping/internals/web/session"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App, sessionManager *session.Manager) {
	// Create handlers
	gradeHandler := handlers.NewGradeHandler(sessionManager)

	// WebSocket middleware
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)

			// Get session and pass session ID to WebSocket connection
			sess, err := sessionManager.Store.Get(c)
			if err != nil {
				return c.Status(500).SendString("Failed to get session")
			}

			sessionID := sess.Get("id")
			if sessionID != nil {
				c.Locals("session_id", sessionID.(string))
			}

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
	app.Get("/about", gradeHandler.HandleAbout)
	app.Post("/grades", gradeHandler.HandleGrades)
	app.Post("/import", gradeHandler.HandleImport)
	app.Get("/search", gradeHandler.HandleSearch)
	app.Get("/api/grades", gradeHandler.HandleGradesAPI)
	app.Get("/print", gradeHandler.HandlePrint)
	app.Get("/print/demo", gradeHandler.HandlePrintDemo)
	app.Get("/export", gradeHandler.HandleExport)
	app.Get("/export/excel", gradeHandler.HandleExcelExport)
}
