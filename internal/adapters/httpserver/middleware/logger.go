package middleware

import (
	"net/http"
	"time"
)

type MiddlewareLogger interface {
	Infow(msg string, kv ...any)
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

type Logger struct {
	log MiddlewareLogger
}

func NewLogger(log MiddlewareLogger) *Logger {
	return &Logger{
		log: log,
	}
}

func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: rw,
			responseData:   responseData,
		}
		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		contentType := r.Header.Get("Content-Type")
		acceptEncoding := r.Header.Get("Accept-Encoding")
		contentEncoding := r.Header.Get("Content-Encoding")

		l.log.Infow(
			"recieved request",
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", int64(duration),
			"size", responseData.size,
			"content-type", contentType,
			"content-encoding", contentEncoding,
			"accept-encoding", acceptEncoding,
		)
	})
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
