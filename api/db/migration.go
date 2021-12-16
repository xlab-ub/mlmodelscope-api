package db

import (
	"api/db/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) (err error) {
	err = db.AutoMigrate(&models.Model{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&models.Framework{})
	if err != nil {
		return
	}

	return
}
