package migrations

import (
	"api/db/models"
	"gorm.io/gorm"
)

func AddUserToTrialInputs(db *gorm.DB) (err error) {
	db.Migrator().AddColumn(&models.TrialInput{}, "user_id")
	return db.Migrator().CreateConstraint(&models.TrialInput{}, "User")
}
