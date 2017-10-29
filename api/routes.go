package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/autlamps/delay-backend-api/output"
	"github.com/autlamps/delay-backend-api/static"
	"github.com/gorilla/mux"
)

func (e *Env) GetRoutes(w http.ResponseWriter, r *http.Request) {
	gr, err := e.Routes.GetRoutes()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	rs := output.Response{
		Success: true,
		Errors: output.Errors{
			Code: 1004,
			Msg:  "Bad behaviour warning",
		},
		Result: struct {
			Count  int           `json:"count"`
			Routes static.Routes `json:"routes"`
		}{
			len(gr),
			gr,
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
}

func (e *Env) GetRoute(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	route_id := vars["route_id"]

	route, err := e.Routes.GetRouteByID(route_id)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	result := struct {
		Routes static.Route
	}{
		route,
	}

	gs := output.Response{
		Success: true,
		Result:  result,
		Errors:  nil,
		Meta:    output.GetMeta(),
	}

	gr, err := json.Marshal(gs)

	w.Write(gr)
}

func (e *Env) GetRouteTrips(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	route_id := vars["route_id"]

	route, err := e.Routes.GetRouteByID(route_id)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	trips, err := e.Trips.GetTripsByRouteID(route.ID)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	result := struct {
		Route static.Route `json:"route"`
		Trips static.Trips `json:"trips"`
	}{
		Route: route,
		Trips: trips,
	}

	gs := output.Response{
		Success: true,
		Result:  result,
		Errors:  nil,
		Meta:    output.GetMeta(),
	}

	gr, err := json.Marshal(gs)

	w.Write(gr)
}
