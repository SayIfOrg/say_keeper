package utils

import (
	"net/http"
	"net/url"
	"strings"
)

//CheckAllowedOrigin returns a function that will check a http request to have the allowed origin header
func CheckAllowedOrigin(allowedOrigins []string, sameOrigin bool) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		origin := r.Header["Origin"]
		u, err := url.Parse(origin[0])
		if err != nil {
			return false
		}

		// Allow SameOrigin
		if sameOrigin && strings.EqualFold(u.Host, r.Host) {
			return true
		}

		for _, allowedOrigin := range allowedOrigins {
			if strings.EqualFold(u.Host, allowedOrigin) == true {
				return true
			}
		}
		return false
	}
}

//CorsMiddleware expects allowedOrigin in format like "https://some.me"
func CorsMiddleware(next http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the CORS headers
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If this is a preflight request, send the headers and return
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next middleware/handler in the chain
		next.ServeHTTP(w, r)
	})
}
