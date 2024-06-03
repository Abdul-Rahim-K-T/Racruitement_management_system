package models

import "github.com/jinzhu/gorm"

type Profile struct {
	gorm.Model
	// ID            uint   `gorm:"primary_key"`
	UserID        uint   `json:"user_id" gorm:"not null;unique"`
	ResumeFileURL string `json:"resume_file_url"`
	Skills        string `json:"skills"`
	Education     string `json:"education"`
	Experience    string `json:"experience"`
	Name          string `json:"name"`  // Assuming this is the name of the applicant
	Email         string `json:"email"` // Assuming this is the email of the applicant
	Phone         string `json:"phone"` // Assuming this is the phone number of the applicant
}
