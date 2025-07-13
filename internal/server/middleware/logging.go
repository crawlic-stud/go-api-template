package middleware

import (
	"net/http"
	"time"
	"validation-api/internal/util/helper"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	errorMsg   string
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode >= 400 && lrw.statusCode < 500 {
		lrw.errorMsg = string(b)
	}
	return lrw.ResponseWriter.Write(b)
}

// NewLoggingMiddleware creates new logging middleware
func NewLoggingMiddleware(helper *helper.ServerHelper) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lrw := &loggingResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			start := time.Now()
			next.ServeHTTP(lrw, r)

			helper.LogRequest(lrw.statusCode, r, start, lrw.errorMsg)
		})
	}
}
