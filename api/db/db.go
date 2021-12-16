package db

import (
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
	database  gorm.DB
	dialector gorm.Dialector
)

func OpenDb() (db *gorm.DB, err error) {
	err = readConfiguration()
	if err != nil {
		return
	}

	dialector, err = getDriver()
	if err != nil {
		return
	}

	return gorm.Open(dialector)
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
