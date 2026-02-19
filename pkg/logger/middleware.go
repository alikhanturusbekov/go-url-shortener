package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// requestInfo contains basic request metadata
type requestInfo struct {
	Method string `json:"method"`
	URI    string `json:"uri"`
}

// responseInfo contains response metadata
type responseInfo struct {
	Status int `json:"status"`
	Size   int `json:"size"`
}

// loggingResponseWriter wraps http.ResponseWriter to capture response data
type loggingResponseWriter struct {
	http.ResponseWriter
	responseInfo *responseInfo
}

// Write writes the response body and tracks its size
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseInfo.Size += size
	return size, err
}

// WriteHeader writes the HTTP status code and records it
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseInfo.Status = statusCode
}

// RequestLogger provides HTTP middleware for structured request logging
func RequestLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			responseInfo := &responseInfo{}
			lw := loggingResponseWriter{
				ResponseWriter: w,
				responseInfo:   responseInfo,
			}

			next.ServeHTTP(&lw, r)
			duration := time.Since(start)

			Log.Info("request completed",
				zap.Any("request", requestInfo{
					Method: r.Method,
					URI:    r.RequestURI,
				}),
				zap.Any("response", responseInfo),
				zap.String("duration", duration.String()),
			)
		})
	}
}
