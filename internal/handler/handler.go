package handler

import (
	"github.com/andriidelzz/go-activity-tracker/internal/repository"
)

type Handler struct {
	repo repository.RepositoryInterface
}

func NewHandler(repo repository.RepositoryInterface) *Handler {
	return &Handler{repo: repo}
}
