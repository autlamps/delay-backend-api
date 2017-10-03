package api

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/autlamps/delay-backend-api/data"
	"github.com/autlamps/delay-backend-api/output"
)

func (e *Env) CreateNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tk, ok := ctx.Value("Token").(data.Token)

	if !ok {
		log.Printf("Token from context not of type ctx: %v", tk)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	jni := struct {
		Type  data.NotifyType
		Value string
		Name  string
	}{}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&jni)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	ni, err := e.NotificationInfo.New(tk.UserID, jni.Type, jni.Name, jni.Value)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	rs := output.Response{
		Success: true,
		Errors:  nil,
		Meta:    output.GetMeta(),
		Result:  ni,
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

func (e *Env) GetAllUserNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tk, ok := ctx.Value("Token").(data.Token)

	if !ok {
		log.Printf("Token from context not of type ctx: %v", tk)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(output.JSON500Response))
		return
	}

	uni, err := e.NotificationInfo.GetAll(tk.UserID)

	rs := output.Response{
		Success: true,
		Errors:  nil,
		Meta:    output.GetMeta(),
		Result:  uni,
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
