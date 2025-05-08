package core

import (
	"log"

	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/utils"
	"github.com/joho/godotenv"
)

// DBConfig holds the configuration for the database connection
type DBConfig struct {
	Address     string
	MaxConns    int
	MaxIdleTime string
}

// AppConfig holds all the application configuration
type AppConfig struct {
	AppName    string
	AppAddress string
	AppPort    string
	DBConfig   *DBConfig
}

// NewAppConfig initializes the application configuration
func NewAppConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ WARNING: Could not load .env file")
	}

	appName := utils.GetString("APP_NAME", "")
	appAddress := utils.GetString("APP_ADDRESS", "")
	appPort := utils.GetString("APP_PORT", "8080")

	dbConfig := &DBConfig{
		Address:     utils.GetString("DB_ADDRESS", "postgres://user:password@localhost:5432/dbname?sslmode=disable"),
		MaxConns:    utils.GetInt("DB_MAX_CONNECTIONS", 30),
		MaxIdleTime: utils.GetString("DB_MAX_IDLE_TIME", "15m"),
	}

	return &AppConfig{
		AppName:    appName,
		AppAddress: appAddress,
		AppPort:    appPort,
		DBConfig:   dbConfig,
	}
}
