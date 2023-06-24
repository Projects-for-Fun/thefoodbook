package webservice

import (
	"net/http"

	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/logging"

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

	if err == domain.ErrInvalidUsernameOrPassword {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
}
