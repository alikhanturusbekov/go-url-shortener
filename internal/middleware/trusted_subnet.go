package middleware

import (
	"net"
	"net/http"
)

// TrustedSubnet returns a middleware that restricts access by client IP.
func TrustedSubnet(trustedSubnet string) func(http.Handler) http.Handler {
	// Parse CIDR once when middleware is created
	var subnet *net.IPNet
	if trustedSubnet != "" {
		_, parsedSubnet, err := net.ParseCIDR(trustedSubnet)
		if err == nil {
			subnet = parsedSubnet
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If subnet is not configured or invalid, deny all access
			if subnet == nil {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			realIP := r.Header.Get("X-Real-IP")
			if realIP == "" {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			ip := net.ParseIP(realIP)
			if ip == nil {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			if !subnet.Contains(ip) {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
