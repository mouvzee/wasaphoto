package database

var AddPhoto = `INSERT Photo(PhotoID, UserID, Caption) WITH VALUES (?,?,?)`

func (db *appdbimpl) uploadPhoto(y Photo, ImageData []byte) error {
	_, err := db.c.Exec(AddPhoto, y.PhotoID, y.User.UserID, y.Caption)
	return err
}

//gestire l'errore al richiamo della funzione

//func(y Photo, ImageData []byte) (y, error){
//profileUserID := y.User.UserID
//path := users.GetPostPhotoPath(profileUserID, y.PostID)

// Save the image
//err = os.WriteFile(path, data, 0666)
//if err != nil {
//	return p, err
//}

//}
