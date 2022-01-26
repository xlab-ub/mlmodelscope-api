package db

import (
	"api/db/migrations"
	"gorm.io/gorm"
	"log"
	"time"
)

type Migrator interface {
	Migrate() error
}

func (d *Db) Migrate() (err error) {
	log.Println("[INFO] running database migrations")
	nextMigration := lookupNextMigration(d)

	for index, migrator := range migrators[nextMigration:] {
		log.Printf("[INFO] running migration %d", nextMigration + index)
		err = migrator(d.database)

		if err != nil {
			return
		}

		recordMigration(d, nextMigration+index+1)
	}

	return
}

func lookupNextMigration(d *Db) (lastMigration int) {
	if d.database.Migrator().HasTable(&migrations.Migrations{}) {
		m := migrations.Migrations{}
		d.database.Order("migration desc").First(&m)
		lastMigration = m.Migration
	} else {
		lastMigration = 0
	}

	return
}

func recordMigration(d *Db, number int) {
	d.database.Create(&migrations.Migrations{
		Migration:  number,
		MigratedAt: time.Now(),
	})
}

// ONLY append new migrator functions to the end of this list
// NEVER remove any migrator function from this list that has
//   already been run against a database that can't be recreated
//   from scratch (for example the Staging and Production databases)
var migrators = [](func(*gorm.DB) error){
	migrations.CreateMigrationsTable,
	migrations.CreateFrameworksTable,
	migrations.CreateModelsTable,
	migrations.CreateTrialsTable,
	migrations.AddFrameworkArchitectures,
	migrations.AddSearchIndices,
	migrations.CreateExperimentsTable,
}
