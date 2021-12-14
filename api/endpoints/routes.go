package endpoints

import "github.com/gin-gonic/gin"

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	r.GET("/", Version)
	r.POST("/predict", Predict)

	return r
}

