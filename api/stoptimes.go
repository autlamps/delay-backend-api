package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/autlamps/delay-backend-api/output"
	"github.com/autlamps/delay-backend-api/static"
	"github.com/gorilla/mux"
)

func (e *Env) GetTrip(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	trip_ip := vars["trip_id"]

	sts, err := e.StopTime.GetStopTimesByTripID(trip_ip)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	rs := output.Response{
		Success: true,
		Result: struct {
			Count     int                  `json:"count"`
			StopTimes static.StopTimeArray `json:"stop_time"`
		}{
			len(sts),
			sts,
		},
		Errors: nil,
		Meta:   output.GetMeta(),
	}
	rj, err := json.Marshal(rs)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
	}

	w.Write(rj)
}
