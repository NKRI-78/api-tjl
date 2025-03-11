package middlewares

import (
	"net/http"
	helper "superapps/helpers"
	"sync"
	"time"
)

var rateLimit = struct {
	requests map[string][]time.Time
	mutex    sync.Mutex
}{requests: make(map[string][]time.Time)}

// RateLimitingMiddleware allows up to 2 requests per minute per IP
func RateLimitingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}

		rateLimit.mutex.Lock()
		timestamps := rateLimit.requests[ip]

		// Remove expired requests older than 1 minute
		now := time.Now()
		validTimestamps := []time.Time{}
		for _, t := range timestamps {
			if now.Sub(t) < time.Minute {
				validTimestamps = append(validTimestamps, t)
			}
		}

		// Check if request limit is exceeded
		if len(validTimestamps) >= 2 {
			helper.Logger("error", "In Server: Too many requests, slow down!")
			helper.Response(w, 500, true, "Too many requests, slow down!", map[string]any{})
			rateLimit.mutex.Unlock()
			return
		}

		// Add the current request timestamp
		validTimestamps = append(validTimestamps, now)
		rateLimit.requests[ip] = validTimestamps
		rateLimit.mutex.Unlock()

		next.ServeHTTP(w, r)
	})
}
