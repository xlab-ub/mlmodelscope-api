package db

import (
	"api/db/models"
)

type Migrator interface {
	Migrate() error
}

func (d *Db) Migrate() (err error) {
	err = d.database.AutoMigrate(&models.Model{})
	if err != nil {
		return
	}

	err = d.database.AutoMigrate(&models.Framework{})

	return
}
