package myrouter

import (
	"fmt"
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
func (r *Router) Use(middleware ...MiddlewareFunc) *Router {
    fmt.Println("Entering Use method")

    for _, m := range middleware {
        fmt.Printf("Middleware function: %p\n", m)
    }

    for i := range r.routes {
        fmt.Printf("Before appending middleware to route %d: %v\n", i, r.routes[i].Middleware)
        r.routes[i].Middleware = append(r.routes[i].Middleware, middleware...)
        fmt.Printf("After appending middleware to route %d: %v\n", i, r.routes[i].Middleware)
    }

    fmt.Printf("Exiting Use method. Router address: %p\n", r)
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

func (r *Router) PrintDetails() {
	fmt.Printf("Number of routes: %d\n", len(r.routes))
	for i, route := range r.routes {
		fmt.Printf("Route %d:\n", i+1)
		fmt.Printf("  Path: %s\n", route.Path)
		fmt.Printf("  Method: %s\n", route.Method)
		fmt.Printf("  Middleware: %v\n", route.Middleware)
	}
}
