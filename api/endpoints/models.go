package endpoints

import (
	"api/api_db"
	"api/db/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ModelListResponse struct {
	Models []models.Model `json:"models,omitempty"`
}

func ListModels(c *gin.Context) {
	db, _ := api_db.GetDatabase()
	m, _ := db.GetAllModels()

	c.JSON(200, ModelListResponse{Models: m})
}

func ListModelsByFrameworkId(c *gin.Context) {
	db, _ := api_db.GetDatabase()
	idParam := c.Param("id")
	frameworkId, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(400, &ErrorResponse{Error: fmt.Sprintf("invalid Framework ID: %s", idParam)})
		return
	}

	m, _ := db.GetModelsForFramework(frameworkId)

	c.JSON(200, ModelListResponse{Models: m})
}

func ListModelsByTask(c *gin.Context) {
	db, _ := api_db.GetDatabase()
	task := c.Param("task")
	m, _ := db.GetModelsByTask(task)

	c.JSON(200, ModelListResponse{Models: m})
}
