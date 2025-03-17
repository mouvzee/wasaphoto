package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mouvzee/wasaphoto/service/api/reqcontext"
)

/*
* ChangeCaption is the handler for the PUT /profiles/:profileUserID/posts/:PhotoID endpoint
* It allows a user to change the caption of his own post
 */
func (rt *_router) changeCaption(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the PhotoID from the URL
	PhotoID, err := strconv.Atoi(ps.ByName("PhotoID"))
	if err != nil {
		http.Error(w, "Invalid PhotoID", http.StatusBadRequest)
		return
	}

	// Get the profileUserID from the URL
	profileUserID, err := strconv.Atoi(ps.ByName("profileUserID"))
	if err != nil {
		http.Error(w, "Invalid profileUserID", http.StatusBadRequest)
		return
	}

	// Get the userID from the context
	userID := ctx.UserID

	// Check if the user is the owner of the post
	if userID != profileUserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Read the request body, save the new caption in the variable NewCaptionBody
	type NewCaptionBody struct {
		NewCaption string `json:"caption"`
	}

	var newCaptionBody NewCaptionBody
	err = json.NewDecoder(r.Body).Decode(&newCaptionBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validCaption := regexp.MustCompile(`^[^\/\\]{0,64}$`)
	if !validCaption.MatchString(newCaptionBody.NewCaption) {
		http.Error(w, "Invalid caption, /, \\, new line and tab are not allowed. Max lenght is 64 characters", http.StatusBadRequest)
		return
	}

	// Update the caption in the database
	err = rt.db.UpdateCaption(userID, PhotoID, newCaptionBody.NewCaption)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error while updating the caption")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
