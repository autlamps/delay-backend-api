package data

import (
	"database/sql"
	"testing"
	"time"

	"fmt"

	_ "github.com/lib/pq"
)

func TestTokenService_New(t *testing.T) {
	db, err := sql.Open("postgres", dburl)
	defer db.Close()

	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping db: %v", err)
	}

	// Insert a new user so we can create their token
	_, err = db.Exec("INSERT into users (user_id, date_created) VALUEs ('9b1e4f4a-9776-4485-b627-071a4c012003', $1)", time.Now())

	if err != nil {
		t.Fatalf("Failed to insert testing user into db: %v", err)
	}

	ts := InitTokenService("hello", db)

	tk, err := ts.New("9b1e4f4a-9776-4485-b627-071a4c012003")

	if err != nil {
		t.Fatalf("Failed to create new token: %v", err)
	}

	row := db.QueryRow("SELECT token_id, user_id, date_created FROM tokens WHERE user_id = '9b1e4f4a-9776-4485-b627-071a4c012003'")

	dbtk := Token{}

	err = row.Scan(&dbtk.ID, &dbtk.UserID, &dbtk.Created)

	if err != nil {

		if err == sql.ErrNoRows {
			t.Fatalf("No token in database: %v", err)
		}

		t.Fatalf("Other error retrieving token: %v", err)
	}

	if dbtk.UserID != tk.UserID {
		t.Fatalf("Database token and created token don't have the same user_id")
	}

	if dbtk.ID != tk.ID {
		t.Fatalf("Database token and created token don't have the same token_id")
	}
}

//func TestTokenService_ToAuth(t *testing.T) {
//	ts := InitTokenService("hello", nil)
//
//	tk := Token{
//		ID:      "b50eb31d-4709-4df5-b65d-b6ddb88fea4a",
//		UserID:  "9b1e4f4a-9776-4485-b627-071a4c012003",
//		Created: time.Now(),
//	}
//
//	tks, err := ts.ToAuth(tk)
//
//	if err != nil {
//		t.Fatalf("Failed to turn token into token string: %v", err)
//	}
//
//	ot, err := ts.FromAuth(tks)
//
//	fmt.Println(ot)
//
//	//if tks != exp {
//	//	t.Fatalf("Generated token doesn't match. Expected %v, got %v", exp, tks)
//	//}
//
//}
