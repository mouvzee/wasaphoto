package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mouvzee/wasaphoto/service/api/reqcontext"
)

/*
getMyBans is the handler for GET /users/:profileUserID/bans endpoint
It return the list of the users that the profileUserID has banned
*/

func (rt *_router) getBannedList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the profileUserID and targetUserID from the URL
	profileUserID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userID := ctx.UserID

	// Check if the user is authorized
	if profileUserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Get the bans from the database
	dbBans, err := rt.db.GetBan(profileUserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting bans")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert the bans to the User struct
	bans := make([]User, len(dbBans))

	for i, dbBan := range dbBans {
		var user User
		err := user.TakeUser(dbBan)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error while converting the user")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		bans[i] = user
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bans); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding bans")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
