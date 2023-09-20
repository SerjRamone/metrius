package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SerjRamone/metrius/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// @todo mayby content-type check
func testRequest(t *testing.T, ts *httptest.Server, method,
	path string,
) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	m := storage.New()
	_ = m.SetCounter("foo", 1)

	ts := httptest.NewServer(Router(m))
	defer ts.Close()

	type want struct {
		statusCode   int
		responseText string
	}

	tests := []struct {
		name        string
		want        want
		method      string
		urlPath     string
		contentType string
	}{
		{
			name: "test #1 - 200 OK",
			want: want{
				statusCode:   http.StatusOK,
				responseText: "OK",
			},
			method:      http.MethodPost,
			urlPath:     "/update/counter/someMetric/527",
			contentType: "text/plain",
		},
		{
			name: "test #2 - Bad method ",
			want: want{
				statusCode:   http.StatusMethodNotAllowed,
				responseText: "",
			},
			method:      http.MethodGet,
			urlPath:     "/update/counter/someMetric/527",
			contentType: "text/plain",
		},
		// {
		// 	name: "test #3 - Bad Content-Type",
		// 	want: want{
		// 		statusCode:   http.StatusBadRequest,
		// 		responseText: "Bad content-type\n",
		// 	},
		// 	method:      http.MethodPost,
		// 	urlPath:     "/update/counter/someMetric/527",
		// 	contentType: "application/json",
		// },
		{
			name: "test #4 - Bad URL",
			want: want{
				statusCode:   http.StatusNotFound,
				responseText: "404 page not found\n",
			},
			method:      http.MethodPost,
			urlPath:     "/update/",
			contentType: "text/plain",
		},
		{
			name: "test #5 - Request without metrics name",
			want: want{
				statusCode:   http.StatusNotFound,
				responseText: "404 page not found\n",
			},
			method:      http.MethodPost,
			urlPath:     "/update/counter/",
			contentType: "text/plain",
		},
		{
			name: "test #6 - Request with unknown type",
			want: want{
				statusCode:   http.StatusBadRequest,
				responseText: "Metrics type not set or unknown\n",
			},
			method:      http.MethodPost,
			urlPath:     "/update/foo/someMetric/527",
			contentType: "text/plain",
		},
		{
			name: "test #7 - Request without value",
			want: want{
				statusCode:   http.StatusNotFound,
				responseText: "404 page not found\n",
			},
			method:      http.MethodPost,
			urlPath:     "/update/counter/someMetric/",
			contentType: "text/plain",
		},
		{
			name: "test #8 - Request with bad value",
			want: want{
				statusCode:   http.StatusBadRequest,
				responseText: "Invalid metrics value\n",
			},
			method:      http.MethodPost,
			urlPath:     "/update/counter/someMetric/bad",
			contentType: "text/plain",
		},
		{
			name: "test #9 - Get metrics value - valid request",
			want: want{
				statusCode:   http.StatusOK,
				responseText: "1",
			},
			method:      http.MethodGet,
			urlPath:     "/value/counter/foo",
			contentType: "text/plain",
		},
		{
			name: "test #10 - Get metrics value - request without metrics name",
			want: want{
				statusCode:   http.StatusNotFound,
				responseText: "404 page not found\n",
			},
			method:      http.MethodGet,
			urlPath:     "/value/counter/",
			contentType: "text/plain",
		},
		{
			name: "test #11 - Get metrics value - unknown metrics",
			want: want{
				statusCode:   http.StatusNotFound,
				responseText: "not found\n",
			},
			method:      http.MethodGet,
			urlPath:     "/value/counter/unknown",
			contentType: "text/plain",
		},
		{
			name: "test #12 - Get metrics value - invalid metrics type",
			want: want{
				statusCode:   http.StatusBadRequest,
				responseText: "unknown type\n",
			},
			method:      http.MethodGet,
			urlPath:     "/value/unknown/foo",
			contentType: "text/plain",
		},
		{
			name: "test #13 - Main page - OK",
			want: want{
				statusCode:   http.StatusOK,
				responseText: "",
			},
			method:      http.MethodGet,
			urlPath:     "/",
			contentType: "text/plain",
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			resp, respText := testRequest(t, ts, v.method, v.urlPath)
			resp.Body.Close()
			assert.Equal(t, v.want.statusCode, resp.StatusCode)
			if v.want.responseText != "" {
				assert.Equal(t, v.want.responseText, respText)
			}
		})
	}
}
