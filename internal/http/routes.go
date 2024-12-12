package http

import (
	"golang-auth/internal/http/controllers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/", controllers.HelloHandler)
}
