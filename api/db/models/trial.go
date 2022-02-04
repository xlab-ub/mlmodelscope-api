package models

import (
	"gorm.io/gorm"
	"time"
)

type Trial struct {
	ID           string         `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	CompletedAt  *time.Time     `json:"completed_at,omitempty"`
	Experiment   *Experiment    `json:"experiment,omitempty"`
	ExperimentID string         `json:"-"`
	Model        *Model         `json:"model,omitempty"`
	ModelID      uint           `json:"-"`
	Inputs       []TrialInput   `json:"inputs,omitempty"`
	Result       string         `json:"result,omitempty"`
}

type TrialInput struct {
	gorm.Model
	TrialID string `json:"-"`
	URL     string `json:"url"`
	User    *User  `json:"-"`
	UserID  string `json:"user_id"`
}
