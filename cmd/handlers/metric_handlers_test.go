package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/kholodmv/go-service/cmd/metrics"
	"github.com/kholodmv/go-service/internal/logger"
	dataBase "github.com/kholodmv/go-service/internal/store"
)

func BenchmarkUpdateJSONMetric(b *testing.B) {
	storage := dataBase.NewMemoryStorage()
	storage.AddMetric(context.TODO(), metrics.Counter, int64(5), "metric1")
	storage.AddMetric(context.TODO(), metrics.Gauge, 56.4, "metric2")
	router := chi.NewRouter()
	log := logger.Initialize()
	h := NewHandler(router, storage, *log, "")

	type want struct {
		code int
	}
	benchmarks := []struct {
		name string
		body []byte
		want want
	}{
		{
			name: "UpdateJSONMetric benchmark - type counter",
			body: []byte(`{"type": "counter", "delta": 10, "id": "metric1"}`),
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "UpdateJSONMetric benchmark - type gauge",
			body: []byte(`{"type": "gauge", "value": 10, "id": "metric2"}`),
			want: want{
				code: http.StatusOK,
			},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()

				r := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewReader(bm.body))

				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				b.StartTimer()

				h.UpdateJSONMetric(w, r)

				resp := w.Result()
				resp.Body.Close()
				assert.Equal(b, bm.want.code, resp.StatusCode)
			}
		})
	}
}

func BenchmarkUpdatesMetrics(b *testing.B) {
	storage := dataBase.NewMemoryStorage()
	storage.AddMetric(context.TODO(), metrics.Counter, int64(5), "metric3")
	storage.AddMetric(context.TODO(), metrics.Gauge, 56.4, "metric4")
	router := chi.NewRouter()
	log := logger.Initialize()
	h := NewHandler(router, storage, *log, "")

	type want struct {
		code int
	}
	benchmarks := []struct {
		name string
		body []byte
		want want
	}{
		{
			name: "UpdatesMetrics benchmark - type counter",
			body: []byte(`[{"type": "counter", "delta": 10, "id": "metric3"}]`),
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "UpdatesMetrics benchmark - type gauge",
			body: []byte(`[{"type": "gauge", "value": 10, "id": "metric4"}]`),
			want: want{
				code: http.StatusOK,
			},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()

				r := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewReader(bm.body))

				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				b.StartTimer()

				h.UpdatesMetrics(w, r)

				resp := w.Result()
				resp.Body.Close()
				assert.Equal(b, bm.want.code, resp.StatusCode)
			}
		})
	}
}

func BenchmarkGetJSONMetric(b *testing.B) {
	storage := dataBase.NewMemoryStorage()
	storage.AddMetric(context.TODO(), metrics.Counter, int64(5), "metric3")
	storage.AddMetric(context.TODO(), metrics.Gauge, 56.4, "metric4")
	router := chi.NewRouter()
	log := logger.Initialize()
	h := NewHandler(router, storage, *log, "")

	type want struct {
		code int
	}
	benchmarks := []struct {
		name string
		body []byte
		want want
	}{
		{
			name: "GetJSONMetric benchmark - type counter",
			body: []byte(`{"type": "counter", "delta": 10, "id": "metric3"}`),
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "GetJSONMetric benchmark - type gauge",
			body: []byte(`{"type": "gauge", "value": 10, "id": "metric4"}`),
			want: want{
				code: http.StatusOK,
			},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()

				r := httptest.NewRequest(http.MethodGet, "/value/", bytes.NewReader(bm.body))

				r.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				b.StartTimer()

				h.GetJSONMetric(w, r)

				resp := w.Result()
				resp.Body.Close()
				assert.Equal(b, bm.want.code, resp.StatusCode)
			}
		})
	}
}

func BenchmarkGetAllMetric(b *testing.B) {
	storage := dataBase.NewMemoryStorage()
	storage.AddMetric(context.TODO(), metrics.Counter, int64(5), "metric3")
	storage.AddMetric(context.TODO(), metrics.Gauge, 56.4, "metric4")
	router := chi.NewRouter()
	log := logger.Initialize()
	h := NewHandler(router, storage, *log, "")

	type want struct {
		code int
	}
	benchmarks := []struct {
		name string
		want want
	}{
		{
			name: "GetAllMetric benchmark - type counter",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "GetAllMetric benchmark - type gauge",
			want: want{
				code: http.StatusOK,
			},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()

				r := httptest.NewRequest(http.MethodGet, "/", nil)

				r.Header.Set("Content-Type", "text/html; charset=utf-8")
				w := httptest.NewRecorder()

				b.StartTimer()

				h.GetAllMetric(w, r)

				resp := w.Result()
				resp.Body.Close()
				assert.Equal(b, bm.want.code, resp.StatusCode)
			}
		})
	}
}
