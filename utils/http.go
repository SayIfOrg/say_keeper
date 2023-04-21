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
