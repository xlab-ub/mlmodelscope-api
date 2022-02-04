package migrations

import "gorm.io/gorm"

func CreateUsersTable(db *gorm.DB) (err error) {
	type User struct {
		gorm.Model
		ID string `gorm:"primaryKey"`
	}

	return db.Migrator().CreateTable(&User{})
}
