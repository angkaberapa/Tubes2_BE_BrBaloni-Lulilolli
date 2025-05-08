package app

import (
	"log"

	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/api"
	v1 "github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/api/v1"
	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/core"
	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/db"
	"github.com/gin-gonic/gin"
)

// Run initializes the application context and starts the HTTP server.
// It sets up all necessary components, such as the database connection pool.
// It also registers the API routes.
func Run() {
	cfg := core.NewAppConfig()

	dbPool, err := db.NewPostgresPool(
		cfg.DBConfig.Address,
		cfg.DBConfig.MaxConns,
		cfg.DBConfig.MaxIdleTime,
	)

	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer dbPool.Close()

	appCtx := core.AppContext{
		Config: cfg,
		DBPool: dbPool,
	}

	handlers := api.InitHandlers(&appCtx)

	r := gin.Default()

	api.RegisterRoutes(r, handlers)
	v1.RegisterRoutes(r, handlers, &appCtx)

	// r.Run(cfg.AppAddress)
	r.Run(":8080") // For testing purposes, run on port 8080
}
