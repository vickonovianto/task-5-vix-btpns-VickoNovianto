package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"user-photo-api/router"

	"github.com/gin-gonic/gin"
)

type (
	server struct {
		httpServer *gin.Engine
	}

	Server interface {
		Run()
	}
)

func InitServer() Server {
	app := gin.Default()

	// Serving uploads folder
	app.Static("/uploads", "./uploads")

	return &server{
		httpServer: app,
	}
}

func (s *server) Run() {
	rootFolderPath, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}

	uploadFolderPath := filepath.Join(rootFolderPath, "uploads")
	if _, err := os.Stat(uploadFolderPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(uploadFolderPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	photosFolderPath := filepath.Join(uploadFolderPath, "users")
	if _, err := os.Stat(photosFolderPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(photosFolderPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	tempFolderPath := filepath.Join(uploadFolderPath, "temp")
	if _, err := os.Stat(tempFolderPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(tempFolderPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	router.MountRoutes(s.httpServer)

	if err := s.httpServer.Run(fmt.Sprintf(":%v", os.Getenv("PORT"))); err != nil {
		log.Panic(err)
	}
}
