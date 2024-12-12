package controllers

import (
	"context"
	"golang-auth/internal/database"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func DatabaseHealthCheck(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.DB.Ping(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "error",
			"message": "Database connection failed: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Database connection success",
	})
}
