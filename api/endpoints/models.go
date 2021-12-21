package endpoints

import (
	"api/api_db"
	"api/db/models"
	"github.com/gin-gonic/gin"
)

type ModelListResponse struct {
	Models []models.Model `json:"models,omitempty"`
}

func ListModels(c *gin.Context) {
	db, _ := api_db.GetDatabase()
	m, _ := db.GetAllModels()

	c.JSON(200, ModelListResponse{Models: m})
}
