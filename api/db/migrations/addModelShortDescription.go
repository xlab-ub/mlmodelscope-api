package migrations

import (
	"api/db/models"
	"gorm.io/gorm"
)

func AddModelShortDescription(db *gorm.DB) (err error) {
	err = db.Migrator().AddColumn(&models.Model{}, "short_description")
	if err != nil {
		return
	}

	return db.Migrator().CreateIndex(&models.Model{}, "ShortDescription")
}
