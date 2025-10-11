package server

import (
	"github.com/andriidelzz/go-activity-tracker/internal/handler"
	"github.com/andriidelzz/go-activity-tracker/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterRoutes(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Cors())
	r.Use(middleware.RateLimiter())

	r.POST("/events", h.HandleCreateEvent)
	r.GET("/events", h.HandleGetEvents)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
