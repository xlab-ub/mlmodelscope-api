package db

import (
	"api/db/models"
	"time"
)

type Migrator interface {
	Migrate() error
}

type migrations struct {
	Migration  int
	MigratedAt time.Time
}

func (d *Db) Migrate() (err error) {
	nextMigration := lookupNextMigration(d)

	for index, migrator := range migrators[nextMigration:] {
		err = migrator(d)

		if err != nil {
			return
		}

		recordMigration(d, nextMigration + index + 1)
	}

	return
}

func lookupNextMigration(d *Db) (lastMigration int) {
	if d.database.Migrator().HasTable(&migrations{}) {
		m := migrations{}
		d.database.Order("migration desc").First(&m)
		lastMigration = m.Migration
	} else {
		lastMigration = 0
	}

	return
}

func recordMigration(d *Db, number int) {
	d.database.Create(&migrations{
		Migration: number,
		MigratedAt: time.Now(),
	})
}

func CreateMigrationsTable(db *Db) (err error) {
	if !db.database.Migrator().HasTable(&migrations{}) {
		err = db.database.Migrator().CreateTable(&migrations{})
	}

	return
}

func CreateFrameworksTable(db *Db) (err error) {
	return db.database.Migrator().CreateTable(&models.Framework{})
}

func CreateModelsTable(db *Db) error {
	return db.database.Migrator().CreateTable(&models.Model{})
}

// ONLY append new migrator functions to the end of this list
// NEVER remove any migrator function from this list that has
//   already been run against a database that can't be recreated
//   from scratch (for example the Staging and Production databases)
var migrators = [](func(*Db) error) {
	CreateMigrationsTable,
	CreateFrameworksTable,
	CreateModelsTable,
}
