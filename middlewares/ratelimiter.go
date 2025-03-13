package middlewares

import (
	"log"
	"net"
	"net/http"
	helper "superapps/helpers"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*rate.Limiter
	lastUsed map[string]time.Time // 🆕 Menyimpan waktu terakhir limiter digunakan
	r        rate.Limit
	b        int
}

// NewRateLimiter membuat instance RateLimiter baru
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	rl := &RateLimiter{
		clients:  make(map[string]*rate.Limiter),
		lastUsed: make(map[string]time.Time),
		r:        r,
		b:        b,
	}

	go rl.cleanupExpiredLimiters()

	return rl
}

// getClientIP mendapatkan IP asli dari request, memperhitungkan reverse proxy
func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// getLimiter mengembalikan rate limiter untuk IP tertentu
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.clients[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.clients[ip] = limiter
	}
	rl.lastUsed[ip] = time.Now() // 🆕 Update waktu terakhir digunakan

	return limiter
}

// LimitMiddleware menerapkan rate limiting pada endpoint tertentu
func (rl *RateLimiter) LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		limiter := rl.getLimiter(ip)
		if !limiter.Allow() {
			helper.Logger("error", "In Server: Too many requests. Please try again later")
			helper.Response(w, 429, true, "Too many requests. Please try again later", map[string]any{}) // 🆕 429 bukan 500
			return
		}

		next.ServeHTTP(w, r)
	})
}

// cleanupExpiredLimiters menghapus limiter yang tidak aktif selama 10 menit
func (rl *RateLimiter) cleanupExpiredLimiters() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Running cleanupExpiredLimiters...") // 🆕 Logging untuk debug

		rl.mu.Lock()
		now := time.Now()
		for ip, last := range rl.lastUsed {
			if now.Sub(last) > 10*time.Minute {
				log.Printf("Deleting limiter for IP: %s\n", ip) // 🆕 Logging IP yang dihapus
				delete(rl.clients, ip)
				delete(rl.lastUsed, ip)
			}
		}
		rl.mu.Unlock()
	}
}
