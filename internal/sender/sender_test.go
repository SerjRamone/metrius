package sender

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SerjRamone/metrius/internal/metrics"
)

func TestSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/update/", r.URL.String())
		assert.Equal(t, "gzip", r.Header.Get("Content-Encoding"))

		gr, err := gzip.NewReader(r.Body)
		assert.NoError(t, err)
		defer gr.Close()

		body, _ := io.ReadAll(gr)
		assert.Equal(t, `{"id":"Alloc","type":"gauge","value":134024}`, string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// @todo json response
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
	var pubKey []byte
	sender := NewMetricsSender(u.Host, "testkey", pubKey)
	err = sender.Send(c)
	assert.NoError(t, err)
	// assert.Equal(t, http.StatusOK, resp.StatusCode)
}
