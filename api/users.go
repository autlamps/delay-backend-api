package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/autlamps/delay-backend-api/data"
	"github.com/autlamps/delay-backend-api/output"
)

// CreateNewUser creates a new user
func (e *Env) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	nu := data.NewUser{}

	err := decoder.Decode(&nu)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	var user data.User

	if nu.Email == "" && nu.Password == "" {
		user, err = e.Users.NewAnonUser()

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(output.JSON500Response))
			return
		}
	} else {
		user, err = e.Users.NewUser(nu)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(output.JSON500Response))
			return
		}
	}

	token, err := e.Tokens.New(user.ID.String())

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	tks, err := e.Tokens.ToAuth(token)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	out := struct {
		ID        string
		Token     string
		CreatedOn int64
	}{
		user.ID.String(),
		tks,
		user.Created.Unix(),
	}

	rs := output.Response{
		Success: true,
		Result:  out,
		Errors:  output.Errors{},
		Meta:    output.GetMeta(),
	}

	rj, err := json.Marshal(rs)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	w.Write(rj)
}

// AuthenticateUser returns a token for a valid login
func (e *Env) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	pr := struct {
		Email    string
		Password string
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pr)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	u, err := e.Users.Authenticate(pr.Email, pr.Password)

	if err != nil {
		if err == data.ErrInvalidEmailOrPassword {
			rs := output.Response{
				Success: false,
				Result:  nil,
				Errors: output.Errors{
					Code: 1002,
					Msg:  "Incorrect email or password",
				},
				Meta: output.GetMeta(),
			}

			rj, err := json.Marshal(rs)

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(output.JSON500Response))
				return
			}

			w.Write(rj)
			return
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(output.JSON500Response))
			return
		}
	}

	tk, err := e.Tokens.New(u.ID.String())

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	tks, err := e.Tokens.ToAuth(tk)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	result := struct {
		UserID    string
		AuthToken string
	}{
		u.ID.String(),
		tks,
	}

	rs := output.Response{
		Success: true,
		Result:  result,
		Errors:  output.Errors{},
		Meta:    output.GetMeta(),
	}

	rj, err := json.Marshal(rs)

	w.Write(rj)
}
