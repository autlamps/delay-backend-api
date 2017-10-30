package api

import (
	"context"
	"log"
	"net/http"

	"github.com/autlamps/delay-backend-api/data"
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
			if err == data.ErrTokenInvalid {
				log.Printf("middleware - AuthUser: failed to parse token: %v\n", err)
			} else {
				log.Printf("middleware - AuthUser: failed to parse token: %v\n", err)
			}

			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(output.JSON401Response))
			return
		}

		ctx := context.WithValue(r.Context(), "Token", token)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (e *Env) CheckEmailConfirmed(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tk, ok := ctx.Value("Token").(data.Token)

		if !ok {
			log.Printf("CheckEmailConfirmed - Token from context not of type ctx: %v", tk)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(output.JSON500Response))
			return
		}

		u, err := e.Users.GetUser(tk.UserID)

		if err != nil {
			log.Printf("CheckEmailConfirmed - failed to get user: %v", tk)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(output.JSON500Response))
			return
		}

		if !u.EmailConfirmed {
			w.Header().Set("X-DELAY-CONFIRMED", "false")
		}

		h.ServeHTTP(w, r)
	})
}

// JSONContentType sets content type of request to json
func JSONContentType(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}
