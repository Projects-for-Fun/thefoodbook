package webservice

import (
	"net/http"

	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/logging"

	"golang.org/x/crypto/bcrypt"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
)

func MapErrorResponse(rw http.ResponseWriter, r *http.Request, err error) {
	logger := logging.GetLogger(r.Context())

	wrappedErrors, isWrapped := err.(interface{ Unwrap() []error })

	if isWrapped {
		var errs []error

		for _, err := range wrappedErrors.Unwrap() {
			errs = append(errs, err)
			logger.Info().AnErr("error", err).Msg(err.Error())
		}

		header(rw, errs[0])
	}

	if !isWrapped {
		header(rw, err)
	}
}

func header(rw http.ResponseWriter, err error) {
	if err == domain.ErrUserExists {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
}

func encryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
