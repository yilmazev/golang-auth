package controllers

import (
	"golang-auth/internal/database"
	"golang-auth/internal/domain/auth"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var authService *auth.AuthService

func InitAuthService() {
	if database.DB == nil {
		log.Fatal("Database connection is not initialized")
	}
	authService = &auth.AuthService{
		Repo: &auth.AuthRepository{DB: database.DB},
	}
}

func RegisterHandler(c echo.Context) error {
	if authService == nil {
		log.Println("RegisterHandler: authService is not initialized")
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "AuthService is not initialized",
		})
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid request format",
		})
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Username, email, and password are required",
		})
	}

	token, err := authService.Register(auth.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Error creating user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "User registered successfully",
		Data:    map[string]string{"token": token},
	})
}

func LoginHandler(c echo.Context) error {
	if authService == nil {
		log.Println("LoginHandler: authService is not initialized")
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "AuthService is not initialized",
		})
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		log.Printf("LoginHandler: Error binding request: %v", err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid request format",
		})
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		log.Println("LoginHandler: Missing required fields")
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Username and password are required",
		})
	}

	token, err := authService.Login(req.Username, req.Password)
	if err != nil {
		log.Printf("LoginHandler: Invalid username or password: %v", err)
		return c.JSON(http.StatusUnauthorized, Response{
			Status:  "error",
			Message: "Invalid username or password",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Login successful",
		Data:    map[string]string{"token": token},
	})
}

func GetUserByUsernameHandler(c echo.Context) error {
	username := c.Param("username")
	if strings.TrimSpace(username) == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Username is required",
		})
	}

	user, err := authService.GetUserByUsername(username)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, Response{
				Status:  "error",
				Message: "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Failed to fetch user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "User fetched successfully",
		Data:    user,
	})
}

func ProfileHandler(c echo.Context) error {
	if authService == nil {
		log.Println("ProfileHandler: authService is not initialized")
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "AuthService is not initialized",
		})
	}

	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, Response{
			Status:  "error",
			Message: "Authorization header is missing",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return c.JSON(http.StatusUnauthorized, Response{
			Status:  "error",
			Message: "Invalid Authorization header format",
		})
	}

	username, err := authService.ParseToken(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, Response{
			Status:  "error",
			Message: "Invalid token",
		})
	}

	user, err := authService.GetUserByUsername(username)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, Response{
				Status:  "error",
				Message: "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Failed to fetch user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Profile fetched successfully",
		Data:    user,
	})
}
