package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mouvzee/wasaphoto/service/api/reqcontext"
)

/*
deletePhoto is the hanlder for the DELETE /users/:profileUserID/posts/:PhotoID endpoint.
It deletes a photo from the database and returns a 200 OK if the photo is deleted.
*/
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the user ID from the URL
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the post ID from the URL
	photoID, err := strconv.Atoi(ps.ByName("photoID"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}
	// Print all parameters in ps
	for i := 0; i < len(ps); i++ {
		ctx.Logger.Infof("Param %s: %s", ps[i].Key, ps[i].Value)
	}

	// token, err := strconv.Atoi(ps.ByName("Token"))
	// if err != nil {
	// 	http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// Check if the user is authorized
	if userID != ctx.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Delete the photo from the database
	if err := rt.db.Delete_Photo(photoID); err != nil {
		ctx.Logger.WithError(err).Error("Error while deleting the post")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
}
