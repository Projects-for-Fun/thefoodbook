package webservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/auth"

	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/logging"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
)

type UserRequest struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"member_since"`
	LastLogin time.Time `json:"last_login"`
}

func (w *Webservice) HandleSignUp(rw http.ResponseWriter, r *http.Request) {
	var user UserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		MapErrorResponse(rw, r, err)
		return
	}

	user.Password, err = auth.EncryptPassword(user.Password)
	if err != nil {
		MapErrorResponse(rw, r, err)
		return
	}

	_, err = w.CreateUser(r.Context(), domain.User{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password, // encrypted
	})
	if err != nil {
		MapErrorResponse(rw, r, err)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (w *Webservice) HandleLogin(jwtKey []byte) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if !ok {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := w.LoginUser(r.Context(), username, password)
		if err != nil {
			MapErrorResponse(rw, r, err)
			return
		}

		expirationTime, tokenString, err := auth.CreateTokenForUser(*user, jwtKey)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(rw, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
	}
}

func (w *Webservice) HandleLogout(rw http.ResponseWriter, _ *http.Request) {
	http.SetCookie(rw, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
}

func (w *Webservice) HandleRefresh(jwtKey []byte) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		logger := logging.GetLogger(r.Context())
		token := auth.GetToken(r.Context())

		claims, ok := token.Claims.(*domain.Claims)
		if !ok && !token.Valid {
			logger.Info().Msg("Unauthorized user")
			rw.WriteHeader(http.StatusUnauthorized)
		}

		if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
			logger.Info().Msg("Token expired")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		expirationTime, tokenString, err := auth.CreateTokenFromExistingClaims(claims, jwtKey)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(rw, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
	}
}

func (w *Webservice) HandleWelcome(rw http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger(r.Context())
	token := auth.GetToken(r.Context())

	if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
		_, err := rw.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
		if err != nil {
			logger.Error().Err(err).Msg("Unexpected error")
		}
	} else {
		logger.Info().Msg("Unauthorized user")
		rw.WriteHeader(http.StatusUnauthorized)
	}
}
