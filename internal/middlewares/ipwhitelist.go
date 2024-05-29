package middlewares

import (
	"net"
	"net/http"
)

// IPWhitelist checks if the incoming request's IP address belongs to a trusted subnet
func IPWhitelist(trustedNet *net.IPNet) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hIP := r.Header.Get("X-Real-IP")
			if hIP == "" {
				http.Error(w, "X-Real-IP header is required", http.StatusBadRequest)
				return
			}

			// parse IP from header
			ip := net.ParseIP(hIP)
			if ip == nil {
				http.Error(w, "invalid X-Real-IP header value", http.StatusBadRequest)
				return
			}

			// compare IP with trusted subnet
			if !trustedNet.Contains(ip) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
