package db

import (
	"api/db/models"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSetupDbReturnsErrorForMissingSqliteConfiguration(t *testing.T) {
	os.Setenv("DB_HOST", "test.db")

	_, err := OpenDb()
	configError, ok := err.(*ConfigurationError)
	assert.True(t, ok, "OpenDb() should return a ConfigurationError")
	assert.Equal(t, "DB_DRIVER", configError.Field)
	assert.Equal(t, "missing configuration environment variable", configError.Err.Error())

	os.Clearenv()
	os.Setenv("DB_DRIVER", "sqlite")

	_, err = OpenDb()
	configError, ok = err.(*ConfigurationError)
	assert.True(t, ok, "OpenDb() should return a ConfigurationError")
	assert.Equal(t, "DB_HOST", configError.Field)
	assert.Equal(t, "missing configuration environment variable", configError.Err.Error())
}

func TestOpenDbReturnsErrorForUnsupportedDatabaseType(t *testing.T) {
	os.Setenv("DB_DRIVER", "bad_db")
	os.Setenv("DB_HOST", "test")

	_, err := OpenDb()
	configError, ok := err.(*ConfigurationError)
	assert.True(t, ok, "OpenDb() should return a ConfigurationError")
	assert.Equal(t, "DB_DRIVER", configError.Field)
	assert.Equal(t, "unsupported database: bad_db", configError.Err.Error())
}

//func TestOpenDbReturnsADatabase(t *testing.T) {
//	os.Setenv("DB_DRIVER", "sqlite")
//	os.Setenv("DB_HOST", "test.sqlite")
//
//	db, err := OpenDb()
//	Migrate(db)
//
//	assert.Nil(t, err)
//}

func TestMarshalModelToJson(t *testing.T) {
	m := &models.Model{
		Attributes: models.ModelAttributes{
			Top1:            "54.92",
			Top5:            "78.03",
			Kind:            "CNN",
			ManifestAuthor:  "Cheng Li",
			TrainingDataset: "ImageNet",
		},
		Description: "MXNet Image Classification model, which is trained on the ImageNet dataset. Use AlexNet from GluonCV model zoo.",
		Details: models.ModelDetails{
			GraphChecksum:   "4abd57ec8863ff3e3e29ecd4ead43d1f",
			GraphPath:       "model-symbol.json",
			WeightsChecksum: "906234b2a6b14bedac2dcccba8178529",
			WeightsPath:     "model-0000.params",
		},
		Framework: models.Framework{
			Name:    "MXNet",
			Version: "1.7.0",
		},
		FrameworkID: 0,
		Input: models.ModelOutput{
			Description: "the input image",
			Type:        "image",
		},
		License: "unrestricted",
		Name:    "AlexNet",
		Output: models.ModelOutput{
			Description: "the output label",
			Type:        "classification",
		},
		Version: "1.0",
	}

	j, err := json.Marshal(m)

	assert.Nil(t, err)
	println(string(j))
}
