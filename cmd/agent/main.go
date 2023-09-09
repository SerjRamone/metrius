// Package main ...
package main

import (
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/SerjRamone/metrius/internal/metrics"
)

var (
	pollInterval   = 2
	reportInterval = 10
)

const ServerURL = "http://localhost:8080"

func main() {
	collections := make(map[int64]metrics.Collection, 4)
	m := runtime.MemStats{}
	reportedAt := time.Now()
	polledAt := time.Now()

	for {
		// collect metrics
		if seconds := int((time.Since(polledAt)).Seconds()); seconds >= pollInterval {
			// getting metrics from runtime
			runtime.ReadMemStats(&m)

			// add metrics to collection
			c := metrics.NewCollection(m)
			collections[time.Now().UnixMicro()] = c
			log.Println("🗄  metrics added. Collections.len() is: ", len(collections))

			polledAt = time.Now()
		}

		// send metrics
		if seconds := int((time.Since(reportedAt)).Seconds()); seconds >= reportInterval {
			if len(collections) > 0 {
				for k, c := range collections {
					postCollection(c)
					log.Println("📨  sended")
					delete(collections, k)
					time.Sleep(50 * time.Millisecond)
				}
			}

			reportedAt = time.Now()
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// send whole Collection
func postCollection(c metrics.Collection) {
	for _, m := range c {
		r, err := postMetrics(ServerURL, m)
		if err != nil {
			log.Println(err)
		}
		r.Body.Close()
	}
}

// single request
func postMetrics(sURL string, m map[string]string) (*http.Response, error) {
	url := sURL + "/update/" + m["type"] + "/" + m["name"] + "/" + m["value"]
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
