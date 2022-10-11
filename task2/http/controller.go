package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"spotlas/types"
	"spotlas/util"
)

type Controller struct {
	SpotService *Service
}

func NewController(spotService *Service) *Controller {
	return &Controller{
		SpotService: spotService,
	}
}

func (con *Controller) Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}

func (con *Controller) GetSpots(c *gin.Context) {
	queryDto := &types.QueryDto{}
	if err := util.ConvertUrlQueryToStruct(queryDto, c); err != nil {
		// log.info simplify by println
		util.BuildErrorResponse(c, http.StatusBadRequest, err, nil)
		return
	}
	spots, err := con.SpotService.GetAllSpotsByCondition(*queryDto)
	if err != nil {
		fmt.Println("con.SpotService.GetAllSpotsByCondition.Error =>", err)
		util.BuildErrorResponse(c, http.StatusInternalServerError, nil, nil)
		return
	}
	util.BuildSuccessResponse(c, 200, spots)
}
