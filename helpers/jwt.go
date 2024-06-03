package helpers

import (
	"os"

	"time"

	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	// "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))
var refreshSecretKey = []byte("secret")

type Claims struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	UserType string `json:"userType"`
	jwt.StandardClaims
}

func GenerateToken(id uint, name, userType string) (map[string]string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Id:   id,
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		UserType: userType, // Add usersType in the claims
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["id"] = id
	rtClaims["name"] = name
	rtClaims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	rt, err := refreshToken.SignedString(refreshSecretKey)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  tokenString,
		"refresh_token": rt,
	}, nil
}

func VerifyToken(tokenString string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	user := &models.User{
		Model: gorm.Model{
			ID: claims.Id, //
		},
		Name:     claims.Name,
		UserType: claims.UserType,
	}

	return user, nil
}
