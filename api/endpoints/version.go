package endpoints

import "github.com/gin-gonic/gin"

const version string = "v0.1.0"

func Version(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to the MLModelScope API!",
		"version": version,
	})
}