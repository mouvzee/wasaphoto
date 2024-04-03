package database

var DelPhoto = "DELETE Photo(PhotoID, UserID, Caption) WITH VALUE(?,?,?)"

func (db *appdbimpl) Delete_Photo(y Photo, URL string) error {
	
}
