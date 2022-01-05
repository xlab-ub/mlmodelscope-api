package migrations

import (
	"gorm.io/gorm"
	"time"
)

type Migrations struct {
	Migration  int
	MigratedAt time.Time
}

func CreateMigrationsTable(db *gorm.DB) (err error) {
	if !db.Migrator().HasTable(&Migrations{}) {
		err = db.Migrator().CreateTable(&Migrations{})
	}

	return
}

