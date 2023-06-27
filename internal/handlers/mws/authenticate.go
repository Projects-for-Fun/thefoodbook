package mws

import (
	"net/http"

	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/auth"

	"github.com/Projects-for-Fun/thefoodbook/internal/core/domain"
	"github.com/Projects-for-Fun/thefoodbook/pkg/sys/logging"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtKey []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger := logging.GetLogger(r.Context())

			c, err := r.Cookie("token")
			if err != nil {
				logger.Info().AnErr("error", err).Msg("Unauthorized user. Cookie not set.")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tkn, err := jwt.ParseWithClaims(c.Value, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					logger.Info().AnErr("error", err).Msg("Unauthorized user. Invalid signature.")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				logger.Info().AnErr("error", err).Msg("Bad request.")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if !tkn.Valid {
				logger.Info().AnErr("error", err).Msg("Unauthorized user. Token is not valid.")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			logger.Info().Msg("User logged in")
			next.ServeHTTP(w, r.WithContext(auth.AttachToken(r.Context(), tkn)))
		}

		return http.HandlerFunc(fn)
	}
}
