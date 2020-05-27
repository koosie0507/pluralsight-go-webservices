package middleware

import (
	"log"
	"net/http"
	"time"
)

// Log is a middleware that logs some stats to standard output
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		handler.ServeHTTP(w, r)
		log.Printf("%s %s finished in %s\n", r.Method, r.URL.Path, time.Since(now))
	})
}
