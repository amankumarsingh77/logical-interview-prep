package main

import (
	"log"
	"sync"
	"time"
)

type rateLimiter struct {
	requests        map[string][]time.Time
	window          time.Duration
	limit           int
	cleanUpInterval time.Duration
	sync.Mutex
}

type RateLimiter interface {
	IsAllowed(ipAddress string) bool
	CleanInactiveUserLog()
}

func NewRateLimiter(window time.Duration, limit int, cleanUpInterval time.Duration) RateLimiter {
	limiter := &rateLimiter{
		requests:        make(map[string][]time.Time),
		window:          window,
		limit:           limit,
		cleanUpInterval: cleanUpInterval,
	}
	go limiter.CleanInactiveUserLog()
	return limiter
}
func (r *rateLimiter) IsAllowed(ipAddress string) bool {
	r.Lock()
	defer r.Unlock()
	timeStamps, ok := r.requests[ipAddress]
	if !ok {
		r.requests[ipAddress] = append(r.requests[ipAddress], time.Now())
		return true
	}
	now := time.Now()
	cutOff := now.Add(-r.window)
	var recentRequests []time.Time
	for _, t := range timeStamps {
		if t.After(cutOff) {
			recentRequests = append(recentRequests, t)
		}
	}
	if len(recentRequests) >= r.limit {
		r.requests[ipAddress] = recentRequests
		return false
	}
	recentRequests = append(recentRequests, now)
	r.requests[ipAddress] = recentRequests
	return true
}

func (r *rateLimiter) CleanInactiveUserLog() {
	ticker := time.NewTicker(r.cleanUpInterval)
	for range ticker.C {
		log.Println("cleaning inactive user logs")
		r.Lock()
		cutoff := time.Now().Add(-r.window)
		for ip, timestamps := range r.requests {
			if len(timestamps) > 0 && timestamps[len(timestamps)-1].Before(cutoff) {
				delete(r.requests, ip)
			}
		}
		r.Unlock()
	}
}
