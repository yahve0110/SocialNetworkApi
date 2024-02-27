package myrouter

import (
	"social/internal/handlers"
	"social/internal/handlers/comment"
	followHandlers "social/internal/handlers/follows"
	"social/internal/handlers/post"
	"social/internal/middleware"
)

// DefineRoutes defines the routes and middleware
func DefineRoutes() *Router {
	router := NewRouter()

	// Register middleware for specific routes

	router.Handle("POST", "/register", handlers.UserRegister)
	router.Handle("POST", "/login", handlers.UserLogin)
	router.Handle("GET", "/logout", handlers.UserLogout)

	router.Handle("POST", "/addpost", middleware.LogMiddleware(postHandler.AddPost), middleware.AuthMiddleware)
	router.Handle("GET", "/getposts", middleware.LogMiddleware(postHandler.GetUserPosts), middleware.AuthMiddleware)
	router.Handle("POST", "/addPostLike", middleware.LogMiddleware(postHandler.AddPostLike), middleware.AuthMiddleware)

	router.Handle("POST", "/addcomment", middleware.LogMiddleware(commentHandlers.AddComment), middleware.AuthMiddleware)
	router.Handle("GET", "/comments", middleware.LogMiddleware(commentHandlers.GetCommentsForPost), middleware.AuthMiddleware)
	router.Handle("POST", "/addCommentLike", middleware.LogMiddleware(commentHandlers.AddCommentLike), middleware.AuthMiddleware)


	router.Handle("GET", "/isSessionValid", middleware.LogMiddleware(handlers.IsSessionValid))
	router.Handle("GET", "/getAllUsers", middleware.LogMiddleware(handlers.GetAllUsers), middleware.AuthMiddleware)

	router.Handle("POST", "/followUser", middleware.LogMiddleware(followHandlers.FollowUser), middleware.AuthMiddleware)
	router.Handle("POST", "/unfollowUser", middleware.LogMiddleware(followHandlers.UnfollowUserHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/getFollowers", middleware.LogMiddleware(followHandlers.GetFollowersHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/getFollowing", middleware.LogMiddleware(followHandlers.GetFollowingHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/getPendingFollowers", middleware.LogMiddleware(followHandlers.GetFollowersWithPendingStatusHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/acceptPendingFollowers", middleware.LogMiddleware(followHandlers.AcceptPendingFollowerHandler), middleware.AuthMiddleware)


	return router
}
