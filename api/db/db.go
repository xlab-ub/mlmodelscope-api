package db

import (
	"api/db/models"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
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
	dbName    string
	driver    string
	dialector gorm.Dialector
	host      string
	password  string
	port      string
	user      string
)

type Db interface {
	CreateFramework(*models.Framework) error
	CreateModel(*models.Model) error
	GetAllFrameworks() ([]models.Framework, error)
	GetAllModels() ([]models.Model, error)
	GetModelsByTask(string) ([]models.Model, error)
	GetModelsForFramework(int) ([]models.Model, error)
	Migrate() error
	QueryFrameworks(*models.Framework) (*models.Framework, error)
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

func (d *db) GetModelsByTask(task string) (m []models.Model, err error) {
	d.database.Where(&models.Model{Output: models.ModelOutput{Type: task}}).Joins("Framework").Find(&m)

	return
}

func (d *db) GetModelsForFramework(frameworkId int) (m []models.Model, err error) {
	d.database.Where(&models.Model{FrameworkID: frameworkId}).Joins("Framework").Find(&m)

	return
}

func (d *db) QueryFrameworks(query *models.Framework) (*models.Framework, error) {
	var framework models.Framework
	r := d.database.Where(query).First(&framework)

	if r.Error != nil {
		return nil, nil
	}

	return &framework, nil
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
		return
	}

	host = os.Getenv("DB_HOST")
	if host == "" {
		err = &ConfigurationError{
			Field: "DB_HOST",
			Err:   errors.New("missing configuration environment variable"),
		}
		return
	}

	if driver == "postgres" {
		return readServerConfiguration()
	}

	return
}

func readServerConfiguration() (err error) {
	dbName = os.Getenv("DB_DBNAME")
	if dbName == "" {
		err = &ConfigurationError{
			Field: "DB_DBNAME",
			Err:   errors.New("missing configuration environment variable"),
		}
		return
	}

	password = os.Getenv("DB_PASSWORD")
	if password == "" {
		err = &ConfigurationError{
			Field: "DB_PASSWORD",
			Err:   errors.New("missing configuration environment variable"),
		}
		return
	}

	port = os.Getenv("DB_PORT")
	if port == "" {
		err = &ConfigurationError{
			Field: "DB_PORT",
			Err:   errors.New("missing configuration environment variable"),
		}
		return
	}

	user = os.Getenv("DB_USER")
	if user == "" {
		err = &ConfigurationError{
			Field: "DB_USER",
			Err:   errors.New("missing configuration environment variable"),
		}
		return
	}

	return
}

func getDriver() (d gorm.Dialector, err error) {
	switch driver {
	case "sqlite":
		d = sqlite.Open(host)
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)
		d = postgres.Open(dsn)
	default:
		err = &ConfigurationError{
			Field: "DB_DRIVER",
			Err:   fmt.Errorf("unsupported database: %s", driver),
		}
	}

	return
}
