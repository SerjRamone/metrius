package sender

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/update/gauge/Alloc/134024.000000", r.URL.String())

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK\n"))
	}))
	defer server.Close()

	c := []metrics.Collection{
		{
			metrics.CollectionItem{Name: "Alloc", Variation: "gauge", Value: 134024},
		},
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
