package models

import "time"

type Application struct {
	ID        uint `gorm:"primary_key"`
	JobID     uint
	UserID    uint
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
