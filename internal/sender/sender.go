// Package sender ...
package sender

import (
	"fmt"
	"io"
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
func (sender *metricsSender) Send(collections map[int64]metrics.Collection) error {
	for k, c := range collections {

		for _, m := range c {
			r, err := sender.sendMetrics(m)
			if err != nil {
				return err
			}

			r.Body.Close()
		}

		delete(collections, k)
		time.Sleep(50 * time.Millisecond)
	}

	return nil
}

// single request
func (sender *metricsSender) sendMetrics(m map[string]string) (*http.Response, error) {
	url := fmt.Sprintf("%s/update/%s/%s/%s", sender.sURL, m["type"], m["name"], m["value"])
	r, err := http.Post(url, "text/plain", nil)
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
