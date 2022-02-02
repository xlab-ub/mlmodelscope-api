package db

import (
	"api/db/models"
	"errors"
)

type UserInteractor interface {
	CreateUser(*models.User) error
	GetUserById(id string) (*models.User, error)
}

func (d *Db) CreateUser(user *models.User) (err error) {
	if user.ID == "" {
		return errors.New("attempt to create User without ID")
	}

	return d.database.Create(user).Error
}

func (d *Db) GetUserById(id string) (user *models.User, err error) {
	err = d.database.First(&user, "users.id = ?", id).Error

	return
}
