package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"os"
	"path/filepath"

	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/database"
	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/middleware"
	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/models"
	"github.com/gin-gonic/gin"
)

func ViewJobs(c *gin.Context) {
	// Fetch jobs from the database
	db := database.InitDB()

	var jobs []models.Job
	if err := db.Preload("PostedBy").Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database connection"})
		return
	}
	defer sqlDB.Close()

	// Map jobResponse []JobResponse
	jobResponses := []gin.H{}
	for _, job := range jobs {
		jobResponse := gin.H{
			"id":                 job.ID,
			"title":              job.Title,
			"description":        job.Description,
			"posted_on":          job.PostedOn,
			"total_applications": job.TotalApplications,
			"company_name":       job.CompanyName,
			"posted_by": gin.H{
				"id":    job.PostedBy.ID,
				"name":  job.PostedBy.Name,
				"email": job.PostedBy.Email,
			},
		}
		jobResponses = append(jobResponses, jobResponse)
	}

	// Return the list of jobs
	c.JSON(http.StatusOK, gin.H{"jobs": jobResponses})
}

func ApplyJob(c *gin.Context) {
	// Get job ID from request parameters
	jobID := c.Query("job_id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID is required"})
		return
	}

	// Fetch job details from the database
	var job models.Job
	if err := database.DB.Where("id = ?", jobID).First(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job details"})
		return
	}

	// Check if the job exists
	if job.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	// Fetch user ID from the context
	userID, exists := c.Get(string(middleware.ContextKeyUserID))
	fmt.Println("userID:", userID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	// Fetch user type from the context
	userType, exists := c.Get(string(middleware.ContextKeyUserType))
	fmt.Println("userType:", userType)
	fmt.Println("exists:", exists)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User type not found"})
		return
	}

	// Ensure userID is of the correct type (uint)
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Check if the user is an applicant
	if userType != "Applicant" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "only applicants can apply for jobs"})
		return
	}

	// Create a new application
	application := models.Application{
		JobID:  job.ID,
		UserID: userIDUint,
		Status: "Pending", // You can set an initial status here
	}

	// Save the application to the database
	if err := database.DB.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application submitted successfully"})
}

func UploadResume(c *gin.Context) {
	// Ensure the user is an applicant
	userType, exists := c.Get(string(middleware.ContextKeyUserType))
	if !exists || userType != "Applicant" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only applicants can upload resumes"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get(string(middleware.ContextKeyUserID))
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	// Ensure the user ID is of the correct type (uint)
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Parse the multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	// Retrieve file
	file, handler, err := c.Request.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
		return
	}
	defer file.Close()

	// Validate file type
	ext := filepath.Ext(handler.Filename)
	if ext != ".pdf" && ext != ".docx" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format. Only PDF or DOCX allowed"})
		return
	}

	// Define the file path
	dirPath := "uploads/resumes/"
	filePath := fmt.Sprintf("%s%d%s", dirPath, userIDUint, ext)

	// Create the directory if it doesn't exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
			return
		}
	}

	// Save the file
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Process the resume using APILayer API
	resumeData, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read resume file"})
		return
	}

	parsedResume, err := parseResumeWithAPILayer(resumeData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse resume: %v", err)})
		return
	}

	// Log parsed resume data for debugging
	fmt.Printf("Parsed Resume: %+v\n", parsedResume)

	// Check if a profile exists for the user
	var profileCount int64
	if err := database.DB.Model(&models.Profile{}).Where("user_id = ?", userIDUint).Count(&profileCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check profile existence"})
		return
	}

	// Create or update profile based on existence
	if profileCount == 0 {
		// No profile found, create a new one
		profile := models.Profile{
			UserID:        userIDUint,
			ResumeFileURL: filePath,
			Skills:        parsedResume.Skills,
			Education:     parsedResume.Education,
			Experience:    parsedResume.Experience,
			Name:          parsedResume.Name,
			Email:         parsedResume.Email,
			Phone:         parsedResume.Phone,
		}
		if err := database.DB.Create(&profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
			return
		}
		fmt.Printf("Profile created: %+v\n", profile)
	} else {
		// Profile found, update it
		var profile models.Profile
		if err := database.DB.Where("user_id = ?", userIDUint).First(&profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profile"})
			return
		}
		profile.ResumeFileURL = filePath
		profile.Skills = parsedResume.Skills
		profile.Education = parsedResume.Education
		profile.Experience = parsedResume.Experience
		profile.Name = parsedResume.Name
		profile.Email = parsedResume.Email
		profile.Phone = parsedResume.Phone
		if err := database.DB.Save(&profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
			return
		}
		fmt.Printf("Profile updated: %+v\n", profile)
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "Resume uploaded successfully"})
}

// parseResumeWithAPILayer function calls the APILayer API to parse the resume and extract details
func parseResumeWithAPILayer(resumeData []byte) (*models.Profile, error) {
	url := "https://api.apilayer.com/resume_parser/upload"

	// Prepare the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(resumeData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("apiKey", "0bWeisRWoLj3UdXt3MXMSMWptYFIpQfS")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Log the status code and headers for debugging
	fmt.Printf("API response status: %d\n", resp.StatusCode)
	fmt.Printf("API response headers: %v\n", resp.Header)

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to parse resume, status code: %d, response body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read the raw response body for logging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Log the raw response body for debugging
	fmt.Printf("API raw response body: %s\n", string(bodyBytes))

	// Define a struct to match the structure of the JSON response
	type APIResponse struct {
		Name      string   `json:"name"`
		Address   string   `json:"address"`
		Email     string   `json:"email"`
		Phone     string   `json:"phone"`
		Skills    []string `json:"skills"`
		Education []struct {
			Name  string   `json:"name"`
			Dates []string `json:"dates"`
		} `json:"education"`
		Experience []struct {
			Title        string   `json:"title"`
			Dates        []string `json:"dates,omitempty"`
			DateStart    string   `json:"date_start,omitempty"`
			DateEnd      string   `json:"date_end,omitempty"`
			Organization string   `json:"organization"`
		} `json:"experience"`
	}

	// Parse the response body
	var result APIResponse

	// Decode the JSON response
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("error decoding response body: %v", err)
	}

	// Log the parsed data for debugging
	fmt.Printf("API Layer parsed Result: %+v\n", result)

	// Check if parsed data is empty and log a warning
	if len(result.Skills) == 0 && len(result.Education) == 0 && len(result.Experience) == 0 && result.Name == "" && result.Email == "" && result.Phone == "" {
		fmt.Println("Warning: Parsed result is empty, please check  the API response and the request data.")
	}

	// Convert skills and experience slices to a single string
	skills := strings.Join(result.Skills, ", ")
	experience := ""
	for _, exp := range result.Experience {
		experience += fmt.Sprintf("%s at %s (%s to %s); ", exp.Title, exp.Organization, exp.DateStart, exp.DateEnd)
	}
	// Convert education struct to a readable string format
	var educationDetails []string
	for _, edu := range result.Education {
		educationDetails = append(educationDetails, fmt.Sprintf("%s (%s)", edu.Name, strings.Join(edu.Dates, ", ")))
	}
	education := strings.Join(educationDetails, "; ")

	// Return the parsed resume details
	return &models.Profile{
		Skills:     skills,
		Education:  education,
		Experience: experience,
		Name:       result.Name,
		Email:      result.Email,
		Phone:      result.Phone,
	}, nil
}
