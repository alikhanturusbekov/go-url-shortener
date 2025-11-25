package compress

import (
	"net/http"
	"strings"
)

func GzipCompressor() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ow := w

			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			if supportsGzip {
				cw := newCompressWriter(w)
				ow = cw
				defer cw.Close()
			}

			contentEncoding := r.Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")
			if sendsGzip {
				cr, err := newCompressReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				r.Body = cr
				defer cr.Close()
			}

			next.ServeHTTP(ow, r)
		})
	}
}

func isCompressible(contentType string) bool {
	return strings.HasPrefix(contentType, "text/plain") ||
		strings.HasPrefix(contentType, "application/json")
}
