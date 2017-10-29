package api

import (
	"fmt"
	"net/http"

	"database/sql"

	"github.com/autlamps/delay-backend-api/data"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/autlamps/delay-backend-api/objstore"
	"github.com/autlamps/delay-backend-api/static"
	_ "github.com/lib/pq"
)

type Conf struct {
	RDURL string
	DBURL string
	Key   string
}

type Env struct {
	Users            data.UserStore
	Tokens           data.TokenStore
	Routes           static.RouteStore
	StopTime         static.StopTimeStore
	ObjStore         objstore.Store
	NotificationInfo data.NotifyInfoStore
	Subscriptions    data.SubscriptionStore
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

	obj, err := objstore.InitService(c.RDURL)

	if err != nil {
		return nil, err
	}

	env := Env{
		Users:            data.InitUserService(db),
		Tokens:           data.InitTokenService(c.Key, db),
		Routes:           static.RouteServiceInit(db),
		StopTime:         static.StopTimeServiceInit(db),
		NotificationInfo: data.InitNotifyInfoService(db),
		Subscriptions:    data.InitSubscriptionService(db),
		ObjStore:         obj,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", CurrentRoutes)
	r.Handle("/users", alice.New(JSONContentType).ThenFunc(env.CreateNewUser)).Methods("POST")
	r.Handle("/tokens", alice.New(JSONContentType).ThenFunc(env.AuthenticateUser)).Methods("POST")
	r.Handle("/routes", alice.New(JSONContentType, env.AuthUser).ThenFunc(env.GetRoutes)).Methods("GET")
	r.Handle("/routes/{route_id}", alice.New(JSONContentType, env.AuthUser).ThenFunc(env.GetRoute)).Methods("GET")
	r.Handle("/delays", alice.New(JSONContentType, env.AuthUser).ThenFunc(env.GetDelays)).Methods("GET")
	r.Handle("/delays/subscribed", alice.New(JSONContentType, env.AuthUser).ThenFunc(env.GetSubedDelays)).Methods("GET")
	r.Handle("/stoptimes/{trip_id}", alice.New(JSONContentType, env.AuthUser).ThenFunc(env.GetStoptimesByTrip)).Methods("GET")
	r.Handle("/notifications", alice.New(JSONContentType, env.AuthUser).ThenFunc(env.CreateNotification)).Methods("POST")
	r.Handle("/notifications", alice.New(JSONContentType, env.AuthUser).ThenFunc(env.GetAllUserNotifications)).Methods("GET")
	r.Handle("/subscriptions", alice.New(JSONContentType, env.AuthUser).ThenFunc(env.CreateNewSubscription)).Methods("POST")
	r.Handle("/subscriptions", alice.New(JSONContentType, env.AuthUser).ThenFunc(env.GetAllUserSubscriptions)).Methods("GET")

	return r, nil
}

// CurrentRoutes returns a simple html page listing what routes are currently available
func CurrentRoutes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w,
		"<p>Create New User - POST /users</p>"+
			"<p>Authenitcate User - POST /tokens</p>"+
			"<p>Get All Routes - GET /routes</p>"+
			"<p>Get a Route with an ID - GET /routes/:route_id</p>"+
			"<p>Get Delays - GET /delays</p>"+
			"<p>Get stoptimes of a trip from it's ID - GET /stoptime/:trip_id</p>"+
			"<p>Create Notification Method - POST /notifications</p>"+
			"<p>Get all notifications - GET /notifications</p>"+
			"<p>Create Subscription - POST /subscriptions</p>"+
			"<p>Get all user subscriptions - GET /subscriptions</p>",
	)
}
