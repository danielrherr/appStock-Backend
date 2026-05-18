package middleware

import (
	"net/http"
	"net/url"
	"os"
	"strings"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && isAllowedOrigin(r, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			if origin != "" && !isAllowedOrigin(r, origin) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAllowedOrigin(r *http.Request, origin string) bool {
	parsedOrigin, err := url.Parse(origin)
	if err != nil || parsedOrigin.Host == "" {
		return false
	}

	// Local development
	if isLocalDevHost(parsedOrigin.Hostname()) {
		return true
	}

	// Same origin
	if sameOrigin(r, parsedOrigin) {
		return true
	}

	// Allowed domains from env (comma-separated)
	allowedDomains := os.Getenv("ALLOWED_DOMAINS")
	if allowedDomains != "" {
		for _, domain := range strings.Split(allowedDomains, ",") {
			if strings.TrimSpace(domain) == parsedOrigin.Hostname() {
				return true
			}
		}
	}

	return false
}

func isLocalDevHost(host string) bool {
	host = strings.ToLower(host)
	return host == "localhost" || host == "127.0.0.1"
}

func sameOrigin(r *http.Request, origin *url.URL) bool {
	requestHost := r.Host
	if requestHost == "" {
		return false
	}

	requestScheme := "http"
	if r.TLS != nil {
		requestScheme = "https"
	}
	if forwardedProto := r.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		requestScheme = forwardedProto
	}

	return strings.EqualFold(origin.Scheme, requestScheme) && strings.EqualFold(origin.Host, requestHost)
}
