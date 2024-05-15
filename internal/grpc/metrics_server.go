// Package grpc ...
package grpc

import (
	"context"
	"fmt"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/internal/storage"
	"github.com/SerjRamone/metrius/pkg/logger"
	pb "github.com/SerjRamone/metrius/pkg/metrius_v1"
	"go.uber.org/zap"
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
	_ = ctx
	switch in.Metrics.Type {
	case pb.Metrics_COUNTER:
		err = s.storage.SetCounter(in.Metrics.Id, metrics.Counter(in.Metrics.Delta))
		if err != nil {
			response.Error = fmt.Sprintf("can't set counter. ID: %s, VALUE: %v", in.Metrics.Id, in.Metrics.Delta)
			logger.Error("can't set counter", zap.String("ID", in.Metrics.Id), zap.Int64("delta", in.Metrics.Delta), zap.Error(err))
		}
	case pb.Metrics_GAUGE:
		err = s.storage.SetGauge(in.Metrics.Id, metrics.Gauge(in.Metrics.Value))
		if err != nil {
			response.Error = fmt.Sprintf("can't set gauge. ID: %s, VALUE: %v", in.Metrics.Id, in.Metrics.Value)
			logger.Error("can't set gauge", zap.String("ID", in.Metrics.Id), zap.Float64("delta", in.Metrics.Value), zap.Error(err))
		}
	default:
		response.Error = "unknown metrics type"
	}
	response.Metrics = in.Metrics

	return &response, err
}

// BatchUpdate ...
func (s *MetricsServer) BatchUpdate(ctx context.Context, in *pb.BatchUpdateRequest) (*pb.BatchUpdateResponse, error) {
	var response pb.BatchUpdateResponse
	var err error
	_ = ctx
	var batch []metrics.Metrics
loop:
	for _, m := range in.Metrics {
		var mType string
		switch m.Type {
		case pb.Metrics_GAUGE:
			mType = "gauge"
		case pb.Metrics_COUNTER:
			mType = "counter"
		default:
			response.Error = "unknown metrics type"
			break loop
		}
		batch = append(batch, metrics.Metrics{ID: m.Id, Value: &m.Value, Delta: &m.Delta, MType: mType})
	}
	if err := s.storage.BatchUpsert(batch); err != nil {
		logger.Error("can't do batch upsert", zap.Error(err))
	}
	return &response, err
}

// GetMetrics ...
func (s *MetricsServer) GetMetrics(ctx context.Context, in *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	var response pb.GetMetricsResponse
	var err error
	_ = ctx

	switch in.Metrics.Type {
	case pb.Metrics_COUNTER:
		if m, ok := s.storage.Counter(in.Metrics.Id); ok {
			response.Metrics = &pb.Metrics{
				Id:    in.Metrics.Id,
				Delta: int64(m),
				Type:  pb.Metrics_COUNTER,
			}
		} else {
			response.Error = fmt.Sprintf("counter not found. ID: %s", in.Metrics.Id)
		}
	case pb.Metrics_GAUGE:
		if m, ok := s.storage.Gauge(in.Metrics.Id); ok {
			response.Metrics = &pb.Metrics{
				Id:    in.Metrics.Id,
				Value: float64(m),
				Type:  pb.Metrics_GAUGE,
			}
		} else {
			response.Error = fmt.Sprintf("gauge not found. ID: %s", in.Metrics.Id)
		}
	default:
		response.Error = "unknown metrics type"
	}

	return &response, err
}
