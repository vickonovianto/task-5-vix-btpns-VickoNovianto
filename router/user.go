package router

import (
	"user-photo-api/controllers"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func MountUserRoutes(api *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	userController := controllers.NewUserController()
	userGroup := api.Group("/users")
	userGroup.POST("/register", userController.RegisterUserHandler)
	userGroup.POST("/login", authMiddleware.LoginHandler)
	userGroup.PUT("/:userId", authMiddleware.MiddlewareFunc(), userController.UpdateUserHandler)
	userGroup.DELETE("/:userId", authMiddleware.MiddlewareFunc(), userController.DeleteUserHandler)
}
