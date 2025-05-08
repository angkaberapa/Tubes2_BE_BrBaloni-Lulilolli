package v1

import (
	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/api"
	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/core"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the API routes for the v1 version of the API.
func RegisterRoutes(r *gin.Engine, handlers *api.Handlers, appCtx *core.AppContext) {
	// v1 := r.Group("/v1")

	// authGroup := v1.Group("/auth")
	// {
	// 	authGroup.GET("/signin", handlers.AuthHandler.SignInV1)
	// 	authGroup.GET("/signout", handlers.AuthHandler.SignOutV1)
	// }
}
