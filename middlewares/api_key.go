package middlewares

import (
	"net/http"
)

// APIKey middleware untuk validate API key
func APIKey(validAPIKey string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Get API key dari header
			apiKey := r.Header.Get("X-API-Key")

			// Validate
			if apiKey == "" {
				http.Error(w, "API key required", http.StatusUnauthorized)
				return
			}

			if apiKey != validAPIKey {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			// API key valid, continue
			next(w, r)
		}
	}
}
