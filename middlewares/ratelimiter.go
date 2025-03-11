package middlewares

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// Create a struct to hold each client's rate limiter
type Client struct {
	limiter *rate.Limiter
}

// In-memory storage for clients
var clients = make(map[string]*Client)
var mu sync.Mutex

// Get a client's rate limiter or create one if it doesn't exist
func GetClientLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	// If the client already exists, return the existing limiter
	if client, exists := clients[ip]; exists {
		return client.limiter
	}

	// Create a new limiter with 2 requests per minute
	limiter := rate.NewLimiter(2, 1)
	clients[ip] = &Client{limiter: limiter}
	return limiter
}

func RateLimitingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := GetClientLimiter(ip)

		// Check if the request is allowed
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
