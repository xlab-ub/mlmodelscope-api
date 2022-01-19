package migrations

import (
	"api/db/models"
	"gorm.io/gorm"
)

func AddSearchIndices(db *gorm.DB) (err error) {
	type Model struct {
		Name        string `gorm:"index:idx_models_name,expression:LOWER(name)"`
		Description string `gorm:"index:idx_models_description,expression:LOWER(description)"`
	}

	if !db.Migrator().HasIndex(&models.Model{}, "Name") {
		err = db.Migrator().CreateIndex(&Model{}, "Name")
	}

	if err == nil && !db.Migrator().HasIndex(&models.Model{}, "Description") {
		err = db.Migrator().CreateIndex(&Model{}, "Description")
	}

	return
}
