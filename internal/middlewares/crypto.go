package middlewares

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/SerjRamone/metrius/pkg/logger"
)

// Crypto decrypts request body
func Crypto(privKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error("request body reading error", zap.Error(err))
				return
			}
			r.Body.Close()

			// parse pem-encoded key
			block, _ := pem.Decode(privKey)
			key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error("parsing key error", zap.Error(err))
				return
			}

			msgLen := len(bodyBytes)
			// chunk size
			step := key.PublicKey.Size()
			var decryptedBytes []byte

			// decrypt by chunks
			for start := 0; start < msgLen; start += step {
				finish := start + step
				if finish > msgLen {
					finish = msgLen
				}

				decryptedBlockBytes, err := key.Decrypt(nil, bodyBytes[start:finish], &rsa.OAEPOptions{Hash: crypto.SHA512})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Error("decrypt error", zap.Error(err))
					return
				}

				decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
			}

			// set unread body
			r.Body = io.NopCloser(bytes.NewBuffer(decryptedBytes))

			next.ServeHTTP(w, r)
		})
	}
}
