package util

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	"math"
	"net/http"
	"spotlas/types"
)

var decoder = schema.NewDecoder()

const R = 6371 * 1000 // Radius of the earth in meters
func GetDistanceFromLatLonInMeters(lat1, lon1, lat2, lon2 float64) float64 {
	var dLat = (lat2 - lat1) * (math.Pi / 180)
	var dLon = (lon2 - lon1) * (math.Pi / 180)
	var a = math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*
			math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	var d = R * c // Distance in meters
	return d
}

func MaxLatLongOnBearing(centerLat, centerLong, bearing, distance float64) types.Coordinate {
	var lonRads = centerLong * math.Pi / 180
	var latRads = centerLat * math.Pi / 180
	var bearingRads = bearing * math.Pi / 180
	var maxLatRads = math.Asin(math.Sin(latRads)*math.Cos(distance/R) + math.Cos(latRads)*math.Sin(distance/R)*math.Cos(bearingRads))
	var maxLonRads = lonRads + math.Atan2(math.Sin(bearingRads)*math.Sin(distance/R)*math.Cos(latRads), math.Cos(distance/R)-math.Sin(latRads)*math.Sin(maxLatRads))

	var maxLat = maxLatRads * 180 / math.Pi
	var maxLong = maxLonRads * 180 / math.Pi

	return types.Coordinate{Lat: maxLat, Long: maxLong}
}

func ConvertUrlQueryToStruct(obj interface{}, c *gin.Context) error {
	if err := decoder.Decode(obj, c.Request.URL.Query()); err != nil {
		return err
	}
	return nil
}

func BuildErrorResponse(ctx *gin.Context, status int, err error, body interface{}) {
	if err == nil {
		BuildStandardResponse(ctx, http.StatusInternalServerError, body, ResponseMeta{Code: int64(status), Message: "Internal Server Error"})
	}
	BuildStandardResponse(ctx, status, body, ResponseMeta{Code: int64(status), Message: err.Error()})
}

func BuildSuccessResponse(ctx *gin.Context, status int, body interface{}) {
	BuildStandardResponse(ctx, status, body, ResponseMeta{Code: int64(status), Message: "Successfully"})
}

func BuildStandardResponse(ctx *gin.Context, status int, body interface{}, meta interface{}) {
	ctx.JSON(status, response{Data: body, Meta: meta})
}

type response struct {
	Data interface{} `json:"data,omitempty"`
	Meta interface{} `json:"meta,omitempty"`
}

type ResponseMeta struct {
	Code    interface{} `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
}
