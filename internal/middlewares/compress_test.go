package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipCompressor(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	req.Header.Set("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("Lorem ipsum dolor sit amet"))
	})

	compressedHandler := GzipCompressor(handler)
	compressedHandler.ServeHTTP(rr, req)

	// Check if the content encoding header is set to gzip
	assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))

	// Check if the response body is gzipped
	body := rr.Body.Bytes()
	assert.Greater(t, len(body), 0) // Check if body has some content

	// Decompress the response body
	reader, err := gzip.NewReader(bytes.NewReader(body))
	assert.NoError(t, err)
	defer reader.Close()

	decompressedBody, err := io.ReadAll(reader)
	assert.NoError(t, err)

	assert.Equal(t, "Lorem ipsum dolor sit amet", string(decompressedBody))
}

func TestGzipCompressorWithoutGzipSupport(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	// Do not set Accept-Encoding header to simulate client without gzip support

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Lorem ipsum dolor sit amet"))
	})

	compressedHandler := GzipCompressor(handler)
	compressedHandler.ServeHTTP(rr, req)

	// Check if the content encoding header is not set (since client does not support gzip)
	assert.Empty(t, rr.Header().Get("Content-Encoding"))

	// Check if the response body is not gzipped
	body := rr.Body.String()
	assert.Equal(t, "Lorem ipsum dolor sit amet", body)
}

func TestGzipCompressorWithGzipRequestBody(t *testing.T) {
	// Create a new request with gzip compressed body
	bodyContent := "Lorem ipsum dolor sit amet"
	var requestBody bytes.Buffer
	gz := gzip.NewWriter(&requestBody)
	_, err := gz.Write([]byte(bodyContent))
	assert.NoError(t, err)
	assert.NoError(t, gz.Close())

	req, err := http.NewRequest("POST", "/", &requestBody)
	req.Header.Set("Content-Encoding", "gzip")
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, bodyContent, string(body))

		_, _ = w.Write([]byte("OK"))
	})

	compressedHandler := GzipCompressor(handler)
	compressedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGzipCompressorWithoutGzipRequestBody(t *testing.T) {
	// Create a new request without gzip compressed body
	req, err := http.NewRequest("POST", "/", strings.NewReader("Lorem ipsum dolor sit amet"))
	assert.NoError(t, err)

	// Do not set Content-Encoding header to simulate non-gzip request body

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		assert.Equal(t, "Lorem ipsum dolor sit amet", string(body))
		_, _ = w.Write([]byte("OK"))
	})

	compressedHandler := GzipCompressor(handler)
	compressedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
