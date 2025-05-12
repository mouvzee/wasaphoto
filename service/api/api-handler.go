package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// -------LOGIN AND REGISTER--------//
	rt.router.POST("/login", rt.wrap(rt.doLogin, false))

	// ----------DELETE USER-------------//
	rt.router.DELETE("/profiles/:userID", rt.wrap(rt.deleteUser, true))

	// --------GET USER PROFILE----------//
	rt.router.GET("/profiles/:userID", rt.wrap(rt.getUserProfile, true))

	// --------CHANGE USERNAME---------//
	rt.router.PUT("/profiles/:userID/username", rt.wrap(rt.setMyUserName, true))

	// --------CHANGE PROFILE PIC-------//
	rt.router.PUT("/profiles/:userID/profile-picture", rt.wrap(rt.setMyProfilePic, true))

	// --------DELETE PROFILE PIC---------//
	rt.router.PUT("/profiles/:userID/reset-profile-picture", rt.wrap(rt.resetMyProfilePic, true))

	// ----------FOLLOW USER-----------//
	rt.router.PUT("/profiles/:userID/followings/:targetUserID", rt.wrap(rt.followUser, true))

	// --------UNFOLLOW USER-----------//
	rt.router.DELETE("/profiles/:userID/followings/:targetUserID", rt.wrap(rt.unfollowUser, true))

	// ----------GET FOLLOWINGS--------//
	rt.router.GET("/profiles/:userID/followings", rt.wrap(rt.getMyFollowings, true))

	// ----------GET FOLLOWERS---------//
	rt.router.GET("/profiles/:userID/followers", rt.wrap(rt.getMyFollowers, true))

	// ----------BAN USER--------------//
	rt.router.PUT("/profiles/:userID/bans/:targetUserID", rt.wrap(rt.banUser, true))

	// --------UNBAN USER--------------//
	rt.router.DELETE("/profiles/:userID/bans/:targetUserID", rt.wrap(rt.unbanUser, true))

	// -----------GET BANS------------//
	rt.router.GET("/profiles/:userID/bans", rt.wrap(rt.getMyBans, true))

	// ----------UPLOAD POST----------//
	rt.router.POST("/profiles/:userID/posts", rt.wrap(rt.uploadPhoto, true))

	// ----------GET POSTS------------//
	rt.router.GET("/profiles/:userID/posts", rt.wrap(rt.getPosts, true))

	// ----------CHANGE CAPTION-------------//
	rt.router.PUT("/profiles/:userID/posts/:photoID/caption", rt.wrap(rt.changeCaption, true))

	// ----------DELETE POST-----------//
	rt.router.DELETE("/profiles/:userID/posts/:photoID", rt.wrap(rt.deletePhoto, true))

	// ----------LIKE POST-------------//
	rt.router.PUT("/profiles/:userID/posts/:photoID/likes/:userID", rt.wrap(rt.likePhoto, true))

	// --------UNLIKE POST-------------//
	rt.router.DELETE("/profiles/:userID/posts/:photoID/likes/:userID", rt.wrap(rt.unlikePhoto, true))

	// ---------GET LIKES--------------//
	rt.router.GET("/profiles/:userID/posts/:photoID/likes", rt.wrap(rt.getLikes, true))

	// ----------COMMENT POST----------//
	rt.router.POST("/profiles/:userID/posts/:photoID/comments", rt.wrap(rt.commentPhoto, true))

	// ----------UNCOMMENT POST--------//
	rt.router.DELETE("/profiles/:userID/posts/:photoID/comments/:commentID", rt.wrap(rt.uncommentPhoto, true))

	// ---------GET COMMENTS-----------//
	rt.router.GET("/profiles/:userID/posts/:photoID/comments", rt.wrap(rt.getComments, true))

	// ---------GET FEED--------------//
	rt.router.GET("/profiles/:userID/feed", rt.wrap(rt.getMyStream, true))

	// -----------SEARCH------------//
	rt.router.GET("/profiles", rt.wrap(rt.searchUsers, true))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
