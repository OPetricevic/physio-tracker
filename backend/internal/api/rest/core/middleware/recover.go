package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// RecoverMiddleware catches panics, logs stack + request info, and returns 500.
// This keeps the server alive and provides actionable context on crashes.
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v\nmethod=%s path=%s time=%s\nstack=%s",
					rec, r.Method, r.URL.Path, time.Now().UTC().Format(time.RFC3339), debug.Stack())
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
