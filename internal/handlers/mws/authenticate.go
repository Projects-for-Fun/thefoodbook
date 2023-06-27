package mws

import (
	"net/http"

	"github.com/Projects-for-Fun/thefoodbook/internal/handlers/webservice"

	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/auth"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/logging"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtKey []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			logger := logging.GetLogger(r.Context())

			c, err := r.Cookie("token")
			if err != nil {
				webservice.MapErrorResponse(rw, r, domain.ErrUnauthorized, "Unauthorized user. Cookie not set.")
				return
			}

			tkn, err := jwt.ParseWithClaims(c.Value, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					webservice.MapErrorResponse(rw, r, domain.ErrUnauthorized, "Unauthorized user. Invalid signature.")
					return
				}

				webservice.MapErrorResponse(rw, r, domain.ErrBadRequest, "Unauthorized user. Bad request.")
				return
			}

			if !tkn.Valid {
				webservice.MapErrorResponse(rw, r, domain.ErrBadRequest, "Unauthorized user. Token is not valid.")
				return
			}

			logger.Info().Msg("User logged in")
			next.ServeHTTP(rw, r.WithContext(auth.AttachToken(r.Context(), tkn)))
		}

		return http.HandlerFunc(fn)
	}
}
