package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

// New constructor for gzipResponseWriter to initialize the gzip writer
func NewGzipResponseWriter(w http.ResponseWriter) *gzipResponseWriter {
	gz := gzip.NewWriter(w)
	return &gzipResponseWriter{ResponseWriter: w, writer: gz}
}

// Write wraps the gzip.Writer's Write method to compress data before writing it to the underlying ResponseWriter
func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

// WriteHeader wraps the ResponseWriter's WriteHeader method
func (g *gzipResponseWriter) WriteHeader(statusCode int) {
	g.ResponseWriter.WriteHeader(statusCode)
}

// Header wraps the ResponseWriter's Header method to allow manipulation of response headers
func (g *gzipResponseWriter) Header() http.Header {
	return g.ResponseWriter.Header()
}

// Close should be called to close the gzip writer and flush the compressed data to the client
func (g *gzipResponseWriter) Close() error {
	return g.writer.Close()
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gz := NewGzipResponseWriter(w)
			defer gz.Close()
			w.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(gz, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
