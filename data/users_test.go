package data

import (
	"database/sql"
	"fmt"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestUserService_NewUser(t *testing.T) {
	db, err := sql.Open("postgres", dburl)
	defer db.Close()

	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping db: %v", err)
	}

	us := InitUserService(db)

	nu := NewUser{
		Email:    "bobby.tables@example.com",
		Name:     "Bobby Tables",
		Password: "correcthorsebatterystaple",
	}

	u, err := us.NewUser(nu)

	if err != nil {
		t.Fatalf("Failed to insert user into db: %v", err)
	}

	if u.Name != nu.Name {
		t.Fatalf("New users name not the same as returned. Expected %v, got %v", nu.Name, u.Name)
	}

	if u.Email != nu.Email {
		t.Fatalf("New users email not the same as returned. Expected %v, got %v", nu.Email, u.Email)
	}

	if err := bcrypt.CompareHashAndPassword(u.Password, []byte(nu.Password)); err != nil {
		t.Fatalf("New users password not the same as returned.")
	}

	dbu := User{}

	row := us.db.QueryRow("SELECT user_id, email, name, password, date_created FROM users WHERE email = $1", nu.Email)

	err = row.Scan(&dbu.ID, &dbu.Email, &dbu.Name, &dbu.Password, &dbu.Created)

	if err != nil {
		t.Fatalf("Failed to get new user from db: %v", err)
	}

	if dbu.Name != nu.Name {
		t.Fatalf("New users name not the same as database. Expected %v, got %v", nu.Name, dbu.Name)
	}

	if nu.Email != dbu.Email {
		t.Fatalf("New users email not the same as database. Expected %v, got %v", nu.Email, dbu.Email)
	}

	if err := bcrypt.CompareHashAndPassword(dbu.Password, []byte(nu.Password)); err != nil {
		t.Fatalf("New users password not the same as database.")
	}

	//Clean Up
	_, err = db.Exec("DELETE FROM tokens WHERE user_id = $1", u.ID)

	if err != nil {
		fmt.Printf("Failed to delete created user token: %v\n", err)
	}

	_, err = db.Exec("DELETE FROM users WHERE user_id = $1", u.ID)

	if err != nil {
		fmt.Printf("Failed to delete created user: %v\n", err)
	}
}
