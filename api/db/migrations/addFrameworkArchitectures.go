package migrations

import (
	"api/db/models"
	"gorm.io/gorm"
)

func AddFrameworkArchitectures(db *gorm.DB) (err error) {
	type Architecture struct {
		gorm.Model
	}

	db.Migrator().CreateTable(&Architecture{})
	db.Migrator().AddColumn(&models.Architecture{}, "name")
	db.Migrator().AddColumn(&models.Architecture{}, "framework_id")
	return db.Migrator().CreateConstraint(&models.Framework{}, "Architectures")
}
