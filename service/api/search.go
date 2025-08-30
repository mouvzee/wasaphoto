package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
	"github.com/mouvzee/wasaphoto/service/api/reqcontext"
)

/*
searchUsers is the handler for the GET /profiles?search=query endpoint
It returns the users that match the given search query
*/
func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the search query from the request
	query_search := r.URL.Query().Get("search")

	validQuerySearch := regexp.MustCompile(`^[a-zA-Z0-9_-]{1,13}$`)
	if query_search == "" || !validQuerySearch.MatchString(query_search) {
		http.Error(w, "Bad Request: invalid search query", http.StatusBadRequest)
		return
	}

	userID := ctx.UserID

	dbUsers, err := rt.db.SearchUsers(userID, query_search)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error searching users")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	users := make([]User, len(dbUsers))
	for i, u := range dbUsers {
		var user User
		err := user.TakeUser(u)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error converting users")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		users[i] = user
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding users")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
