package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mouvzee/wasaphoto/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

/*
GetLikes is the handler for the GET /users/:profileUserID/posts/:PhotoID/likes endpoint
It returns the likes of the post with the given PhotoID
*/

func (rt *_router) getLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the profileUserID and PhotoID from the URL
	profileUserID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	PhotoID, err := strconv.Atoi(ps.ByName("PhotoID"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	userID := ctx.UserID

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

	dbLikes, err := rt.db.GetLike(PhotoID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting likes")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	likes := make([]User, len(dbLikes))

	for i, dbLike := range dbLikes {
		var user User
		err := user.TakeUser(dbLike)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error while converting the user")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		likes[i] = user
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(likes); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding likes")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
