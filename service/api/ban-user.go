package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mouvzee/wasaphoto/service/api/reqcontext"
)

/*
banUser is the handler for the POST /users/{profileUserID}/bans/{targetUserID} endpoint.
It bans a user.
A banned user cannot see anything about the user who banned him.
*/
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the profileUserID and targetUserID from the URL
	profileUserID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	targetUserID, err := strconv.Atoi(ps.ByName("targetUserID"))
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

	// Check if the user is trying to ban himself
	if targetUserID == userID {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteFollow(profileUserID, targetUserID) // Rimuovi "io seguo lui"
	if err != nil {
		ctx.Logger.WithError(err).Warning("Error removing follow relationship during ban")
		// Non fermare il processo per questo errore
	}

	err = rt.db.DeleteFollow(targetUserID, profileUserID) // Rimuovi "lui segue me"
	if err != nil {
		ctx.Logger.WithError(err).Warning("Error removing follow relationship during ban")
	}

	// Create the ban in the database
	err = rt.db.CreateBan(profileUserID, targetUserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error creating ban")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
