package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		w, _ := gzip.NewWriterLevel(io.Discard, gzip.HuffmanOnly)
		return w
	},
}

var gzipReaderPool = sync.Pool{
	New: func() interface{} {
		return new(gzip.Reader)
	},
}

func GzipCompressor() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Request Decompression
			if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				gr := gzipReaderPool.Get().(*gzip.Reader)
				if err := gr.Reset(r.Body); err != nil {
					http.Error(w, "invalid gzip body", http.StatusBadRequest)
					return
				}
				defer func() {
					gr.Close()
					gzipReaderPool.Put(gr)
				}()
				r.Body = &readCloser{
					Reader: gr,
					Closer: r.Body,
				}
			}

			// Response Compression
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			gz := gzipWriterPool.Get().(*gzip.Writer)
			defer gzipWriterPool.Put(gz)

			gz.Reset(w)
			defer gz.Close()

			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Del("Content-Length")

			gzw := &gzipResponseWriter{
				ResponseWriter: w,
				writer:         gz,
			}

			next.ServeHTTP(gzw, r)
		})
	}
}

type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.writer.Write(b)
}

type readCloser struct {
	io.Reader
	io.Closer
}
