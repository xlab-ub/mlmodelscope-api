package db

import (
	"api/db/models"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testDb Db

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
	assert.Equal(t, "fw2", frameworks[1].Name)
}

func TestQueryModelsByFrameworkId(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createFrameworkNamed("fw1")
	testDb.CreateModel(&models.Model{Name: "model1", FrameworkID: 1})
	createModelNamed("model2")

	result, _ := testDb.GetModelsForFramework(1)

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "model1", result[0].Name)
	assert.Equal(t, "fw1", result[0].Framework.Name)
}

func TestQueryModelsByTask(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createModelNamed("model1")
	testDb.CreateModel(&models.Model{Name: "model2", Output: models.ModelOutput{Type: "classification"}})

	result, _ := testDb.GetModelsByTask("classification")

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "model2", result[0].Name)
}

func TestQueryModelsByUnknownTask(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createModelNamed("model1")

	result, _ := testDb.GetModelsByTask("classification")

	assert.Equal(t, 0, len(result))
}

func TestQueryFrameworksByNameAndVersion(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()
	createFrameworkNamed("Pytorch")
	testDb.CreateFramework(&models.Framework{Name: "Pytorch", Version: "1.0"})

	result, _ := testDb.QueryFrameworks(&models.Framework{Name: "Pytorch", Version: "1.0"})

	assert.Equal(t, uint(2), result.ID)
}

func createFrameworkNamed(name string) {
	testDb.CreateFramework(&models.Framework{Name: name})
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
