package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/andriidelzz/go-activity-tracker/internal/db"
	"github.com/andriidelzz/go-activity-tracker/internal/model"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(repo *db.Repository) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	r.POST("/events", func(c *gin.Context) {
		var input model.Event
		if err := c.BindJSON(&input); err != nil {
			slog.Error("Failed to bind JSON:", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := repo.CreateEvent(&input); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, input)
	})

	r.GET("/events", func(c *gin.Context) {
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			slog.Error("user_id is required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if userID <= 0 {
			slog.Error("Only positive ids allowed:")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}

		if err != nil {
			slog.Error("Invalid user_id:", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		events, err := repo.GetEvents(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, events)
	})

	return r
}
