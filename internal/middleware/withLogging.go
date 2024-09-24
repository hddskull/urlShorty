package middleware

import (
	"github.com/hddskull/urlShorty/internal/utils"
	"net/http"
	"time"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	responseData := &responseData{
		status: 0,
		size:   0,
	}
	lw := &loggingResponseWriter{
		ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
		responseData:   responseData,
	}
	return lw
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := newLoggingResponseWriter(w)

		h.ServeHTTP(lw, r)

		duration := time.Since(start)

		utils.SugaredLogger.Infow(
			"handler",
			"uri", r.RequestURI,
			"method", r.Method,
			"duration", duration,
			"status", lw.responseData.status,
			"size", lw.responseData.size,
		)
	}
	return http.HandlerFunc(logFn)
}
