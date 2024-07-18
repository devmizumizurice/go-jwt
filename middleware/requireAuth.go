package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/devmizumizurice/go-jwt/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAccessToken(c *gin.Context) {
	fmt.Println("Im middleware")

	token, err := c.Cookie("access_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "MISSING_TOKEN"})
		return
	}

	parsedToken, err := utils.VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if ok && parsedToken.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "EXPIRED_TOKEN"})
			return
		}
	}

	c.Set("sub", claims["sub"])

	c.Next()

}
