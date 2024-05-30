package main

import (
	"github.com/akashdesale98/exoPlanetService/service"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	service.RegisterAPIs(router)

	// Start the Gin server on port 8080
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
