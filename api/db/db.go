package db

import (
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

type Db struct {
	ExperimentInteractor
	FrameworkInteractor
	ModelInteractor
	TrialInteractor
	Migrator

	database *gorm.DB
}

func OpenDb() (result *Db, err error) {
	err = readConfiguration()
	if err != nil {
		return
	}

	dialector, err = getDriver()
	if err != nil {
		return
	}

	database, err := gorm.Open(dialector)
	result = &Db{
		database: database,
	}
	return
}

func readConfiguration() (err error) {
	if driver, err = readConfigurationVariable("DB_DRIVER"); err != nil {
		return
	}

	if host, err = readConfigurationVariable("DB_HOST"); err != nil {
		return
	}

	if driver == "postgres" {
		return readServerConfiguration()
	}

	return
}

func readServerConfiguration() (err error) {
	if dbName, err = readConfigurationVariable("DB_DBNAME"); err != nil {
		return
	}

	if password, err = readConfigurationVariable("DB_PASSWORD"); err != nil {
		return
	}

	if port, err = readConfigurationVariable("DB_PORT"); err != nil {
		return
	}

	if user, err = readConfigurationVariable("DB_USER"); err != nil {
		return
	}

	return
}

func readConfigurationVariable(name string) (variable string, err error) {
	variable = os.Getenv(name)
	if variable == "" {
		err = &ConfigurationError{
			Field: name,
			Err:   errors.New("missing configuration environment variable"),
		}
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
