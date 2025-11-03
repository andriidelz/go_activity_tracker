package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/andriidelzz/go-activity-tracker/internal/metrics"
	"github.com/andriidelzz/go-activity-tracker/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleCreateEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var input model.Event
	if err := c.BindJSON(&input); err != nil {
		slog.Error("Failed to bind JSON:", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateEvent(ctx, &input); err != nil {
		slog.Error("Failed to create event:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics.EventsTotal.WithLabelValues(input.Type).Inc()

	c.JSON(http.StatusCreated, input)
}

func (h *Handler) HandleGetEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		slog.Error("user_id is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		slog.Error("Invalid user_id:", "error", err, "user_id", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	events, err := h.repo.GetEvents(ctx, userID)
	if err != nil {
		slog.Error("Failed to get events:", "error", err, "user_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *Handler) HandleGetStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	stats, err := h.repo.GetStats(ctx)
	if err != nil {
		slog.Error("Failed to get stats:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) HandleAggregate(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	if err := h.repo.AggregateLastPeriod(ctx); err != nil {
		slog.Error("Aggregation failed:", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Aggregation completed"})
}
