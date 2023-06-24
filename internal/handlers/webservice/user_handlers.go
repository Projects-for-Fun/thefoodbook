package webservice

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Projects-for-Fun/thefoodbook/internal/handlers/mws"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
)

type UserRequest struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	CreatedAt  time.Time `json:"member_since"`
	ModifiedAt time.Time `json:"last_login"`
}

func (w *Webservice) HandleSignUp(rw http.ResponseWriter, r *http.Request) {
	var user UserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		MapErrorResponse(rw, r, err)
		return
	}

	// 		encryptedPassword, err := encryptPassword(plainPassword, as.saltRounds)
	//		if err != nil {
	//			return nil, err
	//		}
	user.Password, err = encryptPassword(user.Password, 10)
	if err != nil {
		MapErrorResponse(rw, r, err)
	}

	// Get logger from context
	userID, err := w.CreateUser(r.Context(), mws.GetLogger(r.Context()), domain.User(user))
	if err != nil {
		MapErrorResponse(rw, r, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	userIDString := *userID
	_, err = rw.Write([]byte(userIDString.String()))
	if err != nil {
		MapErrorResponse(rw, r, err)
		return
	}

}
