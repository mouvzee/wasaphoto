package api

import (
	"github.com/mouvzee/wasaphoto/service/database"
)
type Authorization struct {
	User	User
	Token	int
}

//crea sul database una func che permetta di controllare se il token è
//uguale all'user id, prima però devi aggiungere alla creazione
//dell'user il token che si genera e rimane salvato!!!!!



func (a *Authorization) FromDatabase (dbAuthorization database.Authorization) {

}