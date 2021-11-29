package main

import (
	"api/endpoints"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", endpoints.Version)
	r.Run()
}
