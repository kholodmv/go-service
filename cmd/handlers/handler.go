// Package handlers - handler.go - implements HTTP request handlers.
package handlers

import (
	"crypto/rsa"
	"github.com/go-chi/chi/v5"
	mw "github.com/go-chi/chi/v5/middleware"
	"github.com/kholodmv/go-service/internal/middleware/decrypt"
	"github.com/kholodmv/go-service/internal/middleware/gzip"
	"github.com/kholodmv/go-service/internal/middleware/hash"
	"github.com/kholodmv/go-service/internal/middleware/logger"
	"go.uber.org/zap"

	"github.com/kholodmv/go-service/internal/store"
)

// Handler struct.
type Handler struct {
	router           chi.Router
	db               store.Storage
	log              zap.SugaredLogger
	key              string
	cryptoPrivateKey *rsa.PrivateKey
}

// NewHandler creates a new instance of the handler structure.
func NewHandler(router chi.Router, db store.Storage, log zap.SugaredLogger, key string, cryptoPrivateKey *rsa.PrivateKey) *Handler {
	h := &Handler{
		router:           router,
		db:               db,
		log:              log,
		key:              key,
		cryptoPrivateKey: cryptoPrivateKey,
	}

	return h
}

// RegisterRoutes registers routes in the application.
func (mh *Handler) RegisterRoutes(router *chi.Mux) {
	mh.router.Use(gzip.GzipHandler)
	mh.router.Use(logger.LoggerHandler)
	if mh.key != "" {
		mh.router.Use(hash.HashHandler)
	}

	if mh.cryptoPrivateKey != nil {
		mh.router.Use(decrypt.WithRsaDecrypt(mh.cryptoPrivateKey))
	}

	router.Post("/update/{type}/{name}/{value}", mh.UpdateMetric)
	router.Get("/value/{type}/{name}", mh.GetValueMetric)
	router.Get("/", mh.GetAllMetric)
	router.Post("/value/", mh.GetJSONMetric)
	router.Post("/update/", mh.UpdateJSONMetric)
	router.Get("/ping", mh.DBConnection)
	router.Post("/updates/", mh.UpdatesMetrics)

	router.Mount("/debug", mw.Profiler())
}
