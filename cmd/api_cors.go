package main

import (
	"net/http"
	"net/url"
	"os"
	"strings"
)

func apiCORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := configuredFrontendOrigin()
		if origin != "" && r.Header.Get("Origin") == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func configuredFrontendOrigin() string {
	frontendURL, err := url.Parse(strings.TrimSpace(os.Getenv("FRONTEND_URL")))
	if err != nil || frontendURL.Host == "" || (frontendURL.Scheme != "https" && frontendURL.Scheme != "http") {
		return ""
	}
	return frontendURL.Scheme + "://" + frontendURL.Host
}
