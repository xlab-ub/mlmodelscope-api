package models

import (
	"time"
)

type Model struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	Attributes       ModelAttributes `gorm:"embedded;embeddedPrefix:attribute_" json:"attributes"`
	Description      string          `gorm:"index:idx_models_description,expression:LOWER(description)" json:"description"`
	ShortDescription string          `gorm:"index:idx_models_short_description,expression:LOWER(short_description)" json:"short_description"`
	Details          ModelDetails    `gorm:"embedded;embeddedPrefix:detail_" json:"model"`
	Framework        *Framework      `json:"framework"`
	FrameworkID      int             `json:"-"`
	Input            ModelOutput     `gorm:"embedded;embeddedPrefix:input_" json:"input"`
	License          string          `json:"license"`
	Name             string          `gorm:"index:idx_models_name,expression:LOWER(name)" json:"name"`
	Output           ModelOutput     `gorm:"embedded;embeddedPrefix:output_" json:"output"`
	Version          string          `json:"version"`
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
