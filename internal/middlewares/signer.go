package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/SerjRamone/metrius/pkg/logger"
)

// Signer adds header with hash key
func Signer(hashKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerHash := r.Header.Get("HashSHA256")
			srw := &signResponseWriter{
				ResponseWriter: w,
				HashKey:        hashKey,
				Body:           new(bytes.Buffer),
			}
			if headerHash != "" {
				logger.Info("request with hash header", zap.String("hash", headerHash))
				// read request body
				bodyBytes, _ := io.ReadAll(r.Body)
				r.Body.Close()
				// set unread body
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// calculate body hash with key
				b64Hash := CalcHash(bodyBytes, []byte(hashKey))
				logger.Info("hashes", zap.String("calculated", string(b64Hash)), zap.String("header", headerHash))

				// compare hashes calculated and from request header
				if string(b64Hash) != headerHash {
					http.Error(w, "invalid sign", http.StatusBadRequest)
					return
				}
			} else {
				logger.Info("request without 'HashSHA256' header")
			}

			next.ServeHTTP(srw, r)
		})
	}
}

// CalcHash calculate sign for bytes slice with key
// returns encoded in base64 hash string
func CalcHash(body, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(body)
	calculatedHash := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(calculatedHash)
}

type signResponseWriter struct {
	http.ResponseWriter
	Body    *bytes.Buffer
	HashKey string
}

func (rw *signResponseWriter) Write(b []byte) (int, error) {
	rw.Body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (rw *signResponseWriter) WriteHeader(code int) {
	logger.Info("response status", zap.Int("status", code))
	if code == 200 {
		b64Hash := CalcHash(rw.Body.Bytes(), []byte(rw.HashKey))
		rw.ResponseWriter.Header().Add("HashSHA256", b64Hash)
	}
	rw.ResponseWriter.WriteHeader(code)
}
