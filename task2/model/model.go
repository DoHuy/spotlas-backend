package models

import (
	"database/sql"
	geo "github.com/paulmach/go.geo"
)

type Spot struct {
	ID          string          `json:"id" gorm:"primaryKey"`
	Name        sql.NullString  `json:"name"`
	Website     sql.NullString  `json:"website"`
	Coordinates *geo.Point      `json:"coordinates"`
	Description sql.NullString  `json:"description"`
	Rating      sql.NullFloat64 `json:"rating"`
	Distance    float64         `json:"distance"`
}
