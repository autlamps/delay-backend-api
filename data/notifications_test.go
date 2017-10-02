package data

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestNotifyInfoService_New(t *testing.T) {
	db, err := sql.Open("postgres", dburl)
	defer db.Close()

	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping db: %v", err)
	}

	us := InitUserService(db)

	u, err := us.NewUser(NewUser{
		"Bobby Tables",
		"bobby.tables@example.com",
		"correcthorsebatterystaple",
	})

	if err != nil {
		t.Fatalf("Failed to create new user: %v", err)
	}

	nis := InitNotifyInfoService(db)

	ni, err := nis.New(u.ID.String(), PUSH, "iPhone X", "123467876543214")

	if err != nil {
		t.Fatalf("Failed to insert new notification info into db: %v", err)
	}

	dbni := NotifyInfo{}

	row := db.QueryRow("SELECT notification_id, user_id, type, name, value, date_created FROM notification WHERE notification_id = $1", ni.ID)

	err = row.Scan(&dbni.ID, &dbni.UserID, &dbni.Type, &dbni.Name, &dbni.Value, &dbni.Created)

	if err != nil {
		t.Fatalf("Failed to retrieve notify info from db: %v", err)
	}

	dbni.Created = dbni.Created.In(time.Local)

	if !reflect.DeepEqual(ni, dbni) {
		t.Fatalf("Retrieved not the same as saved. Expected %v, got %v", ni, dbni)
	}

	// Cleanup
	_, err = db.Exec("DELETE FROM notification WHERE notification_id = $1", ni.ID)

	if err != nil {
		fmt.Printf("Failed to delete created notification: %v\n", err)
	}

	_, err = db.Exec("DELETE FROM tokens WHERE user_id = $1", u.ID)

	if err != nil {
		fmt.Printf("Failed to delete created user token: %v\n", err)
	}

	_, err = db.Exec("DELETE FROM users WHERE user_id = $1", u.ID)

	if err != nil {
		fmt.Printf("Failed to delete created user: %v\n", err)
	}
}
