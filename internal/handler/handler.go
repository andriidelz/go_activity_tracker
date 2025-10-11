package handler

import "github.com/andriidelzz/go-activity-tracker/internal/repository"

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}
