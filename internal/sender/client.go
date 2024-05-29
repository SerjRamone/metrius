package sender

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/pkg/logger"
	pb "github.com/SerjRamone/metrius/pkg/metrius_v1"
	"github.com/SerjRamone/metrius/pkg/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// APIClient ...
type APIClient interface {
	Do(metrics.CollectionItem) error
	DoBatch([]metrics.Collection) error
}

// HTTPApiClient ...
type HTTPApiClient struct {
	client  *http.Client
	sURL    string
	localIP string
}

// NewHTTPApiClient ...
func NewHTTPApiClient(sURL, localIP string) *HTTPApiClient {
	return &HTTPApiClient{
		client:  &http.Client{},
		sURL:    sURL,
		localIP: localIP,
	}
}

// Do sends metrics to server
func (c *HTTPApiClient) Do(m metrics.CollectionItem) error {
	url := fmt.Sprintf("%s/update/", c.sURL)

	item := metrics.Metrics{
		ID:    m.Name,
		MType: m.Type,
	}

	switch m.Type {
	case "gauge":
		item.Value = &m.Value
	case "counter":
		tmp := int64(m.Value)
		item.Delta = &tmp
	}

	data, err := json.Marshal(item)
	if err != nil {
		logger.Error("metrics encode error", zap.Error(err))
		return err
	}

	var b bytes.Buffer
	gw := gzip.NewWriter(&b)
	if _, err = gw.Write(data); err != nil {
		return err
	}
	if err = gw.Close(); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("X-Real-IP", c.localIP)

	r, err := c.client.Do(req)
	if err != nil {
		var netError net.Error
		if errors.As(err, &netError) {
			err = retry.WithBackoff(func() error {
				retrResp, retrErr := c.client.Do(req)
				if retrErr != nil {
					_, retrErr = io.Copy(io.Discard, retrResp.Body)
					if retrErr != nil {
						return retrErr
					}
				}
				defer retrResp.Body.Close()
				return retrErr
			}, 3)
			return err
		}
		logger.Error("send metrics error", zap.Error(err))
		return err
	}
	defer r.Body.Close()

	_, err = io.Copy(io.Discard, r.Body)
	if err != nil {
		return err
	}

	return nil
}

// DoBatch sends metrics to server
func (c *HTTPApiClient) DoBatch(collections []metrics.Collection) error {
	batch := make([]metrics.Metrics, 0, 200)
	// collect batch of metrics.Metrics
	for _, c := range collections {
		for _, m := range c {
			item := metrics.Metrics{
				ID:    m.Name,
				MType: m.Type,
			}

			switch m.Type {
			case "gauge":
				value := m.Value
				item.Value = &value
			case "counter":
				tmp := int64(m.Value)
				item.Delta = &tmp
			}

			batch = append(batch, item)
		}
	}

	if len(batch) > 0 {
		url := fmt.Sprintf("%s/updates/", c.sURL)
		data, err := json.Marshal(batch)
		if err != nil {
			logger.Error("metrics encode error", zap.Error(err))
			return err
		}

		var b bytes.Buffer
		gw := gzip.NewWriter(&b)
		if _, err = gw.Write(data); err != nil {
			logger.Error("gzipped write data error", zap.Error(err))
			return err
		}
		if err = gw.Close(); err != nil {
			return err
		}

		req, err := http.NewRequest("POST", url, &b)
		if err != nil {
			logger.Error("request object creation error", zap.Error(err))
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("X-Real-IP", c.localIP)

		r, err := c.client.Do(req)
		if err != nil {

			logger.Error("first request error", zap.Error(err))
			var netError net.Error
			if errors.As(err, &netError) {
				err = retry.WithBackoff(func() error {
					retrResp, retrErr := c.client.Do(req)
					if retrErr != nil {
						return retrErr
					}
					_, retrErr = io.Copy(io.Discard, retrResp.Body)
					if retrErr != nil {
						logger.Error("io.Copy error", zap.Error(retrErr))
					}
					defer retrResp.Body.Close()
					return retrErr
				}, 3)
				return err
			}
			logger.Error("send metrics in batch error", zap.Error(err))
			return err
		}
		defer r.Body.Close()

		_, err = io.Copy(io.Discard, r.Body)
		if err != nil {
			logger.Error("io.Copy error", zap.Error(err))
			return err
		}
	}

	return nil

}

// GRPCApiClient ...
type GRPCApiClient struct {
	conn   *grpc.ClientConn
	client pb.MetricsServiceClient
}

// NewGRPCApiClient creates ApiClient
func NewGRPCApiClient(a string) (*GRPCApiClient, error) {
	conn, err := grpc.Dial(a, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("grpc.Dial error", zap.Error(err))
		return nil, err
	}
	return &GRPCApiClient{
		conn:   conn,
		client: pb.NewMetricsServiceClient(conn),
	}, nil
}

// Do sends metrics to server
func (c *GRPCApiClient) Do(m metrics.CollectionItem) error {
	item := &pb.Metrics{
		Id: m.Name,
	}

	switch m.Type {
	case "gauge":
		item.Value = m.Value
		item.Type = pb.Metrics_GAUGE
	case "counter":
		tmp := int64(m.Value)
		item.Delta = tmp
		item.Type = pb.Metrics_COUNTER
	}

	resp, err := c.client.Update(context.TODO(), &pb.UpdateRequest{Metrics: item})
	if err != nil {
		return fmt.Errorf("grps Update error: %w", err)
	}
	if resp.Error != "" {
		return fmt.Errorf("grps Update error: %w", errors.New(resp.Error))
	}

	return nil
}

// DoBatch sends metrics to server
func (c *GRPCApiClient) DoBatch(collections []metrics.Collection) error {
	batch := make([]*pb.Metrics, 0, 200)
	// collect batch of metrics.Metrics
	for _, c := range collections {
		for _, m := range c {
			item := &pb.Metrics{
				Id: m.Name,
			}

			switch m.Type {
			case "gauge":
				item.Value = m.Value
				item.Type = pb.Metrics_GAUGE
			case "counter":
				tmp := int64(m.Value)
				item.Delta = tmp
				item.Type = pb.Metrics_COUNTER
			}

			batch = append(batch, item)
		}
	}

	if len(batch) > 0 {
		resp, err := c.client.BatchUpdate(context.TODO(), &pb.BatchUpdateRequest{Metrics: batch})
		if err != nil {
			return fmt.Errorf("grps BatchUpdate error: %w", err)
		}
		if resp.Error != "" {
			return fmt.Errorf("grps BatchUpdate error: %w", errors.New(resp.Error))
		}
	}

	return nil
}
