// Package sender ...
package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/SerjRamone/metrius/internal/metrics"
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

	r, err := http.Post(url, "application/json", bytes.NewBuffer(data))
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
