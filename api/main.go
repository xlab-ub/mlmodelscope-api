package main

import (
	"api/endpoints"
)

func main() {
	r := endpoints.SetupRoutes()
	r.Run()
}
