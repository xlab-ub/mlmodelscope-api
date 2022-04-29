// +build !integration

package db

import (
	"api/db/models"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testDb *Db

func TestSetupDbReturnsErrorForMissingSqliteConfiguration(t *testing.T) {
	os.Setenv("DB_HOST", "test.db")

	_, err := OpenDb()
	configError, ok := err.(*ConfigurationError)
	assert.True(t, ok, "OpenDb() should return a ConfigurationError")
	assert.Equal(t, "DB_DRIVER", configError.Field)
	assert.Equal(t, "missing configuration environment variable: DB_DRIVER", configError.Error())

	os.Clearenv()
	os.Setenv("DB_DRIVER", "sqlite")

	_, err = OpenDb()
	configError, ok = err.(*ConfigurationError)
	assert.True(t, ok, "OpenDb() should return a ConfigurationError")
	assert.Equal(t, "DB_HOST", configError.Field)
	assert.Equal(t, "missing configuration environment variable: DB_HOST", configError.Error())
}

func TestOpenDbReturnsErrorForMissingPostgresConfiguration(t *testing.T) {
	for _, missing := range []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_DBNAME",
	} {
		t.Run(fmt.Sprintf("Missing_%s", missing), func(t *testing.T) {
			setPostgresEnvVars()
			os.Setenv(missing, "")

			_, err := OpenDb()
			configError, ok := err.(*ConfigurationError)
			assert.True(t, ok, "OpenDb() should return a ConfigurationError")
			assert.Equal(t, missing, configError.Field)
		})
	}
}

func setPostgresEnvVars() {
	os.Setenv("DB_DRIVER", "postgres")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "1234")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_DBNAME", "database")
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

func TestQueryEmptyDb(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	models, _ := testDb.GetAllModels()

	assert.Equal(t, 0, len(models))
}

func TestCreateAndQueryModels(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()

	createModelNamed("test1")
	createModelNamed("test2")

	models, _ := testDb.GetAllModels()

	assert.Equal(t, 2, len(models))
	assert.Equal(t, "test1", models[0].Name)
	assert.Equal(t, "test2", models[1].Name)
}

func createModelNamed(name string) {
	testDb.CreateModel(&models.Model{Name: name})
}

func TestCreateAndQueryFrameworks(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createFrameworkNamed("fw1")
	createFrameworkNamed("fw2")

	frameworks, _ := testDb.GetAllFrameworks()

	assert.Equal(t, 2, len(frameworks))
	assert.Equal(t, "fw1", frameworks[0].Name)
	assert.Equal(t, "amd64", frameworks[0].Architectures[0].Name)
	assert.Equal(t, "fw2", frameworks[1].Name)
	assert.Equal(t, "amd64", frameworks[1].Architectures[0].Name)
}

func TestQueryModelById(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createModelNamed("modelById")
	createModelNamed("modelById2")

	result, _ := testDb.GetModelById(2)

	assert.Equal(t, "modelById2", result.Name)
}

func TestQueryModelByIdIncludesFrameworkAndArchitectures(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createFramework("fw1", "arch1")
	testDb.CreateModel(&models.Model{Name: "model1", FrameworkID: 1})

	result, _ := testDb.GetModelById(1)

	assert.Equal(t, "model1", result.Name)
	assert.Equal(t, "fw1", result.Framework.Name)
	assert.Equal(t, "arch1", result.Framework.Architectures[0].Name)
}

func TestQueryModelsByFrameworkId(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createFrameworkNamed("fw1")
	createFrameworkNamed("fw2")
	testDb.CreateModel(&models.Model{Name: "model1", FrameworkID: 1})
	testDb.CreateModel(&models.Model{Name: "model2", FrameworkID: 2})

	result, _ := testDb.QueryModels(1, "", "", "")

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "model1", result[0].Name)
	assert.Equal(t, "fw1", result[0].Framework.Name)
}

func TestQueryModelsByTask(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createModelNamed("model1")
	testDb.CreateModel(&models.Model{Name: "model2", Output: models.ModelOutput{Type: "classification"}})

	result, _ := testDb.QueryModels(0, "classification", "", "")

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "model2", result[0].Name)
}

func TestQueryModelsByUnknownTask(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createModelNamed("model1")

	result, _ := testDb.QueryModels(0, "classification", "", "")

	assert.Equal(t, 0, len(result))
}

func TestQueryModelsByArchitecture(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createFrameworkNamed("fw1")
	testDb.CreateFramework(&models.Framework{
		Name: "fw2",
		Architectures: []models.Architecture{
			models.Architecture{Name: "arm"},
		},
	})
	testDb.CreateModel(&models.Model{Name: "model1", FrameworkID: 1})
	testDb.CreateModel(&models.Model{Name: "model2", FrameworkID: 2})

	result, _ := testDb.QueryModels(0, "", "arm", "")

	assert.Equal(t, 1, len(result))
	assert.Equal(t, 1, len(result[0].Framework.Architectures))
	assert.Equal(t, "arm", result[0].Framework.Architectures[0].Name)
}

func TestQueryModelsByFrameworkAndTask(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createFrameworkNamed("fw1")
	createFrameworkNamed("fw2")
	testDb.CreateModel(&models.Model{Name: "model1", FrameworkID: 1, Output: models.ModelOutput{Type: "classification"}})
	testDb.CreateModel(&models.Model{Name: "model2", FrameworkID: 1, Output: models.ModelOutput{Type: "segmentation"}})
	testDb.CreateModel(&models.Model{Name: "model3", FrameworkID: 2, Output: models.ModelOutput{Type: "classification"}})

	result, _ := testDb.QueryModels(1, "classification", "", "")

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "model1", result[0].Name)

	result, _ = testDb.QueryModels(2, "classification", "", "")

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "model3", result[0].Name)
}

func TestQueryModelsByTaskAndArchitecture(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createFramework("fw1", "amd64")
	createFramework("fw2", "arm")
	testDb.CreateModel(&models.Model{Name: "model1", FrameworkID: 1, Output: models.ModelOutput{Type: "classification"}})
	testDb.CreateModel(&models.Model{Name: "model2", FrameworkID: 2, Output: models.ModelOutput{Type: "segmentation"}})
	testDb.CreateModel(&models.Model{Name: "model3", FrameworkID: 2, Output: models.ModelOutput{Type: "classification"}})

	result, _ := testDb.QueryModels(0, "segmentation", "arm", "")

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "model2", result[0].Name)
}

func TestQueryModelsBySearchString(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createFramework("MXNet", "amd64")
	createFramework("Onnxruntime", "arm")
	testDb.CreateModel(&models.Model{Name: "AlexNet", FrameworkID: 1, Description: "nothing"})
	testDb.CreateModel(&models.Model{Name: "Inception_v3", FrameworkID: 2, Description: ""})
	testDb.CreateModel(&models.Model{Name: "Xception", FrameworkID: 2, Description: "target"})

	result, _ := testDb.QueryModels(0, "", "", "leX")

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "AlexNet", result[0].Name)

	result, _ = testDb.QueryModels(0, "", "", "target")

	assert.Equal(t, 1, len(result))
	assert.Equal(t, uint(3), result[0].ID)
}

func TestQueryModelsExcludesDeleted(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createFrameworkNamed("Pytorch")
	testDb.CreateModel(&models.Model{Name: "AlexNet", FrameworkID: 1, Description: ""})
	testDb.CreateModel(&models.Model{Name: "Inception_v3", FrameworkID: 1, Description: ""})
	testDb.DeleteModel(&models.Model{ID: 1})

	result, _ := testDb.QueryModels(0, "", "", "")

	assert.Equal(t, 1, len(result))
}

func TestQueryFrameworksByNameAndVersion(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createFrameworkNamed("Pytorch")
	testDb.CreateFramework(&models.Framework{Name: "Pytorch", Version: "1.0"})

	result, _ := testDb.QueryFrameworks(&models.Framework{Name: "Pytorch", Version: "1.0"})

	assert.Equal(t, uint(2), result.ID)
}

func createFramework(name string, architecture string) {
	testDb.CreateFramework(&models.Framework{
		Name: name,
		Architectures: []models.Architecture{
			models.Architecture{Name: architecture},
		},
	})
}

func createFrameworkNamed(name string) {
	testDb.CreateFramework(&models.Framework{
		Name: name,
		Architectures: []models.Architecture{
			models.Architecture{Name: "amd64"},
		},
	})
}

func cleanupTestDatabase() {
	os.Remove("test.sqlite")
}

func CreateTestDatabase() {
	os.Setenv("DB_DRIVER", "sqlite")
	os.Setenv("DB_HOST", "test.sqlite")

	testDb, _ = OpenDb()
	testDb.Migrate()
}
