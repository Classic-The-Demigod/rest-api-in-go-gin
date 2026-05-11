package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (app *application) AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the request header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is required"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == authHeader {
			c.AbortWithStatusJSON(401, gin.H{"error": "Bearer token is required"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(app.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims"})
			return
		}

		userId := claims["user_id"].(float64)

		user, err := app.models.Users.GetUserByID(int(userId))
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized access"})
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
