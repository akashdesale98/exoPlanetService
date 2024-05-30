package service

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/akashdesale98/exoPlanetService/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var inMemoryStore = make(map[string]model.ExoPlanet)
var validate = validator.New(validator.WithRequiredStructEnabled())

func AddExoplanet() func(c *gin.Context) {
	return func(c *gin.Context) {

		var planet model.ExoPlanet

		// Bind JSON request body to the ExoPlanet struct
		if err := c.ShouldBindJSON(&planet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := validate.Struct(planet)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				log.Println("error while validating structs", err)
				return
			}

			field := ""
			for _, err := range err.(validator.ValidationErrors) {
				field = err.Field()
			}

			c.JSON(http.StatusBadRequest, gin.H{"error": "Please, send valid parameters , " + strings.ToLower(field)})
			return
		}

		ok := IsExoPlanetExists(planet.Name)
		if ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Oops! Planet already exists. Please, try different name."})
			return
		}

		planet.ID = uuid.NewString()
		inMemoryStore[planet.ID] = planet

		c.JSON(http.StatusOK, planet)
		log.Print("successfully added new exoplanet")
	}
}

func IsExoPlanetExists(name string) bool {

	for _, v := range inMemoryStore {
		if v.Name == name {
			return true
		}
	}

	return false
}

func ListExoplanets() func(c *gin.Context) {
	return func(c *gin.Context) {

		// Extract values from map
		var planets []model.ExoPlanet
		for _, planet := range inMemoryStore {
			planets = append(planets, planet)
		}

		// Sort by radius
		sort.Slice(planets, func(i, j int) bool {
			return planets[i].Radius < planets[j].Radius
		})

		c.JSON(http.StatusOK, planets)
		log.Print("successfully listed all exoplanets")
	}
}

func GetExoplanet() func(c *gin.Context) {
	return func(c *gin.Context) {

		exoPlanetID := c.Query("id")
		val, ok := inMemoryStore[exoPlanetID]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "exoplanet not found"})
			return
		}

		c.JSON(http.StatusOK, val)
		log.Print("successfully fetched required exoplanet")
	}
}

func UpdateExoplanet() func(c *gin.Context) {
	return func(c *gin.Context) {

		var planet model.ExoPlanet

		// Bind JSON request body to the ExoPlanet struct
		if err := c.ShouldBindJSON(&planet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := validate.Struct(planet)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				log.Println("error while validating structs", err)
				return
			}

			field, tag := "", ""
			for _, err := range err.(validator.ValidationErrors) {
				field, tag = err.Field(), err.Tag()
			}

			log.Println("error field", field)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please, send valid parameters." + tag})
			return
		}

		_, ok := inMemoryStore[planet.ID]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "exoplanet not found"})
			return
		}

		inMemoryStore[planet.ID] = planet
		c.JSON(http.StatusOK, planet)
		log.Print("successfully updated exoplanet")
	}
}

func DeleteExoplanet() func(c *gin.Context) {
	return func(c *gin.Context) {

		var planet model.ExoPlanet

		// Bind JSON request body to the ExoPlanet struct
		if err := c.ShouldBindJSON(&planet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := validate.Struct(planet)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				log.Println("error while validating structs", err)
				return
			}

			field, tag := "", ""
			for _, err := range err.(validator.ValidationErrors) {
				field, tag = err.Field(), err.Tag()
			}

			log.Println("error field", field)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please, send valid parameters." + tag})
			return
		}

		_, ok := inMemoryStore[planet.ID]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "exoplanet not found"})
			return
		}

		delete(inMemoryStore, planet.ID)
		c.JSON(http.StatusOK, "successfully deleted the exoplanet")
		log.Print("successfully deleted the exoplanet")
	}
}

func EstimateFuel() func(c *gin.Context) {
	return func(c *gin.Context) {

		var fuel model.Fuel

		// Bind JSON request body to the fuel struct
		if err := c.ShouldBindJSON(&fuel); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := validate.Struct(fuel)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				log.Println("error while validating structs", err)
				return
			}

			field := ""
			for _, err := range err.(validator.ValidationErrors) {
				field = err.Field()
			}

			c.JSON(http.StatusBadRequest, gin.H{"error": "Please, send valid parameters , " + strings.ToLower(field)})
			return
		}

		val, ok := inMemoryStore[fuel.ID]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "exoplanet not found"})
			return
		}

		m := 0.5
		if val.ExoPlanetType == "Terrestrial" {
			m = val.Mass
		}

		r := val.Radius
		g := m / (r * r)
		f := (val.Distance / (g * g)) * float64(fuel.Capacity)

		c.JSON(http.StatusOK, fmt.Sprintf("estimated fuel to reach %s exoplanet from earth is %f units", val.Name, f))
	}
}
