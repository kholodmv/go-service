package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/handlers"
	"github.com/kholodmv/go-service/cmd/metrics"
	dataBase "github.com/kholodmv/go-service/internal/db"
	"github.com/kholodmv/go-service/internal/logger"
	"github.com/kholodmv/go-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetAllMetric(t *testing.T) {
	router := chi.NewRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	type want struct {
		url             string
		status          int
		contentType     string
		responseGauge   string
		responseCounter string
	}

	var tests = []struct {
		name string
		want want
	}{
		{"StatusOK test #1 - response return metrics",
			want{
				url:             "/",
				status:          http.StatusOK,
				contentType:     "text/plain",
				responseGauge:   `"test_gauge_metric" : 56.4`,
				responseCounter: `"test_counter_metric" : 5`,
			},
		},
	}

	log := logger.Initialize()
	storage := storage.NewMemoryStorage()
	db := dataBase.NewStorage("")
	storage.AddMetric(metrics.Gauge, 56.4, "test_gauge_metric")
	storage.AddMetric(metrics.Counter, int64(5), "test_counter_metric")
	getAllHandler := handlers.NewHandler(router, storage, db, *log)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			getAllHandler.GetAllMetric(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.status, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			parts := strings.Split(string(resBody), "\n")

			require.NoError(t, err)
			assert.Equal(t, parts[0], tt.want.responseGauge)
			assert.Equal(t, parts[1], tt.want.responseCounter)
		})
	}
}

func TestGetValueMetric(t *testing.T) {
	router := chi.NewRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	type want struct {
		url         string
		status      int
		contentType string
	}

	var tests = []struct {
		name string
		want want
	}{
		{"Test #1 - not found metric",
			want{
				url:         "/value/gauge/nameGaugeMetricKek",
				status:      http.StatusNotFound,
				contentType: "text/plain",
			},
		},
	}

	log := logger.Initialize()
	storage := storage.NewMemoryStorage()
	storage.AddMetric(metrics.Gauge, 56.4, "nameGaugeMetric")
	storage.AddMetric(metrics.Counter, int64(5), "nameCounterMetric")
	db := dataBase.NewStorage("")
	getValueHandler := handlers.NewHandler(router, storage, db, *log)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.want.url, nil)
			w := httptest.NewRecorder()
			getValueHandler.GetValueMetric(w, request)

			res := w.Result()
			assert.Equal(t, tt.want.status, res.StatusCode)
			defer res.Body.Close()
		})
	}
}
