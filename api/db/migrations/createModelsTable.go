package migrations

import (
	"api/db/models"
	"gorm.io/gorm"
)

func CreateModelsTable(db *gorm.DB) (err error) {
	type Model struct {
		gorm.Model
	}

	db.Migrator().CreateTable(&Model{})
	db.Migrator().AddColumn(&models.Model{}, "attribute_top1")
	db.Migrator().AddColumn(&models.Model{}, "attribute_top5")
	db.Migrator().AddColumn(&models.Model{}, "attribute_kind")
	db.Migrator().AddColumn(&models.Model{}, "attribute_manifest_author")
	db.Migrator().AddColumn(&models.Model{}, "attribute_training_dataset")
	db.Migrator().AddColumn(&models.Model{}, "description")
	db.Migrator().AddColumn(&models.Model{}, "detail_graph_checksum")
	db.Migrator().AddColumn(&models.Model{}, "detail_graph_path")
	db.Migrator().AddColumn(&models.Model{}, "detail_weights_checksum")
	db.Migrator().AddColumn(&models.Model{}, "detail_weights_path")
	db.Migrator().AddColumn(&models.Model{}, "framework_id")
	db.Migrator().CreateConstraint(&models.Model{}, "Framework")
	db.Migrator().AddColumn(&models.Model{}, "input_description")
	db.Migrator().AddColumn(&models.Model{}, "input_type")
	db.Migrator().AddColumn(&models.Model{}, "license")
	db.Migrator().AddColumn(&models.Model{}, "name")
	db.Migrator().AddColumn(&models.Model{}, "output_description")
	db.Migrator().AddColumn(&models.Model{}, "output_type")
	return db.Migrator().AddColumn(&models.Model{}, "version")
}
