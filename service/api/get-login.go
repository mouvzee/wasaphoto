package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mouvzee/wasaphoto/service/api/reqcontext"
)

func (rt *_router) getLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, "Error in the request.", http.StatusBadRequest)
		return
	}

	if !user.isValid() {
		http.Error(w, "Not Valid Username.", http.StatusBadRequest)
		return
	}

	b, err := rt.db.CheckIfExist(user.Username)

	if err != nil {
		ctx.Logger.WithError(err).Error("Don't really know if the user exist")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !b {
		_, err := rt.create_user(user)
		if err != nil {
			ctx.Logger.WithError(err).Error("Cannot create the user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		dbUser, err := rt.db.Get_User_ByUsername(user.Username)
		if err != nil {
			ctx.Logger.WithError(err).Error("Cannot find the user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = user.takingUser(dbUser)
		if err != nil {
			ctx.Logger.WithError(err).Error("Cannot create the user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
