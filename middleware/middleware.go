package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"sync"
	"time"
)

// Logger is a middleware that logs request details
func Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Call the next handler
			next.ServeHTTP(w, r)
			
			// Log the request
			log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
		})
	}
}

// Recover is a middleware that recovers from panics
func Recover() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Panic: %v\n%s", err, debug.Stack())
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimit is a middleware that limits the number of requests in a period
func RateLimit(requests int, period time.Duration) func(http.Handler) http.Handler {
	type client struct {
		count    int
		lastSeen time.Time
	}
	
	var (
		clients = make(map[string]*client)
		mu      sync.Mutex
	)
	
	// Clean up old clients periodically
	go func() {
		for {
			time.Sleep(period)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > period {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			
			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{count: 0, lastSeen: time.Now()}
			}
			
			client := clients[ip]
			client.lastSeen = time.Now()
			
			if client.count >= requests {
				mu.Unlock()
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			
			client.count++
			mu.Unlock()
			
			next.ServeHTTP(w, r)
			
			// Decrement the count after the period
			time.AfterFunc(period, func() {
				mu.Lock()
				defer mu.Unlock()
				if c, found := clients[ip]; found {
					c.count--
					if c.count < 0 {
						c.count = 0
					}
				}
			})
		})
	}
}

