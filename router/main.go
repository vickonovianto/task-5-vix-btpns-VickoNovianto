package router

import (
	"log"
	"os"

	"user-photo-api/middlewares"

	"github.com/gin-gonic/gin"
)

func MountRoutes(server *gin.Engine) {
	api := server.Group(os.Getenv("API_PREFIX"))
	authMiddleware, err := middlewares.GetJwtMiddleware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	MountPhotoRoutes(api, authMiddleware)
	MountUserRoutes(api, authMiddleware)

}
