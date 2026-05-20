package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const maxLoggedBodyBytes = 2048

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
		requestPath := r.URL.Path
		if r.URL.RawQuery != "" {
			requestPath += "?" + r.URL.RawQuery
		}

		bodyForLog := captureRequestBody(r)
		authHeader := redactAuthorization(r.Header.Get("Authorization"))

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		clientIP := r.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = r.RemoteAddr
		}

		if lrw.statusCode >= 400 {
			log.Printf("HTTP_ERROR method=%s path=%s status=%d duration=%s ip=%s auth=%q ua=%q body=%q", r.Method, requestPath, lrw.statusCode, duration, clientIP, authHeader, r.UserAgent(), bodyForLog)
			return
		}

		log.Printf("HTTP method=%s path=%s status=%d duration=%s ip=%s", r.Method, requestPath, lrw.statusCode, duration, clientIP)
	})
}

func captureRequestBody(r *http.Request) string {
	if r.Body == nil || r.Method == http.MethodGet || r.Method == http.MethodHead {
		return ""
	}

	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if strings.HasPrefix(contentType, "multipart/form-data") {
		return "<multipart-body>"
	}

	if r.ContentLength < 0 || r.ContentLength > maxLoggedBodyBytes {
		return "<body-skipped>"
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, maxLoggedBodyBytes+1))
	if err != nil {
		return "<unreadable-body>"
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return sanitizeBody(body)
}

func sanitizeBody(body []byte) string {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return ""
	}

	var payload map[string]any
	if err := json.Unmarshal(trimmed, &payload); err == nil {
		redactSensitiveFields(payload)
		sanitized, marshalErr := json.Marshal(payload)
		if marshalErr == nil {
			return string(sanitized)
		}
	}

	return string(trimmed)
}

func redactSensitiveFields(payload map[string]any) {
	for key, value := range payload {
		switch strings.ToLower(key) {
		case "password", "token", "access_token", "refresh_token", "authorization":
			payload[key] = "<redacted>"
			continue
		}

		nested, ok := value.(map[string]any)
		if ok {
			redactSensitiveFields(nested)
		}
	}
}

func redactAuthorization(value string) string {
	if value == "" {
		return ""
	}

	parts := strings.SplitN(value, " ", 2)
	if len(parts) == 2 {
		return parts[0] + " <redacted>"
	}

	return "<redacted>"
}
