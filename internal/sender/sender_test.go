package sender

import (
	"io"
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
		assert.Equal(t, "/update/", r.URL.String())

		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, `{"id":"Alloc","type":"gauge","value":134024}`, string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK\n"))
	}))
	defer server.Close()

	c := []metrics.Collection{
		{
			metrics.CollectionItem{Name: "Alloc", Type: "gauge", Value: 134024},
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
