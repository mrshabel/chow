package router

import (
	"chow/docs"
	"chow/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(router *gin.Engine, authHandler *handler.AuthHandler, jointHandler *handler.JointHandler, complaintHandler *handler.ComplaintHandler, middleware *handler.Middleware) {
	// cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "X-Forwarded-For", "Origin", "Content-Type", "Content-Length"},
		AllowCredentials: true,
	}))

	// base api router
	apiRouter := router.Group("/api")
	apiRouter.GET("/health", func(ctx *gin.Context) { ctx.JSON(200, "OK") })

	// swagger
	docs.SwaggerInfo.BasePath = "/api"
	apiRouter.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// auth
	auth := apiRouter.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}

	// joints
	joints := apiRouter.Group("/joints")
	{
		// public
		joints.GET("", jointHandler.GetAllJoints)
		joints.GET("/nearby", jointHandler.GetNearByJoints)
		joints.GET("/search", jointHandler.SearchJoints)
		joints.GET("/:id", jointHandler.GetJoint)

		// protected
		protectedJoints := joints.Use(middleware.AuthMiddleware())
		{
			protectedJoints.POST("", jointHandler.CreateJoint)
			protectedJoints.PATCH("/:id", jointHandler.UpdateJoint)
			protectedJoints.DELETE("/:id", jointHandler.DeleteJoint)
			protectedJoints.POST("/:id/vote", jointHandler.VoteJoint)
			protectedJoints.POST("/:id/complaints", jointHandler.CreateJointComplaint)
			protectedJoints.GET("/:id/complaints", jointHandler.GetJointComplaints)
		}
	}

	// complaints
	complaints := apiRouter.Group("/complaints")
	{
		// protected
		protectedComplaints := complaints.Use(middleware.AuthMiddleware())
		{
			protectedComplaints.GET("", complaintHandler.GetAllComplaints)
			protectedComplaints.GET("/me", complaintHandler.GetUserComplaints)
			protectedComplaints.GET("/:id", complaintHandler.GetComplaint)
			protectedComplaints.PATCH("/:id/resolve", complaintHandler.ResolveComplaint)
		}
	}
}
