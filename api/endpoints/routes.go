package endpoints

import "github.com/gin-gonic/gin"

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	r.GET("/", Version)
	r.GET("/models", ListModels)
	r.POST("/predict", Predict)

	return r
}
