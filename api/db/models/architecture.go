package models

import "time"

type Architecture struct {
	ID          uint       `gorm:"primaryKey" json:"-"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	Name        string     `json:"name"`
	FrameworkID uint        `json:"-"`
}
