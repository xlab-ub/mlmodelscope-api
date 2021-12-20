package db

import (
	"api/db/models"
)

func (d *db) Migrate() (err error) {
	err = d.database.AutoMigrate(&models.Model{})
	if err != nil {
		return
	}

	err = d.database.AutoMigrate(&models.Framework{})

	return
}
