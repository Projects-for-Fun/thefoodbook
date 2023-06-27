package webservice

import (
	"net/http"

	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/logging"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"

	"github.com/samber/lo"
)

func MapErrorResponse(rw http.ResponseWriter, r *http.Request, err error, info ...interface{}) {
	logger := logging.GetLogger(r.Context())

	wrappedErrors, isWrapped := err.(interface{ Unwrap() []error })

	if isWrapped {
		var errs []error

		for _, err := range wrappedErrors.Unwrap() {
			errs = append(errs, err)

			logger.Info().AnErr("error", err).Msg(err.Error())
		}

		addErrorHeader(rw, errs[0])
	}

	// If single error, check if there is an error message in info
	if !isWrapped {
		msg := lo.TernaryF(info == nil, func() string { return err.Error() }, func() string { return info[0].(string) })
		logger.Info().AnErr("error", err).Msg(msg)

		addErrorHeader(rw, err)
	}
}

func addErrorHeader(rw http.ResponseWriter, err error) {
	if err == domain.ErrUserExists {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err == domain.ErrInvalidUsernameOrPassword {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	if err == domain.ErrBadRequest {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
}
