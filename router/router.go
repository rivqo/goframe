package router

import (
	"net/http"
	"strings"
	"sync"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)
type MiddlewareFunc func(http.Handler) http.Handler

type Router struct {
	routes      map[string]map[string]http.Handler
	middlewares []MiddlewareFunc
	notFound    http.Handler
	staticDirs  map[string]string
	staticFiles map[string]string
	mu          sync.RWMutex // Protects concurrent access to routes and static maps
}

type RouteGroup struct {
	prefix      string
	router      *Router
	middlewares []MiddlewareFunc
}

func New() *Router {
	return &Router{
		routes:      make(map[string]map[string]http.Handler),
		middlewares: make([]MiddlewareFunc, 0),
		notFound:    http.NotFoundHandler(),
		staticDirs:  make(map[string]string),
		staticFiles: make(map[string]string),
	}
}

func (r *Router) Use(mw MiddlewareFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.middlewares = append(r.middlewares, mw)
}

func (r *Router) Group(prefix string) *RouteGroup {
	return &RouteGroup{
		prefix:      strings.TrimSuffix(prefix, "/"),
		router:      r,
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (g *RouteGroup) Use(mw MiddlewareFunc) {
	g.middlewares = append(g.middlewares, mw)
}

func (r *Router) Handle(method, path string, handler http.Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	method = strings.ToUpper(method)
	path = "/" + strings.Trim(path, "/")
	
	if _, ok := r.routes[method]; !ok {
		r.routes[method] = make(map[string]http.Handler)
	}
	r.routes[method][path] = handler
}

func (r *Router) register(method, path string, handler HandlerFunc) {
	var wrapped http.Handler = http.HandlerFunc(handler)
	
	r.mu.RLock()
	middlewares := make([]MiddlewareFunc, len(r.middlewares))
	copy(middlewares, r.middlewares)
	r.mu.RUnlock()
	
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}
	
	r.Handle(method, path, wrapped)
}

func (r *Router) Get(path string, handler HandlerFunc)    { r.register("GET", path, handler) }
func (r *Router) Post(path string, handler HandlerFunc)   { r.register("POST", path, handler) }
func (r *Router) Put(path string, handler HandlerFunc)    { r.register("PUT", path, handler) }
func (r *Router) Delete(path string, handler HandlerFunc) { r.register("DELETE", path, handler) }

func (g *RouteGroup) register(method, path string, handler HandlerFunc) {
	fullPath := g.prefix + "/" + strings.Trim(path, "/")
	var wrapped http.Handler = http.HandlerFunc(handler)

	// Combine global and group middlewares
	g.router.mu.RLock()
	globalMiddlewares := make([]MiddlewareFunc, len(g.router.middlewares))
	copy(globalMiddlewares, g.router.middlewares)
	g.router.mu.RUnlock()
	
	allMiddlewares := append(globalMiddlewares, g.middlewares...)
	
	for i := len(allMiddlewares) - 1; i >= 0; i-- {
		wrapped = allMiddlewares[i](wrapped)
	}

	g.router.Handle(method, fullPath, wrapped)
}

func (g *RouteGroup) Get(path string, handler HandlerFunc)    { g.register("GET", path, handler) }
func (g *RouteGroup) Post(path string, handler HandlerFunc)   { g.register("POST", path, handler) }
func (g *RouteGroup) Put(path string, handler HandlerFunc)    { g.register("PUT", path, handler) }
func (g *RouteGroup) Delete(path string, handler HandlerFunc) { g.register("DELETE", path, handler) }

func (r *Router) Static(prefix, dir string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	prefix = "/" + strings.Trim(prefix, "/") + "/"
	r.staticDirs[prefix] = dir
}

func (r *Router) StaticFile(path, file string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	path = "/" + strings.Trim(path, "/")
	r.staticFiles[path] = file
}

func (r *Router) NotFound(handler HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notFound = http.HandlerFunc(handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r == nil {
		http.Error(w, "Internal Server Error: Router not initialized", http.StatusInternalServerError)
		return
	}

	path := req.URL.Path
	method := req.Method

	// Handle static files first
	if method == "GET" {
		r.mu.RLock()
		// Check static files
		if filePath, ok := r.staticFiles[path]; ok {
			r.mu.RUnlock()
			http.ServeFile(w, req, filePath)
			return
		}

		// Check static directories
		for prefix, dir := range r.staticDirs {
			if strings.HasPrefix(path, prefix) {
				r.mu.RUnlock()
				fs := http.StripPrefix(prefix, http.FileServer(http.Dir(dir)))
				fs.ServeHTTP(w, req)
				return
			}
		}
		r.mu.RUnlock()
	}

	// Handle registered routes
	r.mu.RLock()
	if methodRoutes, ok := r.routes[method]; ok {
		if handler, ok := methodRoutes[path]; ok {
			r.mu.RUnlock()
			handler.ServeHTTP(w, req)
			return
		}
	}
	r.mu.RUnlock()

	// Fallback to not found handler
	r.mu.RLock()
	notFound := r.notFound
	r.mu.RUnlock()
	
	notFound.ServeHTTP(w, req)
}