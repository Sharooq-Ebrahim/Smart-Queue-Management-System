package middleware

import (
	"log"
	"net/http"
	"smart-queue/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddlware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		}

		cfg, _ := config.LoadConfig()

		tokenStr := strings.Split(authHeader, "Bearer ")[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		log.Println("claims", claims)
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Next()

	}

}
