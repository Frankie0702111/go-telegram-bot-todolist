package main

import (
	"fmt"
	gorm_utils "go-telegram-bot-todolist/utils/gorm"
	"go-telegram-bot-todolist/utils/log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = gorm_utils.InitMySQL()
)

func main() {
	defer gorm_utils.Close(db)

	// Load .env file
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Panic("Failed to load env file")
	}

	r := gin.Default()
	appPort := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	err := r.Run(appPort)
	if err != nil {
		log.Panic("Failed to run go-telegram")
	}
}
