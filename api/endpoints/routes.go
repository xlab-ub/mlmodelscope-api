package endpoints

import "github.com/gin-gonic/gin"

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	r.GET("/", Version)
	r.GET("/frameworks", ListFrameworks)
	models := r.Group("/models")
	{
		models.GET("", ListModels)
		models.GET("/framework/:id", ListModelsByFrameworkId)
		models.GET("/task/:task", ListModelsByTask)
	}
	r.POST("/predict", Predict)

	return r
}
