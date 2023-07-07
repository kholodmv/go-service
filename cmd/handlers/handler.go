package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/storage"
	"github.com/kholodmv/go-service/internal/logger"
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
	router.Post("/update/{type}/{name}/{value}", logger.RequestLogger(mh.UpdateMetric))
	router.Get("/value/{type}/{name}", logger.RequestLogger(mh.GetValueMetric))
	router.Get("/", logger.RequestLogger(mh.GetAllMetric))

	router.Post("/value/", logger.RequestLogger(mh.GetJSONMetric))
	router.Post("/update/", logger.RequestLogger(mh.UpdateJSONMetric))
}
