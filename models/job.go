package models

type Job struct {
	ID                uint   `gorm:"primary_key" json:"id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	PostedOn          string `json:"posted_on"`
	TotalApplications int    `json:"total_applications"`
	CompanyName       string `json:"company_name"`
	PostedByID        uint   `json:"posted_by_id"`                           // Foreign key for User
	PostedBy          User   `gorm:"foreignKey:PostedByID" json:"posted_by"` // Reference to the user who posted the job  // Reference to the user who posted the job
}
