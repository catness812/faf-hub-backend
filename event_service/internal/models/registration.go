package models

import "gorm.io/gorm"

type Registration struct {
	gorm.Model
	EventID         uint   `gorm:"not null" json:"event_id"`
	UserID          uint   `gorm:"not null" json:"user_id"`
	FirstName       string `gorm:"not null" json:"first_name"`
	LastName        string `gorm:"not null" json:"last_name"`
	Email           string `gorm:"not null" json:"email"`
	PhoneNumber     int    `gorm:"not null" json:"phone_number"`
	AcademicGroup   string `gorm:"not null" json:"academic_group"`
	TeamMembers     string `gorm:"not null" json:"team_members"`
	ShirtSize       string `gorm:"not null" json:"shirt_size"`
	FoodPreferences string `gorm:"not null" json:"food_pref"`
	Motivation      string `gorm:"not null" json:"motivation"`
	Questions       string `gorm:"not null" json:"questions"`
}
