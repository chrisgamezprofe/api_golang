package middleware

import (
	"net/http"

	"github.com/chrisgamezprofe/api_golang/auth"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := auth.ValidarToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}