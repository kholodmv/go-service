package main

import (
	"bytes"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/handlers"
	"github.com/kholodmv/go-service/internal/logger"
	dataBase "github.com/kholodmv/go-service/internal/store"
	"net/http"
	"net/http/httptest"
)

func ExampleMyHandler_HandleRequest() {
	router := chi.NewRouter()
	log := logger.Initialize()
	storage := dataBase.NewMemoryStorage()
	h := handlers.NewHandler(router, storage, *log, "")
	w := httptest.NewRecorder()

	b := []byte(`{"type": "gauge", "value": 10, "id": "metric2"}`)
	r := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewReader(b))
	h.UpdateJSONMetric(w, r)
	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	b = []byte(`[{"type": "counter", "delta": 10, "id": "metric3"}]`)
	r = httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewReader(b))
	h.UpdatesMetrics(w, r)
	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	b = []byte(`{"type": "counter", "delta": 10, "id": "metric3"}`)
	r = httptest.NewRequest(http.MethodGet, "/value/", bytes.NewReader(b))
	h.GetJSONMetric(w, r)
	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	// Output:
	// 200
	// {"id":"metric2","type":"gauge","value":10}
	// 200
	// {"id":"metric2","type":"gauge","value":10}
	// 200
	// {"id":"metric2","type":"gauge","value":10}{"id":"metric3","type":"counter","delta":10}
}
