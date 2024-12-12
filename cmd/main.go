package main

import (
	"golang-auth/config"
	"golang-auth/internal/bootstrap"
	"golang-auth/internal/http"
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	config.InitConfig()

	bootstrap.Initialize()
	defer bootstrap.Shutdown()

	e := echo.New()
	http.SetupRoutes(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Starting server on " + port)
	e.Logger.Fatal(e.Start(":" + port))
}
