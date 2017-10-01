package api

import (
	"testing"

	"database/sql"
	"net/http"

	"net/http/httptest"

	"fmt"

	"io/ioutil"

	"github.com/autlamps/delay-backend-api/data"
	"github.com/autlamps/delay-backend-api/output"
	"github.com/justinas/alice"
)

func TestEnv_AuthUser(t *testing.T) {
	//Setup
	db, err := sql.Open("postgres", dburl)
	defer db.Close()

	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping db: %v", err)
	}

	env := Env{
		Tokens: data.InitTokenService("hello", db),
		Users:  data.InitUserService(db),
	}

	u, err := env.Users.NewUser(data.NewUser{
		"Bobby Tables",
		"bobby.tables@example.com",
		"correcthorsebatterystaple",
	})

	if err != nil {
		t.Fatalf("Failed to create new user for testing: %v", err)
	}

	ctk, err := env.Tokens.New(u.ID.String())

	if err != nil {
		t.Fatalf("Failed to create correct token: %v", err)
	}

	ctks, err := env.Tokens.ToAuth(ctk)

	if err != nil {
		t.Fatalf("Failed to create correct token string: %v", err)
	}

	server := httptest.NewServer(alice.New(env.AuthUser).ThenFunc(env.fakeHTTP))

	//Test
	tests := []struct {
		Token  string
		Result string
	}{
		{ctks, u.ID.String()},
		{"", output.JSON403Response},
		{"1234567", output.JSON401Response},
	}

	for i, test := range tests {
		req, err := http.NewRequest("GET", server.URL, nil)

		if err != nil {
			t.Errorf("%v - Failed to create new req: %v", i, err)
		}

		req.Header.Set("X-DELAY-AUTH", test.Token)

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			t.Errorf("%v - Failed to do request: %v", i, err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			t.Errorf("%v - Failed to read body: %v", i, err)
		}

		bodys := string(body)

		if bodys != test.Result {
			t.Errorf("%v - Body doesn't match result. Expected %v, got %v", i, test.Result, bodys)
		}
	}

	//Clean up
	_, err = db.Exec("DELETE FROM tokens WHERE user_id = $1", u.ID)

	if err != nil {
		fmt.Printf("Failed to delete created user token: %v\n", err)
	}

	_, err = db.Exec("DELETE FROM users WHERE user_id = $1", u.ID)

	if err != nil {
		fmt.Printf("Failed to delete created user: %v\n", err)
	}
}

func (env *Env) fakeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tk, ok := ctx.Value("Token").(data.Token)

	if !ok {
		w.Write([]byte("failed"))
		return
	}

	w.Write([]byte(tk.UserID))
}
