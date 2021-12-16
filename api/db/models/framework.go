package models

import "gorm.io/gorm"

type Framework struct {
	gorm.Model `json:"-"`
	Name       string `json:"name"`
	Version    string `json:"version"`
}
