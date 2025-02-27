package main

import (
	"log"
	"scrapping/internals/utils"
	"scrapping/internals/web/router"

	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func init() {
	// setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
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
