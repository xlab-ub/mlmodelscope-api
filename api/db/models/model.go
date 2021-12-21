package models

import (
	"time"
)

type Model struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time       `json:"-"`
	UpdatedAt   time.Time       `json:"-"`
	Attributes  ModelAttributes `gorm:"embedded;embeddedPrefix:attribute_" json:"attributes"`
	Description string          `json:"description"`
	Details     ModelDetails    `gorm:"embedded;embeddedPrefix:detail_" json:"model"`
	Framework   Framework       `json:"framework"`
	FrameworkID int             `json:"-"`
	Input       ModelOutput     `gorm:"embedded;embeddedPrefix:input_" json:"input"`
	License     string          `json:"license"`
	Name        string          `json:"name"`
	Output      ModelOutput     `gorm:"embedded;embeddedPrefix:output_" json:"output"`
	Version     string          `json:"version"`
}

type ModelAttributes struct {
	Top1            string
	Top5            string
	Kind            string `json:"kind"`
	ManifestAuthor  string `json:"manifest_author"`
	TrainingDataset string `json:"training_dataset"`
}

type ModelDetails struct {
	GraphChecksum   string `json:"graph_checksum"`
	GraphPath       string `json:"graph_path"`
	WeightsChecksum string `json:"weights_checksum"`
	WeightsPath     string `json:"weights_path"`
}

type ModelOutput struct {
	Description string `json:"description"`
	Type        string `json:"type"`
}
