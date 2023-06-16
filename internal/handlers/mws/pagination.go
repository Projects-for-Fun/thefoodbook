package mws

import "net/http"

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// TODO
		next.ServeHTTP(rw, r)
	})
}
