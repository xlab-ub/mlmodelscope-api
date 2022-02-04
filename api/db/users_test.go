// +build !integration

package db

import (
	"api/db/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserInteractor(t *testing.T) {
	CreateTestDatabase()
	defer cleanupTestDatabase()

	t.Run("CannotCreateUserWithEmptyId", func(t *testing.T) {
		user := &models.User{}
		err := testDb.CreateUser(user)

		assert.Equal(t, "attempt to create User without ID", err.Error())
	})

	t.Run("CreateAndQueryUserById", func(t *testing.T) {
		testDb.CreateUser(&models.User{ID: "testUser"})

		user, err := testDb.GetUserById("testUser")

		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "testUser", user.ID)
	})
}
