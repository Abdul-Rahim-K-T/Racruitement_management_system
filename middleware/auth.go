package middleware

import (
	"fmt"
	"net/http"

	"github.com/Abdul-Rahim-K-T/Recruitment_management_system/helpers"
	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	ContextKeyUserID   contextKey = "user_id"
	ContextKeyUserType contextKey = "user_type"
)

func AuthMiddleware(c *gin.Context) {
	// Get the token from the request cookie
	tokenString, err := c.Cookie("UserAuth")
	if err != nil {
		fmt.Println("Error retrieving cookie:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not found"})
		c.Abort()
		return
	}

	// Verify the token
	user, err := helpers.VerifyToken(tokenString)
	if err != nil {
		fmt.Println("Error verifying token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return
	}

	// Debug print the user data
	fmt.Printf("User data in AuthMiddleware: %+v\n", user)

	// Set user ID in the context
	c.Set(string(ContextKeyUserID), user.ID)
	c.Set(string(ContextKeyUserType), user.UserType)
	fmt.Println("userid", user.ID)
	fmt.Println("usertype", user.UserType)
	fmt.Println("AuthMiddleware is completely executed")

	// Continue to the next handler
	c.Next()
}

// type contextKey string

// const (
// 	ContextKeyUser contextKey = "user"
// )

// func AuthMiddleware(c *gin.Context) {
// 	fmt.Println("AuthMiddleware executing")

// 	tokenString, err := c.Cookie("jwt_token")
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
// 		c.Abort()
// 		return
// 	}
// 	fmt.Println("Extracted JWT token:", tokenString)

// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(os.Getenv("SECRET_KEY")), nil
// 	})
// 	if err != nil {
// 		fmt.Println("Error of Invalid token:", err)
// 		fmt.Println("Token while invalid:", token)
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 		c.Abort()
// 		return
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		fmt.Println("Token Claims:", claims)
// 		if float64(time.Now().Unix()) > claims["exp"].(float64) {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
// 			c.Abort()
// 			return
// 		}

// 		var user models.User
// 		if err := database.DB.First(&user, claims["sub"]).Error; err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
// 			c.Abort()
// 			return
// 		}

// 		c.Set(string(ContextKeyUser), user)
// 		c.Next()
// 	} else {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
// 		c.Abort()
// 	}
// }

// func AuthMiddleware(c *gin.Context) {
// 	fmt.Println("AuthMiddleware executing")

// 	// Extract token from cookie
// 	tokenString, err := c.Cookie("jwt_token")
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
// 		c.Abort()
// 		return
// 	}
// 	fmt.Println("Extracted JWT token:", tokenString)

// 	// Verify Secret Key
// 	secretKey := os.Getenv("SECRET_KEY")
// 	fmt.Println("Secret Key:", secretKey)

// 	// Parse and validate token
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		// Ensure the token signing method is as expected
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		// Return the secret key used for signing
// 		return []byte(secretKey), nil
// 	})
// 	if err != nil || !token.Valid {
// 		fmt.Println("Error validating token:", err, "\ntoken:", token)
// 		// Print token claims
// 		fmt.Println("Token Claims:", token.Claims)

// 		// Print token signature algorithm
// 		fmt.Println("Token Signature Algorithm:", token.Method.Alg())
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 		c.Abort()
// 		return
// 	}

// 	// Check token claims
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		if exp, ok := claims["exp"].(float64); ok && float64(time.Now().Unix()) > exp {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
// 			c.Abort()
// 			return
// 		}
// 		email, ok := claims["email"].(string)
// 		if !ok {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email in token claims"})
// 			c.Abort()
// 			return
// 		}
// 		// email, ok := claims["email"].(string)
// 		// if !ok {
// 		//     c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email in token claims"})
// 		//     c.Abort()
// 		//     return
// 		// }

// 		// Create the user object
// 		user := models.User{
// 			// ID:    uint(userID),
// 			Email: email,
// 			// Add other user details as needed
// 		}

// 		// Set the user object in the context using the custom context key
// 		c.Set(string(ContextKeyUser), user)

// 		c.Next()
// 	} else {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 		c.Abort()
// 	}
// }
