package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		log.Printf("Authorization Header: %s\n", tokenString)

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		log.Printf("Trimmed Token: %s\n", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNoLocation
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Printf("Token invalid: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		authorID, ok := claims["author_id"].(float64) // Assuming author_id is a float64 in the token
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Missing author_id"})
			c.Abort()
			return
		}

		c.Set("authorID", uint(authorID))
		log.Printf("authorID set in context: %d\n", uint(authorID))

		c.Next()
	}
}
