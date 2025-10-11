package server

import (
	"github.com/andriidelzz/go-activity-tracker/internal/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterRoutes(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	})) // todo move to /middlewares

	r.POST("/events", h.HandleCreateEvent)
	r.GET("/events", h.HandleGetEvents)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
