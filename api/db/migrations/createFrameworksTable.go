package migrations

import (
	"api/db/models"
	"gorm.io/gorm"
)

func CreateFrameworksTable(db *gorm.DB) (err error) {
	type Framework struct {
		gorm.Model
	}

	if err = db.Migrator().CreateTable(&Framework{}); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Framework{}, "name"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Framework{}, "version"); err != nil {
		return
	}

	return
}

