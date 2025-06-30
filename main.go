package main

import (
	"log"
	"os"
	"scrapping/internals/utils"
	"scrapping/internals/web/router"
	"strings"

	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func init() {
	// setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ENV := strings.ToLower(os.Getenv("ENV"))

	// if ENV is prod or production, do not load .env file else load it
	if ENV != "prod" && ENV != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	}

	if os.Getenv("SCFORM_REMOTE_URL") != "" {
		err := utils.TestChromeDevWS(os.Getenv("SCFORM_REMOTE_URL"))
		if err != nil {
			log.Fatal("Error testing Chrome Dev WS:", err)
		}
	}

}

func main() {
	utils.InitAssets()

	engine := html.New("./views/", ".tpl")
	engine.Reload(false)

	// Initialize the router with all middleware
	app := router.New(engine)

	// Setup all routes
	router.SetupRoutes(app)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
