package api

import (
	"fmt"
	"net/http"

	"database/sql"

	"github.com/autlamps/delay-backend-api/data"
	"github.com/gorilla/mux"

	"encoding/json"

	"log"

	"github.com/autlamps/delay-backend-api/output"
	_ "github.com/lib/pq"
)

type Conf struct {
	RDURL string
	DBURL string
	Key   string
}

type Env struct {
	Users  data.UserStore
	Tokens data.TokenStore
}

// Create returns a router ready to handle requests
func Create(c Conf) (*mux.Router, error) {
	db, err := sql.Open("postgres", c.DBURL)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	env := Env{
		Users:  data.InitUserService(db),
		Tokens: data.InitTokenService(c.Key, db),
	}

	r := mux.NewRouter()
	r.HandleFunc("/", CurrentRoutes)
	r.HandleFunc("/users", env.CreateNewUser).Methods("POST")

	return r, nil
}

// CurrentRoutes returns a simple html page listing what routes are currently available
func CurrentRoutes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<p>Create New User - POST /users</p>")
}

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
