package getvalue

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kholodmv/go-service/cmd/common"
	"github.com/kholodmv/go-service/cmd/storage"
	"io"
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

func (mh *Handler) GetValueMetric(res http.ResponseWriter, req *http.Request) {
	common.CheckGetHTTPMethod(res, req)
	res.Header().Set("Content-Type", "text/plain")

	typeMetric := chi.URLParam(req, "type")
	name := chi.URLParam(req, "name")

	var value interface{}
	var ok bool

	if typeMetric == common.Gauge {
		value, ok = mh.repository.GetValueGaugeMetric(name)
	}
	if typeMetric == common.Counter {
		value, ok = mh.repository.GetValueCounterMetric(name)
	}

	fmt.Println(value)
	if !ok {
		http.NotFound(res, req)
		return
	}
	strValue := fmt.Sprintf("%v", value)

	io.WriteString(res, strValue)
	res.WriteHeader(http.StatusOK)
}
