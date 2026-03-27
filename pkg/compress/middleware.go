package compress

import (
	"compress/gzip"
	"github.com/alikhanturusbekov/go-url-shortener/pkg/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"

	"github.com/alikhanturusbekov/go-url-shortener/pkg/pool"
)

// PooledGzipWriter wraps gzip.Writer to implement Resetter.
type PooledGzipWriter struct {
	W *gzip.Writer
}

// Reset resets the writer to a clean state
func (p *PooledGzipWriter) Reset() {
	p.W.Reset(io.Discard)
}

// PooledGzipReader wraps gzip.Reader to implement Resetter.
type PooledGzipReader struct {
	R *gzip.Reader
}

// Reset closes the reader
func (p *PooledGzipReader) Reset() {
	_ = p.R.Close()
}

var gzipWriterPool = pool.New(func() *PooledGzipWriter {
	w, _ := gzip.NewWriterLevel(io.Discard, gzip.HuffmanOnly)
	return &PooledGzipWriter{W: w}
})

var gzipReaderPool = pool.New(func() *PooledGzipReader {
	return &PooledGzipReader{R: new(gzip.Reader)}
})

// GzipCompressor provides HTTP middleware for gzip compression and decompression
func GzipCompressor() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Request Decompression
			if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				gr := gzipReaderPool.Get()
				if err := gr.R.Reset(r.Body); err != nil {
					http.Error(w, "invalid gzip body", http.StatusBadRequest)
					return
				}
				defer func() {
					if err := gr.R.Close(); err != nil {
						logger.Log.Warn("failed to close gzip reader", zap.Error(err))
					}
					gzipReaderPool.Put(gr)
				}()
				r.Body = &readCloser{
					Reader: gr.R,
					Closer: r.Body,
				}
			}

			// Response Compression
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			gz := gzipWriterPool.Get()
			gz.W.Reset(w)
			defer func() {
				if err := gz.W.Close(); err != nil {
					logger.Log.Warn("failed to close gzip writer", zap.Error(err))
				}
				gzipWriterPool.Put(gz)
			}()

			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Del("Content-Length")

			gzw := &gzipResponseWriter{
				ResponseWriter: w,
				writer:         gz.W,
			}

			next.ServeHTTP(gzw, r)
		})
	}
}

// gzipResponseWriter wraps http.ResponseWriter to write compressed data
type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

// Write writes compressed response data
func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.writer.Write(b)
}

// readCloser combines a Reader and Closer
type readCloser struct {
	io.Reader
	io.Closer
}
