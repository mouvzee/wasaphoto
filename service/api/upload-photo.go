package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mouvzee/wasaphoto/service/api/reqcontext"
)

func (rt *_router) UploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	UserID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	userID := ctx.UserID

	//check if the user is authorized
	if UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Analyse the multipart form to put a limit in the memory
	err = r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		http.Error(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		return
	}

	caption := r.FormValue("caption")

	// Get the file from the form
	file, _, err := r.FormFile("")
	if err != nil {
		http.Error(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		return
	}

	// Read the file
	data, err := io.ReadAll(file)
	if err != nil {
		ctx.Logger.WithError(err).Error("error parse file")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// check if the user is uploading an image/jpeg
	fileType := http.DetectContentType(data)
	if fileType != "image/jpeg" {
		http.Error(w, "Bad Request wrong file type", http.StatusBadRequest)
		return
	}

	defer func() { err = file.Close() }()

	//get the user by the ID 
	dbUser, err := rt.db.GetUserbyID(userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting user")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//analyse the user from the database to create the new post
	var u User
	err = u.TakeUser(dbUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("Not possible to parse the user")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var new = Photo{
		User: u,
		Caption: caption,
	}

	//check if the caption is valid
	if !new.isValid(){
		http.Error(w, "Invalid Caption", http.StatusBadRequest)
		return
	}

	//Saving photo on the database
	dbPhoto := new.SavingPhoto()

	dbNewPost, err := rt.db.CreatePhoto(dbPhoto, data)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error creating post")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse the new post from the database package to the Photo struct in the api package
	err = new.TakingPhoto(dbNewPost)
	if err != nil {
		ctx.Logger.WithError(err).Error("error parsing photo")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return the new post
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(new); err != nil {
		ctx.Logger.WithError(err).Error("Error while encoding the post")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
