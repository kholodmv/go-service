package core

import (
	"context"
	"github.com/kholodmv/go-service/internal/models"
	"github.com/kholodmv/go-service/internal/store"

	"github.com/kholodmv/go-service/internal/converter"
	pb "github.com/kholodmv/go-service/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer
	repo   store.Storage
	logger *zap.Logger
}

func NewMetricsServer(repo store.Storage, logger *zap.Logger) *MetricsServer {
	return &MetricsServer{
		repo:   repo,
		logger: logger,
	}
}

func (s *MetricsServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	result, err := s.repo.GetValueMetric(ctx, req.GetId(), req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to get metric with id %v: %v", req.GetId(), err)
	}
	m := result.(models.Metrics)
	return &pb.GetResponse{Metric: converter.ModelToProto(m)}, nil
}

func (s *MetricsServer) List(ctx context.Context, _ *emptypb.Empty) (*pb.ListResponse, error) {
	data, err := s.repo.GetAllMetrics(ctx, 100)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to get metrics: %v", err)
	}

	return &pb.ListResponse{Metrics: converter.SliceModelToProto(data)}, nil
}

func (s *MetricsServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	var result models.Metrics
	if _, err := s.repo.GetValueMetric(ctx, req.GetMetric().GetId(), req.GetMetric().GetId()); err != nil {
		var v interface{}
		err = s.repo.AddMetric(ctx, req.GetMetric().GetId(), v, req.GetMetric().GetType())
		if err != nil {
			s.logger.Error(err.Error())
			return nil, status.Errorf(codes.Internal, "unable to create metric %+v: %v", req.Metric, err)
		}
	} else {
		var v interface{}
		err = s.repo.AddMetric(ctx, req.GetMetric().GetId(), v, req.GetMetric().GetType())
		if err != nil {
			s.logger.Error(err.Error())
			return nil, status.Errorf(codes.Internal, "unable to update metric %+v: %v", req.Metric, err)
		}
	}

	return &pb.UpdateResponse{Result: converter.ModelToProto(result)}, nil
}
