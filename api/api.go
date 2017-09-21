package api

import (
	"fmt"
	"net/http"

	"database/sql"

	"github.com/autlamps/delay-backend-api/data"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/autlamps/delay-backend-api/static"
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
	Routes static.RouteStore
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
		Routes: static.RouteServiceInit(db),
	}

	r := mux.NewRouter()
	r.HandleFunc("/", CurrentRoutes)
	r.Handle("/users", alice.New(JSONContentType).ThenFunc(env.CreateNewUser)).Methods("POST")
	r.Handle("/tokens", alice.New(JSONContentType).ThenFunc(env.AuthenticateUser)).Methods("POST")
	r.Handle("/routes", alice.New(JSONContentType, env.AuthUser).ThenFunc(env.GetRoutes)).Methods("GET")

	return r, nil
}

// CurrentRoutes returns a simple html page listing what routes are currently available
func CurrentRoutes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<p>Create New User - POST /users</p><p>Authenitcate User - POST /tokens</p>")
}
