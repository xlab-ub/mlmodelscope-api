package db

import (
	"api/db/models"
	"errors"
)

type ExperimentInteractor interface {
	CreateExperiment(*models.Experiment) error
	GetExperimentById(string) (*models.Experiment, error)
}

func (d *Db) CreateExperiment(experiment *models.Experiment) error {
	if experiment.ID == "" {
		return errors.New("attempt to create Experiment without ID")
	}

	return d.database.Create(experiment).Error
}

func (d *Db) GetExperimentById(id string) (experiment *models.Experiment, err error) {
	err = d.database.Preload("Trials").First(&experiment, "experiments.id = ?", id).Error

	return
}