package models

import "gorm.io/gorm"

// User model definition
type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email" gorm:"uniqueIndex"`
}
