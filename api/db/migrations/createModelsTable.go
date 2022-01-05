package migrations

import (
	"api/db/models"
	"gorm.io/gorm"
)

func CreateModelsTable(db *gorm.DB) (err error) {
	type Model struct {
		gorm.Model
	}

	if err = db.Migrator().CreateTable(&Model{}); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "attribute_top1"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "attribute_top5"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "attribute_kind"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "attribute_manifest_author"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "attribute_training_dataset"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "description"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "detail_graph_checksum"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "detail_graph_path"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "detail_weights_checksum"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "detail_weights_path"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "framework_id"); err != nil {
		return
	}

	if err = db.Migrator().CreateConstraint(&models.Model{}, "Framework"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "input_description"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "input_type"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "license"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "name"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "output_description"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "output_type"); err != nil {
		return
	}

	if err = db.Migrator().AddColumn(&models.Model{}, "version"); err != nil {
		return
	}

	return
}

