package models

import (
	"gorm.io/gorm"
	"time"
)

type Trial struct {
	ID          string `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	CompletedAt time.Time
	Model       *Model       `json:"model"`
	ModelID     uint         `json:"-"`
	Inputs      []TrialInput `json:"inputs"`
	Result      string       `json:"result"`
}

type TrialInput struct {
	gorm.Model
	TrialID string `json:"-"`
	URL     string `json:"url"`
}
