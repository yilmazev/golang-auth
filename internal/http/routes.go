package http

import (
	"golang-auth/internal/http/controllers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/", controllers.HelloHandler)
	e.GET("/health", controllers.DatabaseHealthCheck)
	e.POST("/register", controllers.RegisterHandler)
	e.POST("/login", controllers.LoginHandler)
	e.GET("/profile", controllers.ProfileHandler)
	e.GET("/profile/:username", controllers.GetUserByUsernameHandler)
}
