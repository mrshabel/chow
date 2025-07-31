package router

import (
	"chow/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, authHandler *handler.AuthHandler) {
	// base api router
	apiRouter := router.Group("/api")
	apiRouter.GET("/health", func(ctx *gin.Context) { ctx.JSON(200, "OK") })

	// auth
	auth := apiRouter.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/register", authHandler.Register)
}
