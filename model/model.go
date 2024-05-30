package model

type ExoPlanet struct {
	ID            string  `json:"id"`
	Name          string  `json:"name" validate:"required"`
	Description   string  `json:"description" validate:"required"`
	Distance      float64 `json:"distance" validate:"required"`
	Radius        float64 `json:"radius" validate:"required"`
	Mass          float64 `json:"mass" validate:"required_if=ExoPlanetType Terrestrial"`
	ExoPlanetType string  `json:"exo_planet_type" validate:"oneof=Terrestrial GasGiant"`
}

type Fuel struct {
	ID       string `json:"id" validate:"required"`
	Capacity int    `json:"capacity" validate:"required"`
}
