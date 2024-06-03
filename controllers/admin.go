package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/database"
	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/middleware"
	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateJob(c *gin.Context) {
	fmt.Println("Executing Create job function")

	// Get the user object from the context
	UserID, exists := c.Get(string(middleware.ContextKeyUserID))
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Ensure userID is of the correct type (uint)
	userIDUint, ok := UserID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	fmt.Println("currentUser:", userIDUint)
	fmt.Println("OK:", ok)

	// Retrieve the user details from the database to check the UserType
	var dbUser models.User
	if err := database.DB.Preload("Profile").First(&dbUser, userIDUint).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user details from database"})
		return
	}

	fmt.Println("dbUser:", dbUser)

	// Check if the user is an admin
	if dbUser.UserType != "Admin" {
		fmt.Println("usertype:", dbUser.UserType)
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can create jobs"})
		return
	}

	// Parse request body
	var job models.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assign the current admin ID as the job poster ID
	job.PostedByID = dbUser.ID
	job.PostedOn = time.Now().Format("2006-01-02")

	// Save job to database
	if err := database.DB.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	// Preload the PostedBy field and send the job with full user details
	var createdJob models.Job
	if err := database.DB.Preload("PostedBy.Profile").First(&createdJob, job.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created job"})
		return
	}

	// Debug print the createdJob
	fmt.Printf("createdJob: %+v\n", createdJob)

	c.JSON(http.StatusCreated, gin.H{"message": "Job created successfully", "job": createdJob})
}
func ViewJob(c *gin.Context) {
	// Get job ID from the URL parameter
	jobID := c.Param("job_id")

	// Connect to the database
	db := database.InitDB()

	// Query the job from the database
	var job models.Job
	if err := db.First(&job, jobID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	// Return the job details
	c.JSON(http.StatusOK, gin.H{"job": job})
}

func ViewApplicants(c *gin.Context) {
	// Get authenticated user's ID from JWT token
	userID := c.MustGet(string(middleware.ContextKeyUserID)).(uint)

	// Check if the authenticated user is an admin
	isAdmin := IsAdmin(userID)

	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admin users can view applicants"})
		return
	}

	// Fetch users with userType as Admin from the database
	var applicants []models.User
	err := database.DB.Where("user_type = ?", "Applicant").Find(&applicants).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch admins"})
		return
	}

	// Remove password hash before sending the response
	for i := range applicants {
		applicants[i].PasswordHash = ""
	}

	c.JSON(http.StatusOK, applicants)
}

func IsAdmin(userID uint) bool {
	// Query user from the database
	var user models.User
	err := database.DB.First(&user, userID).Error
	if err != nil {
		return false
	}

	// Check if user type is admin
	if user.UserType == "Admin" {
		return true
	}

	return false
}

// CustomUserRespons struct to exclude sensitive fields
type CustomUserRespons struct {
	ID         uint                  `json:"id"`
	Name       string                `json:"name"`
	Email      string                `json:"email"`
	Address    string                `json:"address"`
	UserType   string                `json:"user_type"`
	ProfileHea string                `json:"profile_headline"`
	Profile    CustomProfileResponse `json:"profile,omitempty"`
}

// CustomProfileResponse struct to exclude sensitive fields
type CustomProfileResponse struct {
	UserID        uint   `json:"user_id"`
	ResumeFileURL string `json:"resume_file_url"`
	Skills        string `json:"skills"`
	Education     string `json:"education"`
	Experience    string `json:"experience"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
}

// ViewApplicantData retrieves and returns applicant data
func ViewApplicantData(c *gin.Context) {
	// Get authenticated user's ID from JWT token
	userID := c.MustGet(string(middleware.ContextKeyUserID)).(uint)

	// Check if the authenticated user is an admin
	if !IsAdmin(userID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admin users can view applicants"})
		return
	}

	// Extract applicant_id from the path parameters
	applicantIDstr := c.Param("applicant_id")
	applicantID, err := strconv.ParseUint(applicantIDstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid applicant ID"})
		return
	}

	// Fetch the applicant user data from the database
	var user models.User
	if err := database.DB.First(&user, applicantID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Applicant not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applicant"})
		}
		return
	}

	// Fetch the profile data from the database
	var profile models.Profile
	if err := database.DB.Where("user_id = ?", applicantID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profile data"})
		}
		return
	}

	// Map data to custom response struct
	customUser := CustomUserRespons{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Address:    user.Address,
		UserType:   user.UserType,
		ProfileHea: user.ProfileHead,
		Profile: CustomProfileResponse{
			UserID:        profile.UserID,
			ResumeFileURL: profile.ResumeFileURL, // Adjust field name if different
			Skills:        profile.Skills,
			Education:     profile.Education,
			Experience:    profile.Experience,
			Name:          profile.Name,
			Email:         profile.Email, // Corrected field
			Phone:         profile.Phone,
		},
	}

	// Send the response
	c.JSON(http.StatusOK, customUser)
}
