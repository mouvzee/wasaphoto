package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mattn/go-sqlite3"
	"github.com/mouvzee/wasaphoto/service/api/reqcontext"
)

/*
likePhoto is the handler for the POST /users/:profileUserID/posts/:PhotoID/like endpoint.
It adds a like to the post with the given ID.
*/
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	profileUserID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	userID := ctx.UserID

	isBanned, err := rt.db.IsBanned(profileUserID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user is banned")
		http.Error(w, "Internal Server Error "+err.Error(), http.StatusInternalServerError)
		return
	}

	if isBanned {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	PhotoID, err := strconv.Atoi(ps.ByName("PhotoID"))
	if err != nil {
		http.Error(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.db.CreateLike(userID, PhotoID)
	if err != nil {
		var sqlite3Err sqlite3.Error
		if errors.As(err, &sqlite3Err) {
			if sqlite3Err.Code == sqlite3.ErrConstraint && sqlite3Err.ExtendedCode == 1555 {
				http.Error(w, "Bad Request like already added", http.StatusBadRequest)
				return
			}
		}
		ctx.Logger.WithError(err).Error("error liking post")
		http.Error(w, "Internal Server Error "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
