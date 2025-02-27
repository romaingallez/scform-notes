package router

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"

	"scrapping/internals/utils"
)

// New creates and configures a new Fiber application with all middleware and settings
func New(engine *html.Engine) *fiber.App {
	app := fiber.New(
		fiber.Config{
			Views:       engine,
			ViewsLayout: "layouts/main",
		},
	)
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, HEAD, PUT, PATCH, POST, DELETE",
	}))

	engine.Reload(false)

	// Changing TimeZone & TimeFormat
	app.Use(logger.New(logger.Config{
		Format:     "${pid} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Europe/Paris",
	}))

	// if favicon exists, use it
	if _, err := os.Stat("./assets/src/favicon.ico"); err == nil {
		app.Use(favicon.New(favicon.Config{
			File: "./assets/src/favicon.ico",
			URL:  "/favicon.ico",
		}))
	}

	debug := os.Getenv("DEBUG")

	if strings.Contains(debug, "true") {
		app.Use(func(c *fiber.Ctx) error {
			ip := c.IP()
			log.Println("IP: ", ip)
			log.Println(utils.IsPrivateIP(ip))

			xForwardedFor := c.Get("X-Forwarded-For")
			log.Println("X-Forwarded-For: ", xForwardedFor)
			log.Println(utils.IsPrivateIP(xForwardedFor))

			return c.Next()
		})
	}

	return app
}
