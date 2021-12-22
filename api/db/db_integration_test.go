// +build integration

package db

import (
	"api/db/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatabase(t *testing.T) {
	var (
		d   Db
		err error
	)

	t.Run("Connect", func(t *testing.T) {
		d, err = OpenDb()

		assert.Nil(t, err)
	})

	t.Run("Migrate", func(t *testing.T) {
		err = d.Migrate()

		assert.Nil(t, err)
	})

	t.Run("Create Model", func(t *testing.T) {
		err = d.CreateModel(&models.Model{Name: "integrate", Framework: &models.Framework{Name: "integrate"}})

		assert.Nil(t, err)
	})

	t.Run("Query Model", func(t *testing.T) {
		m, err := d.GetAllModels()

		assert.Nil(t, err)
		assert.Equal(t, "integrate", m[0].Name)
		assert.Equal(t, "integrate", m[0].Framework.Name)
	})
}
