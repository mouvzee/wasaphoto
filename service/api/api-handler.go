package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply, true))
	rt.router.POST("/login", rt.wrap(rt.doLogin, false))
	rt.router.PUT("/user/:userID/username", rt.wrap(rt.setMyUserName, true))
	rt.router.POST("/photo/:photoID/user/:userID", rt.wrap(rt.UploadPhoto, true))
	rt.router.DELETE("/user/:UserID/post/:PostID", rt.wrap(rt.deletePhoto, true))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
