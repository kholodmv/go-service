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
	"github.com/kholodmv/go-service/internal/middleware/subnet"
	"go.uber.org/zap"
	"net"

	"github.com/kholodmv/go-service/internal/store"
)

// Handler struct.
type Handler struct {
	router           chi.Router
	db               store.Storage
	log              zap.SugaredLogger
	key              string
	trustedSubnet    string
	cryptoPrivateKey *rsa.PrivateKey
}

// NewHandler creates a new instance of the handler structure.
func NewHandler(router chi.Router, db store.Storage, log zap.SugaredLogger, key string, trustedSubnet string, cryptoPrivateKey *rsa.PrivateKey) *Handler {
	h := &Handler{
		router:           router,
		db:               db,
		log:              log,
		key:              key,
		trustedSubnet:    trustedSubnet,
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
	if mh.trustedSubnet != "" {
		_, trustedNet, err := net.ParseCIDR(mh.trustedSubnet)
		if err != nil {
			mh.log.Error("failed parse CIDR", zap.Error(err))
		}
		mh.router.Use(subnet.WithCheckSubnet(trustedNet))
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
