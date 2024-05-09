package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mouvzee/wasaphoto/service/api/reqcontext"
)

/*
deletePhoto is the hanlder for the DELETE /users/:profileUserID/posts/:postID endpoint.
It deletes a photo from the database and returns a 200 OK if the photo is deleted.
*/
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the user ID from the URL
	UserID, err := strconv.Atoi(ps.ByName("UserID"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the post ID from the URL
	PhotoID, err := strconv.Atoi(ps.ByName("postID"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}
	//correggi guarda wrap prova con un token (te lo devi inventare tu)
	userID := UserID

	// Check if the user is authorized
	if UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Delete the photo from the database
	if err := rt.db.Delete_Photo(UserID, PhotoID); err != nil {
		ctx.Logger.WithError(err).Error("Error while deleting the post")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
}
