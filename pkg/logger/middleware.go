package logger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type requestInfo struct {
	Method string `json:"method"`
	URI    string `json:"uri"`
}

type responseInfo struct {
	Status int `json:"status"`
	Size   int `json:"size"`
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseInfo *responseInfo
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseInfo.Size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseInfo.Status = statusCode
}

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
