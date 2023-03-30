package mysql

import (
	"log"
	"os"
	"user-photo-api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitGorm() *gorm.DB {
	connection := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(mysql.Open(connection))
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate(
		&models.User{},
		&models.Photo{},
	)
	return db
}
