package bootstrap

import (
	"golang-auth/internal/database"
	"golang-auth/internal/http/controllers"
	"log"
)

func Initialize() {
	database.InitDB()

	controllers.InitAuthService()

	log.Println("All services initialized successfully")
}

func Shutdown() {
	database.CloseDB()
	log.Println("All resources shut down successfully")
}
