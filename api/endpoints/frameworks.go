package endpoints

import (
	"api/api_db"
	"api/db/models"
	"github.com/gin-gonic/gin"
)

type ListFrameworksResponse struct {
	Frameworks []models.Framework `json:"frameworks,omitempty"`
}

func ListFrameworks(c *gin.Context) {
	db, _ := api_db.GetDatabase()
	f, _ := db.GetAllFrameworks()

	c.JSON(200, ListFrameworksResponse{Frameworks: f})
}
