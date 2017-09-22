package api

import (
	"log"
	"net/http"

	"encoding/json"
	"fmt"

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
