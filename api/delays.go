package api

import (
	"log"
	"net/http"

	"encoding/json"
	"fmt"

	"github.com/autlamps/delay-backend-api/data"
	"github.com/autlamps/delay-backend-api/delays"
	"github.com/autlamps/delay-backend-api/output"
)

// Get delays simply returns trips running late from the collection service
func (e *Env) GetDelays(w http.ResponseWriter, r *http.Request) {
	bd, err := e.ObjStore.Get("delays")

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	var d delays.Out

	if err := json.Unmarshal(bd, &d); err != nil {
		log.Println(fmt.Errorf("api - Delays: failed to parse json: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	resp := output.Response{
		Success: true,
		Errors:  nil,
		Result:  d,
		Meta:    output.GetMeta(),
	}

	rb, err := json.Marshal(resp)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	w.Write(rb)
}

func (e *Env) GetSubedDelays(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tk, ok := ctx.Value("Token").(data.Token)

	if !ok {
		log.Printf("Token from context not of type ctx: %v", tk)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	del, err := e.ObjStore.Get("delays")

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	var od delays.Out

	if err := json.Unmarshal(del, &od); err != nil {
		log.Println(fmt.Errorf("api - Delays: failed to parse json: %v", err))
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

	var subedDelays []delays.OutTrip

	// Looping over all delayed trips then all subscribed trips.
	// If they equal added to subedDelays
	for _, d := range od.Trips {
		for _, s := range subs {
			if d.TripID == s.TripID {
				subedDelays = append(subedDelays, d)
			}
		}
	}

	// Reinsert data into outDelays for export

	od.Trips = subedDelays
	od.Count = len(subedDelays)

	resp := output.Response{
		Success: true,
		Errors:  nil,
		Result:  od,
		Meta:    output.GetMeta(),
	}

	rb, err := json.Marshal(resp)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	w.Write(rb)
}
