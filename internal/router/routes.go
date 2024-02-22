package myrouter

import (
	"fmt"
	"net/http"
	"social/internal/handlers/pages"
	"social/internal/handlers"
)

// LogMiddleware is a simple middleware for logging
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Logging middleware: ", r.Method, r.URL.Path)
		next(w, r)
	}
}

// DefineRoutes defines the routes and middleware
func DefineRoutes() *Router {
	router := NewRouter()

	// Register middleware for specific routes
	router.Handle("GET", "/", pagesHandlers.HandleHome, LogMiddleware)
	router.Handle("GET", "/about", pagesHandlers.HandleAbout, LogMiddleware)
	router.Handle("GET", "/post/:postid", handlers.HandlePost, LogMiddleware)
	router.Handle("POST", "/register", handlers.UserRegister)


	return router
}
