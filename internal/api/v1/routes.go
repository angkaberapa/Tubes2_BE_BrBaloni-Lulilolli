package v1

import (
	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/api"
	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/core"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the API routes for the v1 version of the API.
func RegisterRoutes(r *gin.Engine, handlers *api.Handlers, appCtx *core.AppContext) {
	r.Use(cors.Default())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/scrape", ScrapeHandler)
		v1.GET("/images", ImageListHandler)
		v1.GET("/search", SearchHandler)
		v1.POST("/search", SearchHandler)
		// v1.GET("/dfs", DFSHandler)
		// v1.GET("/dfs-concurrent", DFSMultipleRecipeHandler)
	}
	// authGroup := v1.Group("/auth")
	// {
	// 	authGroup.GET("/signin", handlers.AuthHandler.SignInV1)
	// 	authGroup.GET("/signout", handlers.AuthHandler.SignOutV1)
	// }
}
