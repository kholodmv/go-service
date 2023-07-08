package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/storage"
	"github.com/kholodmv/go-service/internal/gzip"
	"github.com/kholodmv/go-service/internal/logger"
)

type Handler struct {
	router     chi.Router
	repository storage.MetricRepository
}

func NewHandler(router chi.Router, repository storage.MetricRepository) *Handler {
	h := &Handler{
		repository: repository,
		router:     router,
	}
	h.router.Use(gzip.GzipHandle)
	return h
}

func (mh *Handler) RegisterRoutes(router *chi.Mux) {
	router.Post("/update/{type}/{name}/{value}", logger.RequestLogger(mh.UpdateMetric))
	router.Get("/value/{type}/{name}", logger.RequestLogger(mh.GetValueMetric))
	router.Get("/", logger.RequestLogger(mh.GetAllMetric))

	router.Post("/value/", logger.RequestLogger(mh.GetJSONMetric))
	router.Post("/update/", logger.RequestLogger(mh.UpdateJSONMetric))
}
