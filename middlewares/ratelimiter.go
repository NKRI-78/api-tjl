package middlewares

import (
	"net"
	"net/http"
	helper "superapps/helpers"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*rate.Limiter
	r       rate.Limit
	b       int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*rate.Limiter),
		r:       r,
		b:       b,
	}

	go rl.cleanupExpiredLimiters()

	return rl
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	limiter, exists := rl.clients[ip]
	rl.mu.Unlock()

	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)

		rl.mu.Lock()
		rl.clients[ip] = limiter
		rl.mu.Unlock()
	}

	return limiter
}

func (rl *RateLimiter) LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Invalid IP address", http.StatusInternalServerError)
			return
		}

		limiter := rl.getLimiter(ip)
		if !limiter.Allow() {
			helper.Logger("error", "In Server: Too many requests. Please try again later")
			helper.Response(w, 500, true, "Too many requests. Please try again later", map[string]any{})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) cleanupExpiredLimiters() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, limiter := range rl.clients {
			if limiter.Burst() == 0 {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}
