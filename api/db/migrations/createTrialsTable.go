package migrations

import (
	"api/db/models"
	"gorm.io/gorm"
)

func CreateTrialsTable(db *gorm.DB) (err error) {
	type Trial struct {
		gorm.Model
		ID string `gorm:"primaryKey"`
	}

	type TrialInput struct {
		gorm.Model
	}

	db.Migrator().CreateTable(&Trial{})
	db.Migrator().AddColumn(&models.Trial{}, "model_id")
	db.Migrator().CreateConstraint(&models.Trial{}, "Model")
	db.Migrator().AddColumn(&models.Trial{}, "completed_at")
	db.Migrator().AddColumn(&models.Trial{}, "result")
	db.Migrator().CreateTable(&TrialInput{})
	db.Migrator().AddColumn(&models.TrialInput{}, "trial_id")
	db.Migrator().CreateConstraint(&models.Trial{}, "Inputs")
	return db.Migrator().AddColumn(&models.TrialInput{}, "url")
}
