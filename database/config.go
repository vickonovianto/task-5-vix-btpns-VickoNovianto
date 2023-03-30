package database

import (
	"user-photo-api/database/mysql"

	"gorm.io/gorm"
)

func Database() *gorm.DB {
	return mysql.InitGorm()
}
