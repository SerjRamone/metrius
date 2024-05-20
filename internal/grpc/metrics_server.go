// Package grpc ...
package grpc

import (
	"context"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/internal/storage"
	"github.com/SerjRamone/metrius/pkg/logger"
	pb "github.com/SerjRamone/metrius/pkg/metrius_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MetricsServer ...
type MetricsServer struct {
	pb.UnimplementedMetricsServiceServer
	storage storage.Storage
}

// NewMetricsServer ...
func NewMetricsServer(store storage.Storage) *MetricsServer {
	return &MetricsServer{
		storage: store,
	}
}

// Update ...
func (s *MetricsServer) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	var response pb.UpdateResponse
	var err error
	switch in.Metrics.Type {
	case pb.Metrics_COUNTER:
		err = s.storage.SetCounter(ctx, in.Metrics.Id, metrics.Counter(in.Metrics.Delta))
		if err != nil {
			logger.Error("can't set counter", zap.String("ID", in.Metrics.Id), zap.Int64("delta", in.Metrics.Delta), zap.Error(err))
			return nil, status.Errorf(codes.Internal, "can't set counter. ID: %s, VALUE: %v", in.Metrics.Id, in.Metrics.Delta)
		}
	case pb.Metrics_GAUGE:
		err = s.storage.SetGauge(ctx, in.Metrics.Id, metrics.Gauge(in.Metrics.Value))
		if err != nil {
			logger.Error("can't set gauge", zap.String("ID", in.Metrics.Id), zap.Float64("delta", in.Metrics.Value), zap.Error(err))
			return nil, status.Errorf(codes.Internal, "can't set gauge. ID: %s, VALUE: %v", in.Metrics.Id, in.Metrics.Value)
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown metrics type: %v", in.Metrics.Type)
	}
	response.Metrics = in.Metrics

	return &response, nil
}

// BatchUpdate ...
func (s *MetricsServer) BatchUpdate(ctx context.Context, in *pb.BatchUpdateRequest) (*pb.BatchUpdateResponse, error) {
	var response pb.BatchUpdateResponse
	var batch []metrics.Metrics
	for _, m := range in.Metrics {
		var mType string
		switch m.Type {
		case pb.Metrics_GAUGE:
			mType = "gauge"
		case pb.Metrics_COUNTER:
			mType = "counter"
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unknown metrics type: %v", mType)
		}
		batch = append(batch, metrics.Metrics{ID: m.Id, Value: &m.Value, Delta: &m.Delta, MType: mType})
	}
	if err := s.storage.BatchUpsert(ctx, batch); err != nil {
		logger.Error("can't do batch upsert", zap.Error(err))
		return nil, status.Error(codes.Internal, "batach upsert error")
	}
	return &response, nil
}

// GetMetrics ...
func (s *MetricsServer) GetMetrics(ctx context.Context, in *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	var response pb.GetMetricsResponse

	switch in.Metrics.Type {
	case pb.Metrics_COUNTER:
		if m, ok := s.storage.Counter(ctx, in.Metrics.Id); ok {
			response.Metrics = &pb.Metrics{
				Id:    in.Metrics.Id,
				Delta: int64(m),
				Type:  pb.Metrics_COUNTER,
			}
		} else {
			logger.Error("counter not found", zap.String("ID", in.Metrics.Id))
			return nil, status.Errorf(codes.NotFound, "counter not found. ID: %s", in.Metrics.Id)
		}
	case pb.Metrics_GAUGE:
		if m, ok := s.storage.Gauge(ctx, in.Metrics.Id); ok {
			response.Metrics = &pb.Metrics{
				Id:    in.Metrics.Id,
				Value: float64(m),
				Type:  pb.Metrics_GAUGE,
			}
		} else {
			logger.Error("counter not found", zap.String("ID", in.Metrics.Id))
			return nil, status.Errorf(codes.NotFound, "counter not found. ID: %s", in.Metrics.Id)
		}
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unknown metrics type: %v", in.Metrics.Type)
	}

	return &response, nil
}
