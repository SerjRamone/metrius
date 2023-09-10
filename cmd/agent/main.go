// Package main is agent main package
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/SerjRamone/metrius/internal/metrics"
)

var (
	reportedAt  time.Time
	polledAt    time.Time
	memStat     runtime.MemStats
	collections map[int64]metrics.Collection
)

func main() {
	// parse flags and envs
	parseConfig()

	collections = make(map[int64]metrics.Collection, 4)
	memStat = runtime.MemStats{}
	reportedAt = time.Now()
	polledAt = time.Now()

	for {
		// collect metrics
		if seconds := int((time.Since(polledAt)).Seconds()); seconds >= pollInterval {
			// getting metrics from runtime
			runtime.ReadMemStats(&memStat)

			// add metrics to collection
			c := metrics.NewCollection(memStat)
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
		r, err := postMetrics("http://"+serverAddress, m)
		if err != nil {
			log.Println(err)
			return
		}
		r.Body.Close()
	}
}

// single request
func postMetrics(sURL string, m map[string]string) (*http.Response, error) {
	url := fmt.Sprintf("%s/update/%s/%s/%s", sURL, m["type"], m["name"], m["value"])
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
