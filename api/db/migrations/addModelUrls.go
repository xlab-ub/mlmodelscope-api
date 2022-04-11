package migrations

import (
	"api/db/models"
	"gorm.io/gorm"
)

func AddModelUrls(db *gorm.DB) (err error) {
	err = db.Migrator().AddColumn(&models.Model{}, "url_github")
	if err != nil {
		return
	}

	err = db.Migrator().AddColumn(&models.Model{}, "url_citation")
	if err != nil {
		return
	}

	err = db.Migrator().AddColumn(&models.Model{}, "url_link1")
	if err != nil {
		return
	}

	return db.Migrator().AddColumn(&models.Model{}, "url_link2")
}
