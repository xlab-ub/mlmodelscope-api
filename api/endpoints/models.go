package endpoints

import (
	"api/api_db"
	"api/db/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ModelListResponse struct {
	Models []models.Model `json:"models"`
}

func ListModels(c *gin.Context) {
	frameworkId, err := getFrameworkId(c)
	if err != nil {
		c.JSON(400, &ErrorResponse{Error: "invalid Framework ID"})
		return
	}
	task := c.Query("task")
	architecture := c.Query("architecture")
	query := c.Query("q")

	db, _ := api_db.GetDatabase()
	m, _ := db.QueryModels(frameworkId, task, architecture, query)

	c.JSON(200, ModelListResponse{Models: m})
}

func getFrameworkId(c *gin.Context) (frameworkId int, err error) {
	framework := c.Query("framework")

	if framework == "" {
		return
	}

	frameworkId, err = strconv.Atoi(framework)
	return
}

func GetModelById(c *gin.Context) {
	modelId, err := getModelId(c)
	if err != nil {
		c.JSON(400, &ErrorResponse{Error: "invalid Model ID"})
		return
	}

	db, _ := api_db.GetDatabase()
	m, _ := db.GetModelById(uint(modelId))

	c.JSON(200, ModelListResponse{Models: []models.Model{ *m }})
}

func getModelId(c *gin.Context) (modelId int, err error) {
	model := c.Param("model_id")

	if model == "" {
		return
	}

	modelId, err = strconv.Atoi(model)
	return
}
