package router

import (
	"user-photo-api/controllers"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func MountPhotoRoutes(api *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	photoController := controllers.NewPhotoController()
	photoGroup := api.Group("/photos")
	photoGroup.POST("", authMiddleware.MiddlewareFunc(), photoController.CreatePhotoHandler)
	photoGroup.GET("", authMiddleware.MiddlewareFunc(), photoController.GetPhotoHandler)
	photoGroup.PUT("/:photoId", authMiddleware.MiddlewareFunc(), photoController.UpdatePhotoHandler)
	photoGroup.DELETE("/:photoId", authMiddleware.MiddlewareFunc(), photoController.DeletePhotoHandler)
}
