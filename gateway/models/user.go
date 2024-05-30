package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email         string `gorm:"not null;unique" json:"email"`
	Password      string `gorm:"not null" json:"password"`
	PhoneNumber   int    `json:"phone_number"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	AcademicGroup string `json:"academic_group"`
	Admin         bool   `gorm:"not null" json:"admin"`
	Verified      bool   `gorm:"not null" json:"verified"`
}

type UserInfo struct {
	Email         string `json:"email"`
	PhoneNumber   int    `json:"phone_number"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	AcademicGroup string `json:"academic_group"`
}
