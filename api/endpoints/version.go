package endpoints

import "github.com/gin-gonic/gin"

const version string = "v0.1.0"

type versionResponse struct {
	Message string `json:"message"`
	Version string `json:"version"`
}

func Version(c *gin.Context) {
	c.JSON(200, versionResponse{
		Message: "Welcome to the MLModelScope API!",
		Version: version,
	})
}