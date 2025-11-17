package utils

import "github.com/gin-gonic/gin"

func ResponseSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(200, gin.H{
		"status":  true,
		"message": message,
		"data":    data,
	})
}

func ResponseError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status":  false,
		"message": message,
	})
}
