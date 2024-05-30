package models

import "gorm.io/gorm"

type Content struct {
	gorm.Model
	Type    string `gorm:"not null" json:"type"`
	Name    string `gorm:"not null" json:"name"`
	Authors string `gorm:"not null" json:"authors"`
	Cover   string `gorm:"not null" json:"cover"`
	Text    string `gorm:"not null" json:"text"`
	Views   int    `gorm:"not null" json:"views"`
	Images  string `gorm:"not null" json:"images"`
}
