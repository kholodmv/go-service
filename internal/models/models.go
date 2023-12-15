// Package models contains information about models.
package models

import (
	"fmt"
	"github.com/kholodmv/go-service/cmd/metrics"
	"github.com/kholodmv/go-service/proto"
	"reflect"
)

// Metrics struct include information about metrics.
type Metrics struct {
	ID    string   `json:"id"`              // metric name
	MType string   `json:"type"`            // parameter taking the value gauge or counter
	Delta *int64   `json:"delta,omitempty"` // metric value in case of transfer counter
	Value *float64 `json:"value,omitempty"` // metric value in case of transmitting gauge
}

func ReadStruct(st interface{}) ([]Metrics, error) {
	val := reflect.ValueOf(st)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	result := make([]Metrics, 0)
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)

		switch f.Kind() {
		case reflect.Struct:
			r, err := ReadStruct(f)
			if err != nil {
				return nil, err
			}
			result = append(result, r...)
		case reflect.Pointer:
			if f.Elem().CanInterface() {
				r, err := ReadStruct(f.Elem().Interface())
				if err != nil {
					return nil, err
				}
				result = append(result, r...)
			}
		case reflect.Slice:
			for j := 0; j < f.Len(); j++ {
				v := f.Index(j).Float()
				result = append(result, Metrics{
					ID:    fmt.Sprintf("%v%v", val.Type().Field(i).Name, j+1),
					MType: metrics.Gauge,
					Value: &v,
				})
			}
		case reflect.Uint64:
			v := int64(f.Uint())
			result = append(result, Metrics{
				ID:    val.Type().Field(i).Name,
				MType: metrics.Counter,
				Delta: &v,
			})
		case reflect.Int64:
			v := f.Int()
			result = append(result, Metrics{
				ID:    val.Type().Field(i).Name,
				MType: metrics.Counter,
				Delta: &v,
			})
		case reflect.Float64:
			v := f.Float()
			result = append(result, Metrics{
				ID:    val.Type().Field(i).Name,
				MType: metrics.Gauge,
				Value: &v,
			})
		case reflect.Chan:
			continue
		default:
			return nil, fmt.Errorf("ivalid metric Kind: %v %v", f.Kind(), val.Type().Field(i).Name)
		}
	}

	return result, nil
}

func (m *Metrics) ToProto() ([]*proto.Metric, error) {
	metrics, err := ReadStruct(m)
	if err != nil {
		return nil, err
	}

	res := make([]*proto.Metric, 0, len(metrics))
	for i := range metrics {
		res = append(res, &proto.Metric{
			Id:    metrics[i].ID,
			Type:  metrics[i].MType,
			Delta: *metrics[i].Delta,
			Value: *metrics[i].Value,
		})
	}

	return res, nil
}
