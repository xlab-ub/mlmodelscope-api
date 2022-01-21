package endpoints

import "github.com/gin-gonic/gin"

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	r.GET("/", Version)

	experiments := r.Group("/experiments")
	{
		experiments.GET("/:experiment_id", GetExperiment)
	}

	r.GET("/frameworks", ListFrameworks)

	models := r.Group("/models")
	{
		models.GET("", ListModels)
	}
	r.POST("/predict", Predict)

	trial := r.Group("/trial")
	{
		trial.GET("/:trial_id", GetTrial)
	}

	return r
}
