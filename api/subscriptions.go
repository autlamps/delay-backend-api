package api

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/autlamps/delay-backend-api/data"
	"github.com/autlamps/delay-backend-api/output"
)

func (e *Env) CreateNewSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tk, ok := ctx.Value("Token").(data.Token)

	if !ok {
		log.Printf("Token from context not of type ctx: %v", tk)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	jns := struct {
		TripID          string     `json:"trip_id"`
		StopTimeID      string     `json:"stop_time_id"`
		Days            []data.Day `json:"days"`
		NotificationIDs []string   `json:"notification_ids"`
	}{}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&jns)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	ns := data.NewSubscription{
		UserID:          tk.UserID,
		TripID:          jns.TripID,
		StopTimeID:      jns.StopTimeID,
		Days:            jns.Days,
		NotificationIDs: jns.NotificationIDs,
	}

	s, err := e.Subscriptions.New(ns)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	st, err := e.StopTime.GetStopTimeByID(s.StopTimeID)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	s.StopTimeInfo = st

	rs := output.Response{
		Success: true,
		Errors:  nil,
		Meta:    output.GetMeta(),
		Result:  s,
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

func (e *Env) GetAllUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tk, ok := ctx.Value("Token").(data.Token)

	if !ok {
		log.Printf("Token from context not of type ctx: %v", tk)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	subs, err := e.Subscriptions.GetAll(tk.UserID)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	for i, _ := range subs {
		st, err := e.StopTime.GetStopTimeByID(subs[i].StopTimeID)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(output.JSON500Response))
			return
		}

		subs[i].StopTimeInfo = st
	}

	rs := output.Response{
		Success: true,
		Errors:  nil,
		Meta:    output.GetMeta(),
		Result:  subs,
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
