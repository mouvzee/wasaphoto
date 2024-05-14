package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mouvzee/wasaphoto/service/api/reqcontext"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the userID from the URL
	UserID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user is authorized
	if UserID != ctx.UserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	//check if user is valid
	if !u.isValid() {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	// Change the username, if the new username is already taken, the request will fail
	if err := rt.db.ChangeUsername(UserID, u.Username); err != nil {
		http.Error(w, "Username already taken. Username must be unique", http.StatusBadRequest)
		return
	}

	// Return a success message
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "plain/text")
	if err := json.NewEncoder(w).Encode("Username changed"); err != nil {
		ctx.Logger.WithError(err).Error("Response not encoted")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
