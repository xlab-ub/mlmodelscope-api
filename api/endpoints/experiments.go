package endpoints

import (
	"api/api_db"
	"github.com/gin-gonic/gin"
)

func GetExperiment(c *gin.Context) {
	experimentId := c.Param("experiment_id")

	if db, err := api_db.GetDatabase(); err != nil {
		c.JSON(500, NewErrorResponse(err))
	} else {
		if experiment, err := db.GetExperimentById(experimentId); err != nil {
			c.JSON(404, NotFoundResponse)
		} else {
			c.JSON(200, experiment)
		}
	}
}
