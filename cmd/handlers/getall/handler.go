package getall

import (
	"fmt"
	"github.com/kholodmv/go-service/cmd/common"
	"github.com/kholodmv/go-service/cmd/storage"
	"net/http"
)

type Handler struct {
	repository storage.MetricRepository
}

func NewHandler(repository storage.MetricRepository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (m *Handler) GetAllMetric(res http.ResponseWriter, req *http.Request) {
	common.CheckGetHTTPMethod(res, req)
	res.Header().Set("Content-Type", "text/plain")

	metrics := m.repository.GetAllMetrics()

	var str string
	for _, metric := range metrics {
		str += fmt.Sprintf("%q : %v\n", metric.Name, metric.Value)
	}

	fmt.Fprint(res, str)
}
