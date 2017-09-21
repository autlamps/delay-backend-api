package api

import (
	"context"
	"log"
	"net/http"

	"github.com/autlamps/delay-backend-api/output"
)

// Auth user checks the token provided in the header
func (e *Env) AuthUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("X-DELAY-AUTH")

		if auth == "" {
			log.Println("middleware - AuthUser: no auth header included")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(output.JSON403Response))
			return
		}

		token, err := e.Tokens.FromAuth(auth)

		if err != nil {
			log.Printf("middleware - AuthUser: failed to parse token: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(output.JSON401Response))
			return
		}

		ctx := context.WithValue(r.Context(), "Token", token)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// JSONContentType sets content type of request to json
func JSONContentType(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}
