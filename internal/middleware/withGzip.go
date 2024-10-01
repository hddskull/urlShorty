package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/hddskull/urlShorty/internal/utils"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func WithGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isGzipEncoded := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		isJSON := strings.Contains(r.Header.Get("Content-Type"), "application/json")
		isHTML := strings.Contains(r.Header.Get("Content-Type"), "text/html")

		if !isGzipEncoded {
			utils.SugaredLogger.Infow("is not compressed", "isGzipEncoded", isGzipEncoded)
			next.ServeHTTP(w, r)
			return
		}
		if !(isHTML || isJSON) {
			utils.SugaredLogger.Infow("is not html or json", "isHTML", isHTML, "isJSON", isJSON)
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
