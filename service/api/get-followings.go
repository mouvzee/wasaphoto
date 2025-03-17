package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mouvzee/wasaphoto/service/api/methods"

	"github.com/julienschmidt/httprouter"
	"github.com/mouvzee/wasaphoto/service/api/reqcontext"
)

/*
getMyFollowings is the handler for the GET /users/:profileUserID/followings endpoint
It returns the followings of the user with the given profileUserID
*/
func (rt *_router) getMyFollowings(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the profileUserID and targetUserID from the URL
	profileUserID, err := strconv.Atoi(ps.ByName("profileUserID"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	limit, offset, err := methods.GetLimitAndOffset(r.URL.Query())
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the followings
	dbFollowings, err := rt.db.GetFollowings(profileUserID, offset, limit)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error while getting the followings")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert the database followings to the API followings
	followings := make([]User, len(dbFollowings))
	for i, dbUser := range dbFollowings {
		var user User
		err := user.TakeUser(dbUser)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error while converting the user")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		followings[i] = user
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(followings); err != nil {
		ctx.Logger.WithError(err).Error("Error while encoding the response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
