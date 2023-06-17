package handlers

import (
	"fmt"
	"github.com/kholodmv/go-service/cmd/storage"
	"net/http"
)

type GetAllHandler struct {
	repository storage.MetricRepository
}

func NewGetAllHandler(repository storage.MetricRepository) *GetAllHandler {
	return &GetAllHandler{
		repository: repository,
	}
}

func (mh *GetAllHandler) GetAllMetric(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")

	metrics := mh.repository.GetAllMetrics()

	var str string
	for _, metric := range metrics {
		str += fmt.Sprintf("%q : %v\n", metric.Name, metric.Value)
	}

	fmt.Fprint(res, str)
}
