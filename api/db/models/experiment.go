package models

import (
	"gorm.io/gorm"
	"time"
)

type Experiment struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Trials    []Trial        `json:"trials"`
}
