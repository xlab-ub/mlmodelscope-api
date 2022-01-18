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
	d.database.Preload("Architectures").Find(&frameworks)

	return
}

func (d *Db) QueryFrameworks(query *models.Framework) (*models.Framework, error) {
	var framework models.Framework
	r := d.database.Joins("Architectures").Where(query).First(&framework)

	if r.Error != nil {
		return nil, nil
	}

	return &framework, nil
}
