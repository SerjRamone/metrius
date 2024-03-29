// Package sender ...
package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/internal/middlewares"
	"github.com/SerjRamone/metrius/pkg/logger"
	"github.com/SerjRamone/metrius/pkg/retry"
)

// metricsSender ...
type metricsSender struct {
	client  *http.Client
	sURL    string
	hashKey string
}

// request middleware transport
type hashTripper struct {
	hashKey string
}

// RoundTrip is RoundTripper implementation
// Adds HashSHA256 header to request headers
func (sender hashTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var buf bytes.Buffer
	// copy data from request body
	_, err := io.Copy(&buf, req.Body)
	if err != nil {
		logger.Error("request body copy error", zap.Error(err))
		return nil, err
	}
	// get bytes from buffer
	b := buf.Bytes()
	// set new body
	req.Body = io.NopCloser(bytes.NewReader(b))
	// calculate hash string
	b64Hash := middlewares.CalcHash(b, []byte(sender.hashKey))
	req.Header.Set("HashSHA256", b64Hash)
	logger.Info("calculated body hash", zap.String("hash", b64Hash))
	return http.DefaultTransport.RoundTrip(req)
}

// NewMetricsSender crates MetricsSender
func NewMetricsSender(sURL, hashKey string) *metricsSender {
	sender := metricsSender{
		sURL:    "http://" + sURL,
		hashKey: hashKey,
		client:  &http.Client{},
	}

	if sender.hashKey != "" {
		logger.Info("requests with hash header")
		sender.client.Transport = hashTripper{hashKey: hashKey}
	}
	return &sender
}

// Send whole Collection
func (sender *metricsSender) Send(collections []metrics.Collection) error {
	for _, c := range collections {

		for _, m := range c {
			r, err := sender.sendMetrics(m)
			if err != nil {
				return err
			}

			r.Body.Close()
		}

		time.Sleep(50 * time.Millisecond)
	}

	return nil
}

// Worker is a async sender
func (sender *metricsSender) Worker(doneCh chan struct{}, jobCh chan []metrics.Collection) {
	logger.Info("worker started")
	for {
		select {
		case collections := <-jobCh:
			logger.Info("worker recived new collections")
			err := sender.SendBatch(collections)
			if err != nil {
				logger.Error("async SendBatch error", zap.Error(err))
			}
		case <-doneCh:
			logger.Info("worker recived done signal")
			return
		}
	}
}

// SendBatch sends metrics in batches
func (sender *metricsSender) SendBatch(collections []metrics.Collection) error {
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
		url := fmt.Sprintf("%s/updates/", sender.sURL)
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

		r, err := sender.client.Do(req)
		if err != nil {
			logger.Error("first request error", zap.Error(err))
			var netError net.Error
			if errors.As(err, &netError) {
				err = retry.WithBackoff(func() error {
					retrResp, retrErr := sender.client.Do(req)
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

// single request
func (sender *metricsSender) sendMetrics(m metrics.CollectionItem) (*http.Response, error) {
	url := fmt.Sprintf("%s/update/", sender.sURL)

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
		log.Println("metrics encode error", err)
	}

	var b bytes.Buffer
	gw := gzip.NewWriter(&b)
	if _, err = gw.Write(data); err != nil {
		return nil, err
	}
	if err = gw.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	r, err := sender.client.Do(req)
	if err != nil {
		var netError net.Error
		if errors.As(err, &netError) {
			err = retry.WithBackoff(func() error {
				retrResp, retrErr := sender.client.Do(req)
				if retrErr != nil {
					_, retrErr = io.Copy(io.Discard, retrResp.Body)
					if retrErr != nil {
						return retrErr
					}
				}
				defer retrResp.Body.Close()
				return retrErr
			}, 3)
			return r, err
		}
		logger.Error("send metrics error", zap.Error(err))
		return nil, err
	}
	defer r.Body.Close()

	_, err = io.Copy(io.Discard, r.Body)
	if err != nil {
		return nil, err
	}

	return r, nil
}
