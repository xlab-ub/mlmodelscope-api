package migrations

import (
	"api/db/models"
	"gorm.io/gorm"
)

func AddUserToExperiment(db *gorm.DB) (err error) {
	db.Migrator().AddColumn(&models.Experiment{}, "user_id")
	return db.Migrator().CreateConstraint(&models.Experiment{}, "User")
}
