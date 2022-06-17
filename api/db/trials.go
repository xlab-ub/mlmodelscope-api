package db

import (
	"api/db/models"
	"fmt"
	"time"
)

type TrialInteractor interface {
	CompleteTrial(*models.Trial, string) error
	CreateTrial(*models.Trial) error
	DeleteTrial(id string) error
	GetAllTrials() ([]models.Trial, error)
	GetTrialById(id string) (*models.Trial, error)
	GetTrialByModelAndInput(modelId uint, inputUrl string) (*models.Trial, error)
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

	if trial.ExperimentID == "" {
		return fmt.Errorf("attempt to create Trial without an associated Experiment")
	}

	return d.database.Create(trial).Error
}

func (d *Db) DeleteTrial(id string) error {
	if trial, err := d.GetTrialById(id); err != nil {
		return err
	} else {
		if experiment, err := d.GetExperimentById(trial.ExperimentID); err != nil {
			return err
		} else {
			if len(experiment.Trials) == 1 {
				return fmt.Errorf("DeleteTrial: Experiment must have at least one Trial")
			}
		}
	}
	return d.database.Delete(&models.Trial{ID: id}).Error
}

func (d *Db) GetAllTrials() (trials []models.Trial, err error) {
	err = d.database.Preload("Inputs").Joins("Model").Find(&trials).Error

	return
}

func (d *Db) GetTrialById(id string) (trial *models.Trial, err error) {
	err = d.database.Preload("Inputs").Joins("Experiment").Joins("Model").First(&trial, "trials.id = ?", id).Error

	if err != nil {
		err = fmt.Errorf("unknown Trial: %s", id)
		return nil, err
	}

	return
}

func (d *Db) GetTrialByModelAndInput(modelId uint, inputUrl string) (trial *models.Trial, err error) {
	inputQuery := d.database.Select("trial_id").
		Where("url = ?", inputUrl).
		Table("trial_inputs")

	err = d.database.
		Preload("Inputs").
		Joins("Experiment").
		Joins("Model").
		Where("trials.model_id = ? AND trials.id IN (?)", modelId, inputQuery).
		First(&trial).Error

	if err != nil {
		err = fmt.Errorf("error querying trial with (Model: %d, InputUrl: %s)", modelId, inputUrl)
		return nil, err
	}

	return
}
