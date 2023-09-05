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

func TestUpdate(t *testing.T) {
	m := storage.New()

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
				statusCode:   http.StatusBadRequest,
				responseText: "Bad method\n",
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
				statusCode:   http.StatusBadRequest,
				responseText: "Can't parse URL\n",
			},
			method:      http.MethodPost,
			urlPath:     "/update/",
			contentType: "text/plain",
		},
		{
			name: "test #5 - Request without metrics name",
			want: want{
				statusCode:   http.StatusNotFound,
				responseText: "Metrics name not set\n",
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
				statusCode:   http.StatusBadRequest,
				responseText: "Metrics value not set or invalid\n",
			},
			method:      http.MethodPost,
			urlPath:     "/update/counter/someMetric/",
			contentType: "text/plain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.urlPath, nil)
			request.Header.Add("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(Update(m))
			h(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			// assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			responseText, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.responseText, string(responseText))
		})
	}
}
