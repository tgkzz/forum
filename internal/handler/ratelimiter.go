package handler

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	visits map[string]time.Time
	mutex  *sync.Mutex
	limit  int
	window time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		visits: make(map[string]time.Time),
		mutex:  &sync.Mutex{},
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		rl.mutex.Lock()
		lastVisit, exists := rl.visits[ip]
		if !exists || time.Since(lastVisit) > rl.window {
			rl.visits[ip] = time.Now()
			rl.mutex.Unlock()
			next.ServeHTTP(w, r)
		} else {
			rl.mutex.Unlock()
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
	}
}
