package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/autlamps/delay-backend-api/output"
	"github.com/autlamps/delay-backend-api/static"
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
			Count  int
			Routes static.Routes
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
