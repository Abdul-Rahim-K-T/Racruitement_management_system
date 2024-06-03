package controllers

import (
	"fmt"
	"net/http"

	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/database"
	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/helpers"
	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/middleware"
	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/models"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

var DB *gorm.DB

var validate *validator.Validate // Validator instance

func init() {
	validate = validator.New()
}

func Signup(c *gin.Context) {
	// Parse request body and create user
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate struct
	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	var existingUser models.User
	if err := database.DB.Where("name = ?", user.Name).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	// Check if email already exists
	var count int64
	if err := database.DB.Model(&models.User{}).Where("email = ?", user.Email).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := HashPassword(user.PasswordHash)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.PasswordHash = hashedPassword

	// Create user in the database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Return success message
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	// Parse request body to get user credentials
	var loginData LoginData
	if err := c.ShouldBindJSON(&loginData); err != nil {
		fmt.Println("Error parsing JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve user from the database based on the provided email
	user, err := database.FetchUserByEmail(loginData.Email)
	if err != nil {
		fmt.Println("Error fetching user:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare password
	status, err := compareHashPassword(user.PasswordHash, loginData.Password)
	if err != nil || !status {
		fmt.Println("Unauthorized access:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token for authenticated user
	tokenMap, err := helpers.GenerateToken(user.ID, user.Name, user.UserType)
	if err != nil {
		fmt.Println("Error generating token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Extract tokens from the map
	accessToken := tokenMap["access_token"]
	refreshToken := tokenMap["refresh_token"]

	// Set JWT in browser cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("UserAuth", accessToken, 3600*24*30, "", "", false, true)
	c.SetCookie("RefreshToken", refreshToken, 3600*24*30, "", "", false, true)

	// Return the generated tokens
	c.JSON(http.StatusOK, gin.H{
		"token": accessToken,
		"user": gin.H{
			"name":     user.Name,
			"email":    user.Email,
			"userType": user.UserType,
		},
	})
}

func Logout(c *gin.Context) {
	// Retrieve the user ID from the context
	userID, exists := c.Get(string(middleware.ContextKeyUserID))
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user data in context"})
		return
	}

	// Type assertion for the user ID
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID is of incorrect type"})
		return
	}

	// Check if the database connection is initialized
	if database.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	// Fetch the user data using the user ID
	var fetchedUser models.User
	if err := database.DB.First(&fetchedUser, userIDUint).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	// Print the user name
	fmt.Printf("Logging out user: %s\n", fetchedUser.Name)

	// Perform actions such as logging out the user (e.g., deleting the cookies)
	c.SetCookie("UserAuth", "", -1, "/", "localhost", false, true)
	c.SetCookie("RefreshToken", "", -1, "/", "localhost", false, true)

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{
		"message":  "Successfully logged out",
		"user":     fetchedUser.Name,
		"email":    fetchedUser.Email,
		"userType": fetchedUser.UserType,
	})
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func compareHashPassword(hashedPassword string, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false, err
	}
	return true, nil
}

// func Login(c *gin.Context) {
// 	// Parse request body to get user credentials
// 	var loginData LoginData
// 	if err := c.ShouldBindJSON(&loginData); err != nil {
// 		fmt.Println("Error parsing JSON:", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Print the parsed login data
// 	fmt.Println("LoginData:", loginData)

// 	// Retrieve user from the database based on the provided email
// 	user, err := database.FetchUserByEmail(loginData.Email)
// 	if err != nil {
// 		fmt.Println("Error fetching user:", err)
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
// 		return
// 	}

// 	// Hash the password before storing it in the database
// 	hashedPassword, err := HashPassword(loginData.Password)
// 	if err != nil {
// 		fmt.Println("Error hashing password:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
// 		return
// 	}
// 	user.PasswordHash = hashedPassword

// 	// Print hashed password retrieved from the database
// 	fmt.Println("Hashed password from the database:", user.PasswordHash)

// 	// Compare the provided password with the hashed password stored in the database during login
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginData.Password)); err != nil {
// 		fmt.Println("Error validating password:", err)
// 		// Print the hashed password retrieved from the database
// 		fmt.Println("Hashed password from the database:", user.PasswordHash)
// 		// Print the hash value of the login password
// 		hashedLoginPassword, _ := bcrypt.GenerateFromPassword([]byte(loginData.Password), bcrypt.DefaultCost)
// 		fmt.Println("Hash value of login password:", hashedLoginPassword)
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
// 		return
// 	}

// 	// Print the user details
// 	hashedLoginPassword, _ := bcrypt.GenerateFromPassword([]byte(loginData.Password), bcrypt.DefaultCost)
// 	fmt.Println("hash value of hashed pass logi:", hashedLoginPassword)
// 	fmt.Println("Authenticated user:", user, hashedLoginPassword)

// 	// Generate JWT token for authenticated user
// 	token, err := helpers.GenerateToken(user.Email)
// 	if err != nil {
// 		fmt.Println("Error generating token:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
// 		return
// 	}

// 	// Print the generated token
// 	fmt.Println("Generated token:", token)

// 	// Response of token
// 	c.JSON(http.StatusOK, gin.H{"token": token})
// }

// func compareHashPassword(hashedPassword, password string) (error, bool) {

// 	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
// 		return err, false
// 	}

// 	return nil, true
// }

// func Login(c *gin.Context) {
// 	type userDetail struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 	}
// 	var userCredentials userDetail
// 	if err := c.Bind(&userCredentials); err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	//finding with username in database
// 	var user models.User
// 	database.DB.Where("user_name=?", userCredentials.Username).First(&user)
// 	//checking user is blocked or not
// 	if user.IsBlocked {
// 		c.JSON(401, gin.H{
// 			"error": "Unautharized access user is blocked",
// 		})
// 		return
// 	}
// 	//comparing password with database
// 	status := compareHashPassword(user.Password, userCredentials.Password)
// 	//checking password and username
// 	if !status || userCredentials.Username != user.User_Name {
// 		c.JSON(401, gin.H{
// 			"error": "Unautharized access Please check username or password",
// 		})
// 		return
// 	}
// 	//generating token with jwt
// 	token, err := helper.GenerateJWTToken(user.User_Name, "user", user.Email, int(user.User_ID))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{})
// 		return
// 	}
// 	//set jwt in browser
// 	c.SetSameSite(http.SameSiteLaxMode)
// 	c.SetCookie("jwt_token", token, 3600*24, "", "", true, true)
// 	//success message
// 	c.JSON(200, gin.H{
// 		"message": user.User_Name + " successfully logged",
// 	})
// }
