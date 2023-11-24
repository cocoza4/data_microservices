package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authorize(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	user, password, hasAuth := c.Request.BasicAuth()
	if len(auth) == 0 && !hasAuth {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authentication required"})
		return
	}

	if hasAuth {
		if user != os.Getenv("SECRET_USER") || password != os.Getenv("SECRET_PASSWORD") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		}
	} else {
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "failed to parse token"})
		}
	}
}
