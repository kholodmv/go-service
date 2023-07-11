package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/storage"
	"github.com/kholodmv/go-service/internal/gzip"
	"github.com/kholodmv/go-service/internal/logger"
	"github.com/kholodmv/go-service/internal/models"
	"os"
)

type Handler struct {
	router     chi.Router
	repository storage.MetricRepository
}

func NewHandler(router chi.Router, repository storage.MetricRepository, filename string, restore bool) *Handler {
	h := &Handler{
		repository: repository,
		router:     router,
	}

	if restore {
		file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Printf("Сan not open file: %s\n", err)
		}
		defer file.Close()

		var metrics []models.Metrics

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&metrics)
		if err != nil {
			fmt.Printf("Сan not restore data: %s\n", err)
		}

		for _, metric := range metrics {
			if metric.MType == "gauge" {
				h.repository.AddGauge(*metric.Value, metric.ID)
			} else if metric.MType == "counter" {
				h.repository.AddCounter(*metric.Delta, metric.ID)
			}
		}
	}

	h.router.Use(gzip.GzipHandler)
	h.router.Use(logger.LoggerHandler)
	return h
}

func (mh *Handler) RegisterRoutes(router *chi.Mux) {
	router.Post("/update/{type}/{name}/{value}", mh.UpdateMetric)
	router.Get("/value/{type}/{name}", mh.GetValueMetric)
	router.Get("/", mh.GetAllMetric)

	router.Post("/value/", mh.GetJSONMetric)
	router.Post("/update/", mh.UpdateJSONMetric)
}
