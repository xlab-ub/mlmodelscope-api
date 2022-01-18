package endpoints

import (
	"api/api_db"
	"api/db/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ModelListResponse struct {
	Models []models.Model `json:"models,omitempty"`
}

func ListModels(c *gin.Context) {
	frameworkId, err := getFrameworkId(c)
	if err != nil {
		c.JSON(400, &ErrorResponse{Error: "invalid Framework ID"})
		return
	}
	task := c.Query("task")
	architecture := c.Query("architecture")

	db, _ := api_db.GetDatabase()
	m, _ := db.QueryModels(frameworkId, task, architecture)

	c.JSON(200, ModelListResponse{Models: m})
}

func getFrameworkId(c *gin.Context) (frameworkId int, err error) {
	framework := c.Query("framework")

	if framework == "" {
		return
	}

	frameworkId, err = strconv.Atoi(framework)
	return frameworkId, err
}
