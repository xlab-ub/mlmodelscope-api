package migrations

import (
	"api/db/models"
	"gorm.io/gorm"
)

func CreateFrameworksTable(db *gorm.DB) (err error) {
	type Framework struct {
		gorm.Model
	}

	db.Migrator().CreateTable(&Framework{})
	db.Migrator().AddColumn(&models.Framework{}, "name")
	return db.Migrator().AddColumn(&models.Framework{}, "version")
}
