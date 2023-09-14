package sender

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/update/gauge/Alloc/134024", r.URL.String())

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK\n"))
	}))
	defer server.Close()

	collection := []map[string]string{
		{
			"name":  "Alloc",
			"type":  "gauge",
			"value": "134024",
		},
	}
	c := map[int64]metrics.Collection{
		time.Now().UnixMicro(): collection,
	}
	u, err := url.Parse(server.URL)
	if err != nil {
		log.Fatal("can't parse testserver url", err)
	}
	sender := NewMetricsSender(u.Host)
	err = sender.Send(c)
	assert.NoError(t, err)
	// assert.Equal(t, http.StatusOK, resp.StatusCode)
}
