package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model

	Name         string   `json:"name"`
	Email        string   `json:"email" gorm:"unique_index"`
	Address      string   `json:"address"`
	UserType     string   `json:"user_type"`
	PasswordHash string   `json:"password"`
	ProfileHead  string   `json:"profile_headline"`
	Profile      *Profile `gorm:"foreignKey:UserID" json:"profile,omitempty"` // Profile is a pointer and can be nil for admins
}
