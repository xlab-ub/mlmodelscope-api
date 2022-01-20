package db

import "api/db/models"

type FrameworkInteractor interface {
	CreateFramework(*models.Framework) error
	GetAllFrameworks() ([]models.Framework, error)
	QueryFrameworks(*models.Framework) (*models.Framework, error)
}

func (d *Db) CreateFramework(f *models.Framework) (err error) {
	d.database.Create(f)

	return
}

func (d *Db) GetAllFrameworks() (frameworks []models.Framework, err error) {
	err = d.database.Preload("Architectures").Find(&frameworks).Error

	return
}

func (d *Db) QueryFrameworks(query *models.Framework) (framework *models.Framework, err error) {
	err = d.database.Preload("Architectures").Where(query).First(&framework).Error

	return
}
