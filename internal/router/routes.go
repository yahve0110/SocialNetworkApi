package myrouter

import (
	"social/internal/handlers"
	commentHandlers "social/internal/handlers/comment"
	followHandlers "social/internal/handlers/follows"
	groupHandlers "social/internal/handlers/group"
	postHandler "social/internal/handlers/post"
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

	router.Handle("POST", "/createGroup", middleware.LogMiddleware(groupHandlers.CreateGroupHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/getallGroups", middleware.LogMiddleware(groupHandlers.GetAllGroupHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/inviteToGroup", middleware.LogMiddleware(groupHandlers.SendGroupInvitationHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/checkGroupInvites", middleware.LogMiddleware(groupHandlers.GetUserInvitationsHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/acceptGroupInvite", middleware.LogMiddleware(groupHandlers.AcceptGroupInvitationHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/leaveGroup", middleware.LogMiddleware(groupHandlers.LeaveGroupHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/decliceGroupInvite", middleware.LogMiddleware(groupHandlers.RefuseGroupInvitationHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/sendGroupEnterRequest", middleware.LogMiddleware(groupHandlers.SendGroupRequestHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/acceptGroupEnterRequest", middleware.LogMiddleware(groupHandlers.AcceptGroupRequestHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/getAllGroupEnterRequests", middleware.LogMiddleware(groupHandlers.GetAllGroupRequestsHandler), middleware.AuthMiddleware)

	return router
}
  