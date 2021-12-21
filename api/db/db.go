package db

import (
	"api/db/models"
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

type ConfigurationError struct {
	Field string
	Err   error
}

func (e *ConfigurationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Field)
}

var (
	host      string
	driver    string
	dialector gorm.Dialector
)

type Db interface {
	CreateFramework(*models.Framework) error
	CreateModel(*models.Model) error
	GetAllFrameworks() ([]models.Framework, error)
	GetAllModels() ([]models.Model, error)
	GetModelsForFramework(int) ([]models.Model, error)
	Migrate() error
}

type db struct {
	database *gorm.DB
}

func (d *db) CreateFramework(f *models.Framework) (err error) {
	d.database.Create(f)

	return
}

func (d *db) CreateModel(m *models.Model) (err error) {
	d.database.Create(m)

	return
}

func (d *db) GetAllFrameworks() (frameworks []models.Framework, err error) {
	d.database.Find(&frameworks)

	return
}

func (d *db) GetAllModels() (m []models.Model, err error) {
	d.database.Joins("Framework").Find(&m)

	return
}

func (d *db) GetModelsForFramework(frameworkId int) (m []models.Model, err error) {
	d.database.Where(&models.Model{FrameworkID: frameworkId}).Joins("Framework").Find(&m)

	return
}

func OpenDb() (result Db, err error) {
	err = readConfiguration()
	if err != nil {
		return
	}

	dialector, err = getDriver()
	if err != nil {
		return
	}

	database, err := gorm.Open(dialector)
	result = &db{
		database: database,
	}
	return
}

func readConfiguration() (err error) {
	driver = os.Getenv("DB_DRIVER")
	if driver == "" {
		err = &ConfigurationError{
			Field: "DB_DRIVER",
			Err:   errors.New("missing configuration environment variable"),
		}
	}

	host = os.Getenv("DB_HOST")
	if host == "" {
		err = &ConfigurationError{
			Field: "DB_HOST",
			Err:   errors.New("missing configuration environment variable"),
		}
	}

	return
}

func getDriver() (d gorm.Dialector, err error) {
	switch driver {
	case "sqlite":
		d = sqlite.Open(host)
	default:
		err = &ConfigurationError{
			Field: "DB_DRIVER",
			Err:   fmt.Errorf("unsupported database: %s", driver),
		}
	}

	return
}
