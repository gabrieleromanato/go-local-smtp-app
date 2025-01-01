package app

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateJWT(username string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}
		tokenString := tokenHeader[7:]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			secret := os.Getenv("JWT_SECRET")
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("claims", claims)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()

	}
}

func HandleLogin(c *gin.Context) {
	loginValid := false
	var body User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email := body.Email
	password := body.Password
	emailStore, err := NewEmailStore(GetDSN())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating the database"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating the database"})
	}
	session := Session{store: emailStore}
	loginValid = session.UserExists(email, password)
	if !loginValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
	token, err := CreateJWT(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating the token"})
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func CheckTokenExpiration(c *gin.Context) {
	secret := os.Getenv("JWT_SECRET")
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		c.JSON(http.StatusOK, gin.H{"error": "Authorization header missing or invalid"})
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm used
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return secret, nil
	})

	if err != nil {
		// Check for specific JWT errors
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				c.JSON(http.StatusOK, gin.H{"error": "Token has expired"})
			} else {
				c.JSON(http.StatusOK, gin.H{"error": "Invalid token"})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"error": "Failed to parse token"})
		}
	}
	if !token.Valid {
		c.JSON(http.StatusOK, gin.H{"error": "Invalid token"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Token is valid"})
}
