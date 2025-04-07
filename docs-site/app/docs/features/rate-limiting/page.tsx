import { DocPagination } from "@/components/doc-pagination"

export default function RateLimitingPage() {
  return (
    <div className="space-y-6">
      <h1>Rate Limiting</h1>
      <p>
        GoFrame provides a rate limiting middleware to protect your application from excessive requests. Rate limiting
        is essential for preventing abuse, reducing server load, and ensuring fair usage of your API.
      </p>

      <h2>Configuration</h2>
      <p>
        Rate limiting settings are configured in the <code>config.yaml</code> file:
      </p>

      <pre>
        <code>{`rateLimit:
  requests: 100
  period: 1m`}</code>
      </pre>

      <p>
        This configuration allows 100 requests per minute per IP address. The <code>period</code> can be specified using
        Go's duration format (e.g., <code>1s</code>, <code>1m</code>, <code>1h</code>).
      </p>

      <h2>Rate Limiting Middleware</h2>
      <p>
        The rate limiting middleware is implemented in the <code>middleware</code> package:
      </p>

      <pre>
        <code>{`// middleware/middleware.go
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
}`}</code>
      </pre>

      <h2>Applying Rate Limiting</h2>
      <p>You can apply rate limiting to your entire application or specific routes:</p>

      <pre>
        <code>{`// main.go
// Apply global middleware
r.Use(middleware.Logger())
r.Use(middleware.RateLimit(cfg.RateLimit.Requests, cfg.RateLimit.Period))
r.Use(middleware.Recover())`}</code>
      </pre>

      <p>Or apply it to specific route groups:</p>

      <pre>
        <code>{`// routes/api.go
// Create API group with rate limiting
api := r.Group("/api")
api.Use(middleware.RateLimit(50, 1*time.Minute)) // 50 requests per minute
api.Use(authProvider.Middleware())`}</code>
      </pre>

      <h2>Custom Rate Limiting</h2>
      <p>
        You can create custom rate limiting strategies by modifying the middleware. For example, you might want to rate
        limit based on user ID for authenticated users, or apply different limits to different endpoints.
      </p>

      <pre>
        <code>{`// middleware/custom_rate_limit.go
// UserRateLimit is a middleware that limits requests based on user ID
func UserRateLimit(requests int, period time.Duration) func(http.Handler) http.Handler {
    type user struct {
        count    int
        lastSeen time.Time
    }

    var (
        users = make(map[uint]*user)
        mu    sync.Mutex
    )

    // Clean up old users periodically
    go func() {
        for {
            time.Sleep(period)
            mu.Lock()
            for id, user := range users {
                if time.Since(user.lastSeen) > period {
                    delete(users, id)
                }
            }
            mu.Unlock()
        }
    }()

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Get user from context
            currentUser := auth.GetUser(r.Context())
            if currentUser == nil {
                // Fall back to IP-based rate limiting for unauthenticated users
                next.ServeHTTP(w, r)
                return
            }

            userID := currentUser.ID
            
            mu.Lock()
            if _, found := users[userID]; !found {
                users[userID] = &user{count: 0, lastSeen: time.Now()}
            }
            
            u := users[userID]
            u.lastSeen = time.Now()
            
            if u.count >= requests {
                mu.Unlock()
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            u.count++
            mu.Unlock()
            
            next.ServeHTTP(w, r)
            
            // Decrement the count after the period
            time.AfterFunc(period, func() {
                mu.Lock()
                defer mu.Unlock()
                if u, found := users[userID]; found {
                    u.count--
                    if u.count < 0 {
                        u.count = 0
                    }
                }
            })
        })
    }
}`}</code>
      </pre>

      <h2>Rate Limiting Headers</h2>
      <p>
        You can enhance the rate limiting middleware to include headers that inform clients about their rate limit
        status:
      </p>

      <pre>
        <code>{`// Enhanced rate limiting middleware with headers
func RateLimitWithHeaders(requests int, period time.Duration) func(http.Handler) http.Handler {
    // ... existing code ...

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := r.RemoteAddr
            
            mu.Lock()
            if _, found := clients[ip]; !found {
                clients[ip] = &client{count: 0, lastSeen: time.Now()}
            }
            
            client := clients[ip]
            client.lastSeen = time.Now()
            
            // Add rate limit headers
            w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requests))
            w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", requests - client.count))
            w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", client.lastSeen.Add(period).Unix()))
            
            if client.count >= requests {
                mu.Unlock()
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            client.count++
            mu.Unlock()
            
            next.ServeHTTP(w, r)
            
            // ... existing code ...
        })
    }
}`}</code>
      </pre>

      <DocPagination
        prev={{
          title: "Resources",
          href: "/docs/features/resources",
        }}
        next={{
          title: "CLI Commands",
          href: "/docs/features/cli-commands",
        }}
      />
    </div>
  )
}

