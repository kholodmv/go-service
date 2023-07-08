package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/handlers"
	"github.com/kholodmv/go-service/cmd/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	return resp, string(respBody)
}

func TestUpdateMetric(t *testing.T) {
	ts := httptest.NewServer(MetricRouter())
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
		{"StatusOK test #1",
			want{
				url:         "/update/gauge/NumForcedGC/517.33",
				status:      http.StatusOK,
				contentType: "text/plain",
			},
		},
		{"StatusOK test #2",
			want{
				url:         "/update/counter/nameMetric/1",
				status:      http.StatusOK,
				contentType: "text/plain",
			},
		},
		{"StatusOK test #3",
			want{
				url:         "/update/counter/nameMetric/1",
				status:      http.StatusOK,
				contentType: "text/plain",
			},
		},
		{"StatusBadRequest test #4 - incorrect value type for counter",
			want{
				url:         "/update/counter/nameMetric/1.4",
				status:      http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{"StatusBadRequest test #5 - incorrect value type for gauge",
			want{
				url:         "/update/gauge/NumForcedGC/-0s80.55",
				status:      http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{"StatusNotFound test #6 - incorrect url",
			want{
				url:         "/update",
				status:      http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{"StatusNotFound test #7 - incorrect url",
			want{
				url:         "/update/gauge",
				status:      http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{"StatusNotFound test #8 - incorrect url",
			want{
				url:         "/update/gauge/",
				status:      http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{"StatusNotFound test #9 - incorrect url",
			want{
				url:         "/update/counter",
				status:      http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{"StatusNotFound test #10 - incorrect url",
			want{
				url:         "/update/counter/",
				status:      http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{"StatusNotFound test #11 - empty value for type gauge",
			want{
				url:         "/update/gauge//",
				status:      http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{"StatusNotFound test #12 - empty value for type counter",
			want{
				url:         "/update/counter//",
				status:      http.StatusNotFound,
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		resp, _ := testRequest(t, ts, "POST", tt.want.url)
		resp.Body.Close()

		assert.Equal(t, tt.want.status, resp.StatusCode)
	}
}

func TestGetAllMetric(t *testing.T) {
	ts := httptest.NewServer(MetricRouter())
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
		{"StatusOK test #1 - response return metriccs",
			want{
				url:             "/",
				status:          http.StatusOK,
				contentType:     "text/plain",
				responseGauge:   `"test_gauge_metric" : 56.4`,
				responseCounter: `"test_counter_metric" : 5`,
			},
		},
	}

	router := chi.NewRouter()
	storage := storage.NewMemoryStorage()
	storage.AddGauge(56.4, "test_gauge_metric")
	storage.AddCounter(5, "test_counter_metric")
	getAllHandler := handlers.NewHandler(router, storage, "", true)

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
	ts := httptest.NewServer(MetricRouter())
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

	router := chi.NewRouter()
	storage := storage.NewMemoryStorage()
	storage.AddGauge(56.4, "nameGaugeMetric")
	storage.AddCounter(5, "nameCounterMetric")
	getValueHandler := handlers.NewHandler(router, storage, "", true)

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
