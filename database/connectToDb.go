package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/models"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	DB = connectDB()
	return DB
}

func connectDB() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	dns := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nConnected to DATABASE: ", db.Name())

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Application{},
		&models.Job{},
		&models.Profile{},
	); err != nil {
		log.Fatal(err)
	}

	return db
}

func FetchUserDetails(username string) (*models.User, error) {
	var user models.User
	result := DB.Where("email = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func ValidateUserCredentials(email, password string) (*models.User, error) {
	// Fetch user details from the database based on the provided email
	user, err := FetchUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Check if the user was found
	if user == nil {
		// Handle the case when the user does not exist
		return nil, errors.New("user not found")
	}

	// Print the hashed password fetched from the database
	fmt.Println("Hashed password from database:", user.PasswordHash)

	// Compare the provided password with the hashed password stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		// Print the hashed password provided during the login request
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Failed to hash provided password for comparison:", err)
		} else {
			fmt.Println("Hashed password provided during login request:", string(hashedPassword))
		}
		// handle the case when the password doesn't match
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func FetchUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
