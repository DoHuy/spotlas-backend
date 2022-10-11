package types

import (
	geo "github.com/paulmach/go.geo"
)

type Figure int

const (
	Circle Figure = iota + 1
	Square
)

type QueryDto struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
	Type      Figure  `json:"type"`
}

type SpotDto struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name"`
	Website     string     `json:"website"`
	Coordinates *geo.Point `json:"coordinates"`
	Description string     `json:"description"`
	Rating      float64    `json:"rating"`
	Distance    float64    `json:"distance"`
}

type Coordinate struct {
	Lat  float64
	Long float64
}
