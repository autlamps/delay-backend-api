package data

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

func subSetup() (NotifyInfo, User, *sql.DB, error) {
	db, err := sql.Open("postgres", dburl)

	if err != nil {
		return NotifyInfo{}, User{}, nil, fmt.Errorf("Failed to connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		return NotifyInfo{}, User{}, nil, fmt.Errorf("Failed to ping db: %v", err)
	}

	us := InitUserService(db)

	u, err := us.NewUser(NewUser{
		"Bobby Tables",
		"bobby.tables@example.com",
		"correcthorsebatterystaple",
	})

	if err != nil {
		return NotifyInfo{}, User{}, nil, fmt.Errorf("Failed to create new user: %v", err)
	}

	ns := InitNotifyInfoService(db)

	n, err := ns.New(u.ID.String(), "p", "iPhone X", "1234456")

	if err != nil {
		return NotifyInfo{}, User{}, nil, fmt.Errorf("Failed to create new notify method: %v", err)
	}

	return n, u, db, nil
}

func TestSubscriptionService_New(t *testing.T) {
	n, u, db, err := subSetup()

	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	ss := InitSubscriptionService(db)

	ns := NewSubscription{
		TripID:          "df688c57-987c-4705-9e22-936342eb6e3f",
		StopTimeID:      "5cce0bca-d489-43d7-b3cb-48e0df054c8a",
		Days:            []Day{"Mon", "Tue", "Wed"},
		NotificationIDs: []string{n.ID},
		UserID:          u.ID.String(),
	}

	s, err := ss.New(ns)

	if err != nil {
		t.Fatalf("Failed to create subscription: %v", err)
	}

	row := db.QueryRow("SELECT sub_id, trip_id, stoptime_id, user_id, archived, date_created, monday, tuesday, wednesday, thursday, friday, saturday, sunday FROM subscription WHERE sub_id = $1", s.ID)

	dbs := Subscription{}

	err = row.Scan(&dbs.ID, &dbs.TripID, &dbs.StopTimeID, &dbs.UserID, &dbs.Archived, &dbs.Created, &dbs.Monday, &dbs.Tuesday, &dbs.Wednesday, &dbs.Thursday, &dbs.Friday, &dbs.Saturday, &dbs.Sunday)

	if err != nil {
		t.Fatalf("Failed to read in sub from db: %v", err)
	}

	dbs.Created = dbs.Created.Local()

	rows, err := db.Query("SELECT notification_id from sub_notification WHERE sub_id = $1", s.ID)

	if err != nil {
		t.Fatalf("Failed to read notification ids: %v", err)
	}

	for rows.Next() {
		var id string

		err := rows.Scan(&id)

		if err != nil {
			t.Fatalf("Failed to read in individual notification id: %v", err)
		}

		dbs.NotificationIDs = append(dbs.NotificationIDs, id)
	}

	if !reflect.DeepEqual(s, dbs) {
		t.Fatalf("Sub not equal to db sub: Expected %v, got %v", s, dbs)
	}

	if dbs.Monday == false {
		t.Fatalf("New Subscription differs from one returned from db: Expected %v, got %v", ns, dbs)
	}

	if dbs.Tuesday == false {
		t.Fatalf("New Subscription differs from one returned from db: Expected %v, got %v", ns, dbs)
	}

	if dbs.Wednesday == false {
		t.Fatalf("New Subscription differs from one returned from db: Expected %v, got %v", ns, dbs)
	}

	if ns.UserID != dbs.UserID {
		t.Fatalf("New Subscription differs from one returned from db: Expected %v, got %v", ns, dbs)
	}

	if ns.TripID != dbs.TripID {
		t.Fatalf("New Subscription differs from one returned from db: Expected %v, got %v", ns, dbs)
	}

	if ns.StopTimeID != dbs.StopTimeID {
		t.Fatalf("New Subscription differs from one returned from db: Expected %v, got %v", ns, dbs)
	}

	if !reflect.DeepEqual(ns.NotificationIDs, dbs.NotificationIDs) {
		t.Fatalf("New Subscription notification ids differs from one returned from db: Expected %v, got %v", ns.NotificationIDs, dbs.NotificationIDs)

	}

	err = subCleanUp(s, n, u, db)

	if err != nil {
		t.Fatalf("Failed to clean up: %v", err)
	}
}

func TestSubscriptionService_Get(t *testing.T) {
	n, u, db, err := subSetup()

	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	ss := InitSubscriptionService(db)

	ns := NewSubscription{
		TripID:          "df688c57-987c-4705-9e22-936342eb6e3f",
		StopTimeID:      "5cce0bca-d489-43d7-b3cb-48e0df054c8a",
		Days:            []Day{"Mon", "Tue", "Wed"},
		NotificationIDs: []string{n.ID},
		UserID:          u.ID.String(),
	}

	s, err := ss.New(ns)

	if err != nil {
		t.Fatalf("Failed to create subscription: %v", err)
	}

	dbs, err := ss.Get(s.ID)

	if err != nil {
		t.Fatalf("Failed to get sub from db: %v", err)
	}

	if !reflect.DeepEqual(s, dbs) {
		t.Fatalf("DB sub and sub different. Expected %v, got %v", s, dbs)
	}

	err = subCleanUp(s, n, u, db)

	if err != nil {
		t.Fatalf("Failed to clean up")
	}
}

func TestSubscriptionService_GetAll(t *testing.T) {
	n, u, db, err := subSetup()

	if err != nil {
		t.Fatalf("Failed to setup: %v", err)
	}

	ss := InitSubscriptionService(db)

	subs := []Subscription{}

	stopTimeIDs := []string{
		"5cce0bca-d489-43d7-b3cb-48e0df054c8a",
		"9a5a2eae-1ed4-49b5-abfe-519bfefc6300",
		"7225fa2b-cba3-4b4f-a0f8-aecb1b1dfb70",
		"fd91dca8-aaef-4548-bed4-fa7d835de25c",
		"98a3144f-7ed3-419b-9276-ff747c46b982",
	}

	for i := 0; i < 5; i++ {
		ns := NewSubscription{
			TripID:          "df688c57-987c-4705-9e22-936342eb6e3f",
			StopTimeID:      stopTimeIDs[i],
			Days:            []Day{"Mon", "Tue", "Wed"},
			NotificationIDs: []string{n.ID},
			UserID:          u.ID.String(),
		}

		s, err := ss.New(ns)

		if err != nil {
			t.Fatalf("Failed to create new sub: %v", err)
		}

		subs = append(subs, s)
	}

	dbsubs, err := ss.GetAll(u.ID.String())

	if err != nil {
		t.Fatalf("Failed to get subs from db: %v", err)
	}

	if !reflect.DeepEqual(subs, dbsubs) {
		t.Fatalf("Subs returned from GetAll not the same as expected")
	}

	// CleanUp
	for _, sub := range subs {
		_, err := db.Exec("DELETE FROM sub_notification WHERE sub_id = $1", sub.ID)

		if err != nil {
			t.Fatalf("Failed to delete created sub notification: %v\n", err)
		}

		_, err = db.Exec("DELETE FROM subscription WHERE sub_id = $1", sub.ID)

		if err != nil {
			t.Fatalf("Failed to delete created sub: %v\n", err)
		}
	}

	err = subCleanUp(Subscription{}, n, u, db)

	if err != nil {
		t.Fatalf("Failed to clean up: %v", err)
	}

}

func subCleanUp(s Subscription, ni NotifyInfo, u User, db *sql.DB) error {
	// We can send a blank Subscription struct if we have already cleaned up
	if s.ID != "" {
		_, err := db.Exec("DELETE FROM sub_notification WHERE sub_id = $1", s.ID)

		if err != nil {
			return fmt.Errorf("Failed to delete created subs: %v\n", err)
		}

		_, err = db.Exec("DELETE FROM subscription WHERE sub_id = $1", s.ID)

		if err != nil {
			return fmt.Errorf("Failed to delete created sub: %v\n", err)
		}
	}

	_, err := db.Exec("DELETE FROM notification WHERE notification_id = $1", ni.ID)

	if err != nil {
		return fmt.Errorf("Failed to delete created notification: %v\n", err)
	}
	_, err = db.Exec("DELETE FROM tokens WHERE user_id = $1", u.ID)

	if err != nil {
		return fmt.Errorf("Failed to delete created user token: %v\n", err)
	}

	_, err = db.Exec("DELETE FROM users WHERE user_id = $1", u.ID)

	if err != nil {
		return fmt.Errorf("Failed to delete created user: %v\n", err)
	}

	db.Close()

	return nil
}
