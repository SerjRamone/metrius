package middlewares

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignerMiddleware(t *testing.T) {
	testKey := "testHashKey"
	bodyBytes := []byte("test body")
	validHash := CalcHash(bodyBytes, []byte(testKey))
	testCases := []struct {
		name               string
		requestHash        string
		requestBody        []byte
		expectedStatusCode int
	}{
		{
			name:               "Test#1. Valid hash",
			requestBody:        bodyBytes,
			requestHash:        validHash,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Test#2. Invalid hash",
			requestBody:        bodyBytes,
			requestHash:        "invalid",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			middleware := Signer(testKey)

			req, err := http.NewRequest("POST", "/test", bytes.NewBuffer(tc.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("HashSHA256", tc.requestHash)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("testing")
			})
			middleware(handler).ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatusCode, rr.Code, "status code doesn't match")
			//
			// if tc.expectedStatusCode == http.StatusOK {
			// 	expectedHash := CalcHash(tc.requestBody, []byte(testKey))
			// 	assert.Equal(t, expectedHash, rr.Header().Get("HashSHA256"), "hash in response doesn't match")
			// }
		})
	}
}
