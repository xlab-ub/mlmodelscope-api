package migrations

import (
	"api/db/models"
	"gorm.io/gorm"
)

func CreateExperimentsTable(db *gorm.DB) (err error) {
	type Experiment struct {
		gorm.Model
		ID string `gorm:"primaryKey"`
	}

	db.Migrator().CreateTable(&Experiment{})
	db.Migrator().AddColumn(&models.Trial{}, "experiment_id")
	return db.Migrator().CreateConstraint(&models.Trial{}, "Experiment")
}