package repository

import (
	"gorm.io/gorm"
	models "spotlas/model"
)

type IOSportRepository interface {
	GetSpotsByOffsetLimit(offset, limit int) ([]models.Spot, error)
	GetNumberOfSpots() (int64, error)
}

type SpotsRepository struct {
	db *gorm.DB
}

func NewSpotsRepository(db *gorm.DB) IOSportRepository {
	return &SpotsRepository{db: db}
}

func (repo *SpotsRepository) GetSpotsByOffsetLimit(offset, limit int) ([]models.Spot, error) {
	var spots []models.Spot

	rows, err := repo.db.Raw("SELECT id, name, website, ST_AsBinary(coordinates) as coordinates, description, rating "+
		"FROM \"MY_TABLE\" OFFSET ? LIMIT ?", offset, limit).Rows()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var spot models.Spot
		rows.Scan(&spot.ID, &spot.Name, &spot.Website, &spot.Coordinates, &spot.Description, &spot.Rating)
		spots = append(spots, spot)
	}
	defer rows.Close()
	return spots, nil
}

func (repo *SpotsRepository) GetNumberOfSpots() (int64, error) {
	var count int64
	if err := repo.db.Table("MY_TABLE").Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}
