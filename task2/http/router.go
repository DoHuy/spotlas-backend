package http

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	Router     *gin.Engine
	Controller *Controller
}

func NewRouter(router *gin.Engine, con *Controller) *Router {
	return &Router{
		Router:     router,
		Controller: con,
	}
}
func (r *Router) Init() {
	r.Router.Use(r.Controller.Options)
	api := r.Router.Group("/api")
	api.GET("/spots", r.Controller.GetSpots)
	return
}
