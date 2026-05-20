package middleware

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		clientIP := r.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = r.RemoteAddr
		}

		if lrw.statusCode >= 400 {
			log.Printf("HTTP_ERROR method=%s path=%s status=%d duration=%s ip=%s ua=%q", r.Method, r.URL.Path, lrw.statusCode, duration, clientIP, r.UserAgent())
			return
		}

		log.Printf("HTTP method=%s path=%s status=%d duration=%s ip=%s", r.Method, r.URL.Path, lrw.statusCode, duration, clientIP)
	})
}
