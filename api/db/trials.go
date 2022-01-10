package db

import (
	"api/db/models"
	"fmt"
	"time"
)

type TrialInteractor interface {
	CompleteTrial(*models.Trial, string) error
	CreateTrial(*models.Trial) error
	GetAllTrials() ([]models.Trial, error)
	GetTrialById(id string) (*models.Trial, error)
}

func (d *Db) CompleteTrial(trial *models.Trial, result string) error {
	now := time.Now()
	d.database.Model(&trial).Updates(models.Trial{CompletedAt: &now, Result: result})

	return d.database.Error
}

func (d *Db) CreateTrial(trial *models.Trial) (err error) {
	if trial.ID == "" {
		return fmt.Errorf("attempt to create Trial without an ID")
	}

	if trial.ModelID == 0 {
		return fmt.Errorf("attempt to create Trial without an associated Model")
	}

	d.database.Create(trial)

	return d.database.Error
}

func (d *Db) GetAllTrials() (trials []models.Trial, err error) {
	d.database.Preload("Inputs").Joins("Model").Find(&trials)

	return trials, d.database.Error
}

func (d *Db) GetTrialById(id string) (trial *models.Trial, err error) {
	d.database.Preload("Inputs").Joins("Model").First(&trial, "trials.id = ?", id)

	if trial.ID == "" {
		return nil, fmt.Errorf("unknown Trial: %s", id)
	}

	return trial, d.database.Error
}
