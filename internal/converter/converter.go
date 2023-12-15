package converter

import (
	"github.com/kholodmv/go-service/internal/models"
	pb "github.com/kholodmv/go-service/proto"
)

func ModelToProto(m models.Metrics) *pb.Metric {
	return &pb.Metric{
		Id:    m.ID,
		Type:  m.MType,
		Delta: *m.Delta,
		Value: *m.Value,
	}
}

func ProtoToModel(p *pb.Metric) models.Metrics {
	value := p.GetValue()
	delta := p.GetDelta()

	return models.Metrics{
		ID:    p.GetId(),
		MType: p.GetType(),
		Value: &value,
		Delta: &delta,
	}
}

func SliceModelToProto(m []models.Metrics) []*pb.Metric {
	res := make([]*pb.Metric, 0, len(m))
	for i := range m {
		res = append(res, ModelToProto(m[i]))
	}

	return res
}
