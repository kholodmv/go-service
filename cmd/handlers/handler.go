package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/storage"
)

type Handler struct {
	repository storage.MetricRepository
}

func NewHandler(repository storage.MetricRepository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (mh *Handler) RegisterRoutes(router *chi.Mux) {
	router.Post("/update/{type}/{name}/{value}", mh.UpdateMetric)
	router.Get("/value/{type}/{name}", mh.GetValueMetric)
	router.Get("/", mh.GetAllMetric)
}
