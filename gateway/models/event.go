package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name                string    `gorm:"not null" json:"name"`
	StartDateTime       time.Time `gorm:"type:timestamp without time zone; not null" json:"start"`
	EndDateTime         time.Time `gorm:"type:timestamp without time zone; not null" json:"end"`
	Location            string    `gorm:"not null" json:"location"`
	ApplicationDeadline time.Time `gorm:"type:date; not null" json:"deadline"`
	Cover               string    `gorm:"not null" json:"cover"`
	Description         string    `gorm:"not null" json:"desc"`
}
