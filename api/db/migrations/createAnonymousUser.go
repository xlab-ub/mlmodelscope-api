package migrations

import (
	"api/db/models"
	"gorm.io/gorm"
)

func CreateAnonymousUser(db *gorm.DB) error {
	return db.Create(&models.User{ID: "anonymous"}).Error
}
