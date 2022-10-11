package http

import (
	"context"
	"fmt"
	"math"
	"sort"
	models "spotlas/model"
	"spotlas/repository"
	"spotlas/types"
	"spotlas/util"
	"time"
)

const (
	Batch       = 1000
	LessThan50m = 50
)

type Service struct {
	SpotsRepository repository.IOSportRepository
}

func NewService(repository repository.IOSportRepository) *Service {
	return &Service{
		SpotsRepository: repository,
	}
}

func toSpotDto(spots []models.Spot) []types.SpotDto {
	var restOfSpots []types.SpotDto
	var lessThan50mSpots []types.SpotDto
	for _, spot := range spots {
		var dto = types.SpotDto{
			ID:          spot.ID,
			Name:        spot.Name.String,
			Website:     spot.Website.String,
			Coordinates: spot.Coordinates,
			Description: spot.Description.String,
			Rating:      spot.Rating.Float64,
			Distance:    spot.Distance,
		}
		if spot.Distance < LessThan50m {
			lessThan50mSpots = append(lessThan50mSpots, dto)
		} else {
			restOfSpots = append(restOfSpots, dto)
		}
	}
	sort.SliceStable(lessThan50mSpots, func(i, j int) bool {
		return lessThan50mSpots[i].Distance > lessThan50mSpots[j].Distance
	})
	sort.SliceStable(restOfSpots, func(i, j int) bool {
		return restOfSpots[i].Rating > restOfSpots[j].Rating
	})
	return append(lessThan50mSpots, restOfSpots...)
}

func isBoundedBySquare(centerLat, centerLong, checkedPointLat, checkedPointLong, distance float64) bool {
	var top = util.MaxLatLongOnBearing(centerLat, centerLong, 45, distance)
	var right = util.MaxLatLongOnBearing(centerLat, centerLong, 135, distance)
	var bottom = util.MaxLatLongOnBearing(centerLat, centerLong, 225, distance)
	var left = util.MaxLatLongOnBearing(centerLat, centerLong, 315, distance)
	if top.Lat >= checkedPointLat && checkedPointLat >= bottom.Lat {
		if left.Long <= right.Long && left.Long <= checkedPointLong && checkedPointLong <= right.Long {
			return true
		} else if left.Long > right.Long && (left.Long <= checkedPointLong || checkedPointLong <= right.Long) {
			return true
		}
	}
	return false
}

func (s *Service) getBatchSpotsByCondition(dataChan chan []models.Spot, offset, limit int, dto types.QueryDto) error {
	spots, err := s.SpotsRepository.GetSpotsByOffsetLimit(offset, limit)
	var aPartOfResults []models.Spot
	if err != nil {
		return err
	}
	for i := range spots {
		distance := util.GetDistanceFromLatLonInMeters(dto.Latitude, dto.Longitude, spots[i].Coordinates.Lat(), spots[i].Coordinates.Lng())
		spots[i].Distance = distance
		if dto.Type == types.Circle && distance <= dto.Radius {
			aPartOfResults = append(aPartOfResults, spots[i])
		}
		if dto.Type == types.Square && isBoundedBySquare(dto.Latitude, dto.Longitude, spots[i].Coordinates.Lat(), spots[i].Coordinates.Lng(), dto.Radius) {
			aPartOfResults = append(aPartOfResults, spots[i])
		}
	}
	dataChan <- aPartOfResults
	return nil
}

func (s *Service) GetAllSpotsByCondition(dto types.QueryDto) ([]types.SpotDto, error) {
	numberOfSpot, err := s.SpotsRepository.GetNumberOfSpots()
	if err != nil {
		return nil, err
	}

	gap := int(math.Ceil(float64(numberOfSpot) / Batch))
	var results []models.Spot
	var dataChan = make(chan []models.Spot, gap)
	var numberOfCompletedRoutine int

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for page := 1; page <= gap; page++ {
		page := page
		go func() {
			err := s.getBatchSpotsByCondition(dataChan, (page-1)*Batch, Batch, dto)
			if err != nil {
				// push err via slack or email
			}
		}()
	}

	for {
		select {
		case data := <-dataChan:
			results = append(results, data...)
			numberOfCompletedRoutine++
			if numberOfCompletedRoutine == gap {
				return toSpotDto(results), err
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("context timeout, ran out of time")
		}
	}
}
