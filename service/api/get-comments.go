package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mouvzee/wasaphoto/service/api/reqcontext"
)

/*
GetComments is the handler for the GET /users/:profileUserID/posts/:PhotoID/comments endpoint
It returns the comments of the post with the given PhotoID
*/
func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the profileUserID from the URL
	profileUserID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get the PhotoID from the URL
	photoID, err := strconv.Atoi(ps.ByName("photoID"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userID := ctx.UserID

	// Check if the user is banned
	if profileUserID != userID {
		isBanned, err := rt.db.IsBanned(profileUserID, userID)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error checking if user is banned")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if isBanned {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	// Get the comments
	dbComments, err := rt.db.GetComments(photoID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting comments")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert the comments to the API struct
	comments := make([]Comment, len(dbComments))
	for i, dbComment := range dbComments {
		var comment Comment
		if err = comment.takingComment(dbComment); err != nil {
			ctx.Logger.WithError(err).Error("Error converting comment")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		comments[i] = comment
	}

	// Write the comments
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(comments); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding comments")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
