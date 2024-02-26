package myrouter

import (
	"social/internal/handlers"
	pagesHandlers "social/internal/handlers/pages"
	"social/internal/middleware"
)

// DefineRoutes defines the routes and middleware
func DefineRoutes() *Router {
	router := NewRouter()



	// Register middleware for specific routes
	router.Handle("GET", "/", middleware.LogMiddleware(pagesHandlers.HandleHome), middleware.AuthMiddleware)
	router.Handle("GET", "/about", middleware.LogMiddleware(pagesHandlers.HandleAbout), middleware.AuthMiddleware)
	router.Handle("GET", "/post/:postid", middleware.LogMiddleware(handlers.HandlePost), middleware.AuthMiddleware)
	router.Handle("POST", "/register", handlers.UserRegister)
	router.Handle("POST", "/login", handlers.UserLogin)
	router.Handle("GET", "/logout", handlers.UserLogout)
	router.Handle("POST", "/addpost", middleware.LogMiddleware(handlers.AddPost), middleware.AuthMiddleware)
	router.Handle("GET", "/getposts", middleware.LogMiddleware(handlers.GetUserPosts), middleware.AuthMiddleware)
	router.Handle("POST", "/addcomment", middleware.LogMiddleware(handlers.AddComment), middleware.AuthMiddleware)
	router.Handle("GET", "/comments", middleware.LogMiddleware(handlers.GetCommentsForPost), middleware.AuthMiddleware)
	router.Handle("GET", "/isSessionValid", middleware.LogMiddleware(handlers.IsSessionValid))
	router.Handle("POST", "/addPostLike", middleware.LogMiddleware(handlers.AddPostLike),middleware.AuthMiddleware)
	router.Handle("POST", "/addCommentLike", middleware.LogMiddleware(handlers.AddCommentLike),middleware.AuthMiddleware)
	


	return router
}
