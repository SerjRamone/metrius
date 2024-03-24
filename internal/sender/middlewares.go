package sender

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"

	"github.com/SerjRamone/metrius/internal/middlewares"
	"github.com/SerjRamone/metrius/pkg/logger"
	"go.uber.org/zap"
)

// Adds HashSHA256 header to request headers
func hasher(hashKey string) middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return internalRoundTripper(func(req *http.Request) (*http.Response, error) {
			var buf bytes.Buffer
			// copy data from request body
			_, err := io.Copy(&buf, req.Body)
			if err != nil {
				logger.Error("request body copy error", zap.Error(err))
				return nil, err
			}
			// get bytes from buffer
			b := buf.Bytes()
			// set new body
			req.Body = io.NopCloser(bytes.NewReader(b))
			// calculate hash string
			b64Hash := middlewares.CalcHash(b, []byte(hashKey))
			req.Header.Set("HashSHA256", b64Hash)
			logger.Info("calculated body hash", zap.String("hash", b64Hash))

			return rt.RoundTrip(req)
		})
	}
}

// encrypt requst body
func crypto(pubKey []byte) middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return internalRoundTripper(func(req *http.Request) (*http.Response, error) {
			// parse pem-encoded key
			block, _ := pem.Decode(pubKey)
			if block == nil {
				logger.Error("parsing public key error", zap.String("reason", "invalid PEM format"))
				return nil, fmt.Errorf("invalid PEM format")
			}

			// convert key to *rsa.PublicKey
			pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				logger.Error("parsing public key error", zap.Error(err))
				return nil, err
			}

			// read request body
			var buf bytes.Buffer
			_, err = io.Copy(&buf, req.Body)
			if err != nil {
				logger.Error("request body copy error", zap.Error(err))
				return nil, err
			}

			// hash func
			hash := sha512.New()

			bodyBytes := buf.Bytes()
			msgLen := len(bodyBytes)

			// chunk size
			step := pub.Size() - 2*hash.Size() - 2
			var encryptedBytes []byte

			// encode body by chunks
			for start := 0; start < msgLen; start += step {
				finish := start + step
				if finish > msgLen {
					finish = msgLen
				}

				encryptedBlockBytes, err := rsa.EncryptOAEP(hash, rand.Reader, pub, bodyBytes[start:finish], nil)
				if err != nil {
					logger.Error("encrypting chunk error", zap.Error(err))
					return nil, err
				}

				encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
			}

			// set encoded request body
			req.Body = io.NopCloser(bytes.NewReader(encryptedBytes))
			req.ContentLength = int64(len(encryptedBytes))

			return rt.RoundTrip(req)
		})
	}
}
