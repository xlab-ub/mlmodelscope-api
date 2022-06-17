package migrations

import (
	"api/db/models"
	"gorm.io/gorm"
)

func AddSourceTrialToTrial(db *gorm.DB) (err error) {
	db.Migrator().AddColumn(&models.Trial{}, "source_trial_id")
	return db.Migrator().CreateConstraint(&models.Trial{}, "SourceTrial")
}
