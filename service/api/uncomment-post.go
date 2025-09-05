package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mouvzee/wasaphoto/service/api/reqcontext"
)

/*
uncommentPhoto is the handler for the DELETE /profiles/:profileUserID/posts/:PhotoID/comments/:commentID endpoint
It deletes the comment with the given commentID from the post with the given PhotoID.
The comment is deleted only if the user with the given profileUserID is the owner of the comment.
The comment is deleted from the database.
*/

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the profileUserID and PhotoID from the URL
	_, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	photoID, err := strconv.Atoi(ps.ByName("photoID"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Get commentID from query
	commentID, err := strconv.Atoi(ps.ByName("commentID"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	userID := ctx.UserID

	// Delete the comment
	err = rt.db.DeleteComment(commentID, photoID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error uncommenting photo")
		// Se il database restituisce errore, significa nessun permesso
		if err.Error() == "comment not found or permission denied" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return a 200 OK status code
	w.WriteHeader(http.StatusOK)
}
