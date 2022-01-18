package db

import (
	"api/db/models"
	"fmt"
)

type ModelInteractor interface {
	CreateModel(*models.Model) error
	GetAllModels() ([]models.Model, error)
	GetModelsByTask(string) ([]models.Model, error)
	GetModelsForFramework(int) ([]models.Model, error)
	QueryModels(uint, string, string) ([]models.Model, error)
}

func (d *Db) CreateModel(m *models.Model) (err error) {
	d.database.Create(m)

	return
}

func (d *Db) GetAllModels() (m []models.Model, err error) {
	d.database.Joins("Framework").Find(&m)

	return
}

func (d *Db) GetModelById(id uint) (m *models.Model, err error) {
	d.database.Joins("Framework").First(&m, id)

	if m.ID != id {
		return nil, fmt.Errorf("Unknown Model Id: %d", id)
	}

	return
}

func (d *Db) GetModelsByTask(task string) (m []models.Model, err error) {
	d.database.Where(&models.Model{Output: models.ModelOutput{Type: task}}).Joins("Framework").Find(&m)

	return
}

func (d *Db) GetModelsForFramework(frameworkId int) (m []models.Model, err error) {
	d.database.Where(&models.Model{FrameworkID: frameworkId}).Joins("Framework").Find(&m)

	return
}

func (d *Db) QueryModels(frameworkId int, task string, architecture string) (m []models.Model, err error) {
	where := make(map[string]interface{})

	if frameworkId > 0 {
		where["models.framework_id"] = frameworkId
	}

	if task != "" {
		where["models.output_type"] = task
	}

	if architecture != "" {
		where["architectures.name"] = architecture
	}

	d.database.
		Joins("Framework").
		Preload("Framework.Architectures").
		Joins("LEFT JOIN architectures ON architectures.framework_id = \"Framework\".id").
		Where(where).
		Find(&m)
	return
}
