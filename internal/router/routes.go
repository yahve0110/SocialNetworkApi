package myrouter

import (
	"social/internal/handlers"
	commentHandlers "social/internal/handlers/comment"
	followHandlers "social/internal/handlers/follows"
	groupHandlers "social/internal/handlers/group"
	groupPostCommentHandlers "social/internal/handlers/group/groupComments"
	groupChat "social/internal/handlers/group/groupMessages"
	groupPostHandlers "social/internal/handlers/group/groupPosts"
	groupInviteHandlers "social/internal/handlers/group/invitesAndRequests"
	messageHandlers "social/internal/handlers/messages"
	"social/internal/handlers/notifications"
	postHandler "social/internal/handlers/post"
	userHandlers "social/internal/handlers/user"
	userFeedHandler "social/internal/handlers/userFeed"
	"social/internal/middleware"

	groupEventHandlers "social/internal/handlers/group/groupEvents"
)

// DefineRoutes defines the routes and middleware
func DefineRoutes() *Router {
	router := NewRouter()

	// Register middleware for specific routes
	router.Handle("POST", "/register", handlers.UserRegister)
	router.Handle("POST", "/login", handlers.UserLogin)
	router.Handle("POST", "/logout", handlers.UserLogout)
	router.Handle("POST", "/privacyUpdate", middleware.LogMiddleware(userHandlers.UpdatePrivacy), middleware.AuthMiddleware)

	router.Handle("POST", "/addpost", middleware.LogMiddleware(postHandler.AddPost), middleware.AuthMiddleware)
	router.Handle("DELETE", "/deletePost", middleware.LogMiddleware(postHandler.DeletePostHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/getposts", middleware.LogMiddleware(postHandler.GetUserPosts), middleware.AuthMiddleware)
	router.Handle("POST", "/addPostLike", middleware.LogMiddleware(postHandler.AddPostLike), middleware.AuthMiddleware)

	router.Handle("POST", "/addcomment", middleware.LogMiddleware(commentHandlers.AddComment), middleware.AuthMiddleware)
	router.Handle("DELETE", "/deleteComment", middleware.LogMiddleware(commentHandlers.DeleteComment), middleware.AuthMiddleware)
	router.Handle("GET", "/comments", middleware.LogMiddleware(commentHandlers.GetCommentsForPost), middleware.AuthMiddleware)
	router.Handle("POST", "/addCommentLike", middleware.LogMiddleware(commentHandlers.AddCommentLike), middleware.AuthMiddleware)

	router.Handle("POST", "/isSessionValid", middleware.LogMiddleware(handlers.IsSessionValid))
	router.Handle("GET", "/getUserInfo", middleware.LogMiddleware(userHandlers.GetUserInfo), middleware.AuthMiddleware)
	router.Handle("GET", "/getAllUsers", middleware.LogMiddleware(handlers.GetAllUsers), middleware.AuthMiddleware)
	router.Handle("GET", "/getUserInfoById", middleware.LogMiddleware(userHandlers.GetUserInfoById), middleware.AuthMiddleware)

	router.Handle("POST", "/followUser", middleware.LogMiddleware(followHandlers.FollowUser), middleware.AuthMiddleware)
	router.Handle("POST", "/unfollowUser", middleware.LogMiddleware(followHandlers.UnfollowUserHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/getFollowers", middleware.LogMiddleware(followHandlers.GetFollowersHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/getFollowing", middleware.LogMiddleware(followHandlers.GetFollowingHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/getPendingFollowers", middleware.LogMiddleware(followHandlers.GetFollowersWithPendingStatusHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/acceptPendingFollowers", middleware.LogMiddleware(followHandlers.AcceptPendingFollowerHandler), middleware.AuthMiddleware)

	router.Handle("POST", "/createGroup", middleware.LogMiddleware(groupHandlers.CreateGroupHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/getallGroups", middleware.LogMiddleware(groupHandlers.GetAllGroupHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/getGroupById", middleware.LogMiddleware(groupHandlers.GetGroupByID), middleware.AuthMiddleware)
	router.Handle("GET", "/GetAllGroupMembers", middleware.LogMiddleware(groupHandlers.GetAllGroupMembers), middleware.AuthMiddleware)

	router.Handle("GET", "/getMyGroups", middleware.LogMiddleware(groupHandlers.GetMyGroups), middleware.AuthMiddleware)
	router.Handle("GET", "/getPendingRequestGroups", middleware.LogMiddleware(groupHandlers.GetRequestedGroups), middleware.AuthMiddleware)
	router.Handle("POST", "/inviteToGroup", middleware.LogMiddleware(groupInviteHandlers.SendGroupInvitationHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/checkGroupInvites", middleware.LogMiddleware(groupInviteHandlers.GetUserInvitationsHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/getAllUninvitedFollowers", middleware.LogMiddleware(groupInviteHandlers.GetUninvitedFollowersHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/acceptGroupInvite", middleware.LogMiddleware(groupInviteHandlers.AcceptGroupInvitationHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/leaveGroup", middleware.LogMiddleware(groupHandlers.LeaveGroupHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/decliceGroupInvite", middleware.LogMiddleware(groupInviteHandlers.RefuseGroupInvitationHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/sendGroupEnterRequest", middleware.LogMiddleware(groupInviteHandlers.SendGroupRequestHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/acceptGroupEnterRequest", middleware.LogMiddleware(groupInviteHandlers.AcceptGroupRequestHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/getAllGroupEnterRequests", middleware.LogMiddleware(groupInviteHandlers.GetAllGroupRequestsHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/createGroupPost", middleware.LogMiddleware(groupPostHandlers.CreateGroupPostHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/getGroupFeed", middleware.LogMiddleware(groupHandlers.GetGroupFeedHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/getAllGroupPosts", middleware.LogMiddleware(groupPostHandlers.GetAllGroupPostsHandler), middleware.AuthMiddleware)
	router.Handle("DELETE", "/deleteGroupPost", middleware.LogMiddleware(groupPostHandlers.DeleteGroupPostHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/addGroupPostComment", middleware.LogMiddleware(groupPostCommentHandlers.AddGroupPostCommentHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/addGroupPostLike", middleware.LogMiddleware(groupPostHandlers.LikeGroupPostHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/getGroupPostComments", middleware.LogMiddleware(groupPostCommentHandlers.GetGroupPostCommentsHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/addGroupCommentLike", middleware.LogMiddleware(groupPostCommentHandlers.LikeGroupCommentHandler), middleware.AuthMiddleware)
	router.Handle("DELETE", "/deleteGroupComment", middleware.LogMiddleware(groupPostCommentHandlers.DeleteGroupCommentHandler), middleware.AuthMiddleware)
	router.Handle("DELETE", "/deleteGroup", middleware.LogMiddleware(groupHandlers.DeleteGroupHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/createEvent", middleware.LogMiddleware(groupEventHandlers.CreateGroupEventHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/joinEvent", middleware.LogMiddleware(groupEventHandlers.JoinGroupEventHandler), middleware.AuthMiddleware)
	router.Handle("POST", "/declineEvent", middleware.LogMiddleware(groupEventHandlers.DeclineEventHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/getGroupEvents", middleware.LogMiddleware(groupEventHandlers.GetGroupEvents), middleware.AuthMiddleware)

	router.Handle("GET", "/getUserFeed", middleware.LogMiddleware(userFeedHandler.GetUserFeedHandler), middleware.AuthMiddleware)

	router.Handle("POST", "/openChat", middleware.LogMiddleware(messageHandlers.OpenChat), middleware.AuthMiddleware)
	router.Handle("GET", "/getChats", middleware.LogMiddleware(messageHandlers.GetUserChats), middleware.AuthMiddleware)
	router.Handle("GET", "/getChatHistory", middleware.LogMiddleware(messageHandlers.GetChatHistory), middleware.AuthMiddleware)
	router.Handle("GET", "/getChatHistory", middleware.LogMiddleware(messageHandlers.GetChatHistory), middleware.AuthMiddleware)

	router.Handle("POST", "/joinGroupChat", middleware.LogMiddleware(groupChat.JoinGroupChatHandler), middleware.AuthMiddleware)
	router.Handle("GET", "/getGroupChatHistory", middleware.LogMiddleware(groupChat.GetGroupChatHistory), middleware.AuthMiddleware)
	router.Handle("GET", "/getGroupChats", middleware.LogMiddleware(groupChat.GetGroupChats), middleware.AuthMiddleware)

	router.Handle("POST", "/sendNotification", middleware.LogMiddleware(notifications.ReceiveNotification), middleware.AuthMiddleware)
	router.Handle("GET", "/getNotifications", middleware.LogMiddleware(notifications.GetNotifications), middleware.AuthMiddleware)
	router.Handle("DELETE", "/deleteNotification", middleware.LogMiddleware(notifications.DeleteNotification), middleware.AuthMiddleware)


	return router
}
