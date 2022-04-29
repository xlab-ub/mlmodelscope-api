package db

import (
	"api/db/models"
	"database/sql"
	"fmt"
	"strings"
)

type ModelInteractor interface {
	CreateModel(*models.Model) error
	DeleteModel(*models.Model) error
	GetAllModels() ([]models.Model, error)
	GetModelsByTask(string) ([]models.Model, error)
	GetModelsForFramework(int) ([]models.Model, error)
	QueryModels(uint, string, string, string) ([]models.Model, error)
}

func (d *Db) CreateModel(m *models.Model) (err error) {
	return d.database.Create(m).Error
}

func (d *Db) DeleteModel(m *models.Model) error {
	return d.database.Delete(m).Error
}

func (d *Db) GetAllModels() (m []models.Model, err error) {
	err = d.database.Joins("Framework").Find(&m).Error

	return
}

func (d *Db) GetModelById(id uint) (m *models.Model, err error) {
	err = d.database.
		Joins("Framework").
		Preload("Framework.Architectures").
		Joins("LEFT JOIN architectures ON architectures.framework_id = \"Framework\".id").
		First(&m, id).Error

	if err != nil {
		err = fmt.Errorf("Unknown Model Id: %d", id)
	}

	return
}

func (d *Db) GetModelsByTask(task string) (m []models.Model, err error) {
	err = d.database.Where(&models.Model{Output: models.ModelOutput{Type: task}}).Joins("Framework").Find(&m).Error

	return
}

func (d *Db) GetModelsForFramework(frameworkId int) (m []models.Model, err error) {
	err = d.database.Where(&models.Model{FrameworkID: frameworkId}).Joins("Framework").Find(&m).Error

	return
}

func (d *Db) QueryModels(frameworkId int, task string, architecture string, query string) (m []models.Model, err error) {
	db := d.database.
		Joins("Framework").
		Preload("Framework.Architectures").
		Joins("LEFT JOIN architectures ON architectures.framework_id = \"Framework\".id")

	if frameworkId > 0 {
		db = db.Where("models.framework_id = ?", frameworkId)
	}

	if task != "" {
		db = db.Where("models.output_type = ?", task)
	}

	if architecture != "" {
		db = db.Where("architectures.name = ?", architecture)
	}

	if query != "" {
		wildcard := fmt.Sprintf("%%%s%%", strings.ToLower(query))
		db = db.Where("LOWER(models.name) LIKE @query OR LOWER(models.description) LIKE @query", sql.Named("query", wildcard))
	}

	err = db.Find(&m).Error

	return
}
