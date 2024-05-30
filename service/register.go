package service

import (
	"github.com/gin-gonic/gin"
)

func RegisterAPIs(router *gin.Engine) {
	router.GET("/exoplanet", GetExoplanet())
	router.POST("/exoplanet", AddExoplanet())
	router.PUT("/exoplanet", UpdateExoplanet())
	router.DELETE("/exoplanet", DeleteExoplanet())
	router.GET("/listexoplanets", ListExoplanets())
	router.POST("/estimate", EstimateFuel())
}
