package api

import (
	_ "github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/docs"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the API routes for the non-versioned endpoints.
func RegisterRoutes(r *gin.Engine, handlers *Handlers) {
	r.GET("/docs/*any", handlers.DocsHandler)
	r.GET("/health", handlers.HealthHandler.HealthCheck)
}
