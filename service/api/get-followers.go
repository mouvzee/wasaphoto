package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mouvzee/wasaphoto/service/api/reqcontext"
)

/*
GetMyFollowers is the handler for the GET /users/:profileUserID/followers endpoint
It return the list of the followers of the user with the given profileUserID
*/
func (rt *_router) listFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the profileUserID and targetUserID from the URL
	profileUserID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
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

	// Get the followers
	dbFollowers, err := rt.db.GetFollowers(profileUserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error while getting the followers")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert the database followers to the API followers
	followers := make([]User, len(dbFollowers))
	for i, dbUser := range dbFollowers {
		var user User
		err := user.TakeUser(dbUser)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error while converting the user")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		followers[i] = user
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(followers); err != nil {
		ctx.Logger.WithError(err).Error("Error while encoding the followers")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
