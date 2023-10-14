// Package sender ...
package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/pkg/logger"
	"go.uber.org/zap"
)

// metricsSender ...
type metricsSender struct {
	sURL string
}

// NewMetricsSender crates MetricsSender
func NewMetricsSender(sURL string) *metricsSender {
	return &metricsSender{
		sURL: "http://" + sURL,
	}
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

// SendBatch sends metrics in batches
func (sender *metricsSender) SendBatch(collections []metrics.Collection) error {
	batch := []metrics.Metrics{}
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
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "gzip")
		if err != nil {
			logger.Error("request object creation error", zap.Error(err))
			return err
		}

		client := &http.Client{}
		r, err := client.Do(req)
		if err != nil {
			logger.Error("doing request error", zap.Error(err))
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	_, err = io.Copy(io.Discard, r.Body)
	if err != nil {
		return nil, err
	}

	return r, nil
}
