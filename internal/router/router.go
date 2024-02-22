package myrouter

import (
	"net/http"
	"strings"
)


// MiddlewareFunc defines the type for middleware functions
type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// Route represents information about a specific route
type Route struct {
	Path       string
	Method     string
	Handler    http.HandlerFunc
	Middleware []MiddlewareFunc
}

// Router represents a simple router
type Router struct {
	routes []Route
}

// NewRouter creates a new instance of the router
func NewRouter() *Router {
	return &Router{}
}

// Use registers middleware for the router
func (r *Router) Use(middleware ...MiddlewareFunc) *Router {
	for i := range r.routes {
		r.routes[i].Middleware = append(r.routes[i].Middleware, middleware...)
	}
	return r
}

// Handle registers a handler for the specified path and method
func (r *Router) Handle(method, path string, handler http.HandlerFunc, middleware ...MiddlewareFunc) *Router {
	route := Route{
		Path:       path,
		Method:     method,
		Handler:    handler,
		Middleware: middleware,
	}
	r.routes = append(r.routes, route)
	return r
}

// ServeHTTP handles HTTP requests
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	// Find the handler for the current path and method
	var handler http.HandlerFunc
	var middleware []MiddlewareFunc

	for _, route := range r.routes {
		if (route.Path == path || matchPath(route.Path, path)) && route.Method == method {
			handler = route.Handler
			middleware = route.Middleware
			break
		}
	}

	// Apply middleware functions in reverse order
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}

	// Call the request handler if found, otherwise return 404
	if handler != nil {
		handler(w, req)
	} else {
		http.NotFound(w, req)
	}
}

// matchPath checks if the request path matches the route
// Added to support dynamic parameters
func matchPath(route, path string) bool {
	routeParts := strings.Split(route, "/")
	pathParts := strings.Split(path, "/")

	if len(routeParts) != len(pathParts) {
		return false
	}

	for i := 0; i < len(routeParts); i++ {
		// If the route part is not a dynamic parameter and does not match, return false
		if !strings.HasPrefix(routeParts[i], ":") && routeParts[i] != pathParts[i] {
			return false
		}
	}

	return true
}


