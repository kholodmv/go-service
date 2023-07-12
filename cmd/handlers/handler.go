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

	return h
}

func (mh *Handler) RegisterRoutes(router *chi.Mux) {
	mh.router.Use(gzip.GzipHandler)
	mh.router.Use(logger.LoggerHandler)

	router.Post("/update/{type}/{name}/{value}", mh.UpdateMetric)
	router.Get("/value/{type}/{name}", mh.GetValueMetric)
	router.Get("/", mh.GetAllMetric)

	router.Post("/value/", mh.GetJSONMetric)
	router.Post("/update/", mh.UpdateJSONMetric)
}
