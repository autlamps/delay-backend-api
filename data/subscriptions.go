package data

import (
	"database/sql"
	"fmt"
	"time"

	"errors"

	"github.com/google/uuid"
)

var ErrNoNotificationMethods = errors.New("users - No notification methods specificed.")

// Day is one of our three letter day codes
type Day string

// Defines our three letter day codes
const (
	MONDAY    Day = "Mon"
	TUESDAY   Day = "Tue"
	WEDNESDAY Day = "Wed"
	THURSDAY  Day = "Thur"
	FRIDAY    Day = "Fri"
	SATURDAY  Day = "Sat"
	SUNDAY    Day = "Sun"
)

// NewSubscription is received from called and transformed into a db backed Subscription
type NewSubscription struct {
	TripID          string
	StopTimeID      string
	Days            []Day
	NotificationIDs []string
	UserID          string
}

// Subscription contains all subscription info
type Subscription struct {
	ID              string
	TripID          string
	StopTimeID      string
	UserID          string
	Archived        bool
	Created         time.Time
	Monday          bool
	Tuesday         bool
	Wednesday       bool
	Thursday        bool
	Friday          bool
	Saturday        bool
	Sunday          bool
	NotificationIDs []string
}

// SubscriptionStore defines methods for a concrete implementation
type SubscriptionStore interface {
	New(NewSubscription) (Subscription, error)
	Get(string) (Subscription, error)
	GetAll(string) ([]Subscription, error)
}

// SubscriptionService is our concrete psql implementation of the SubscriptionStore
type SubscriptionService struct {
	db *sql.DB
}

func InitSubscriptionService(db *sql.DB) *SubscriptionService {
	return &SubscriptionService{db: db}
}

// New creates a new database backed Subscription, or returns an error
func (ss *SubscriptionService) New(ns NewSubscription) (Subscription, error) {
	// If no notification methods are specified then we send an error back
	if len(ns.NotificationIDs) == 0 {
		return Subscription{}, ErrNoNotificationMethods
	}

	id, err := uuid.NewRandom()

	if err != nil {
		return Subscription{}, fmt.Errorf("subscriptions - New: failed to generate uuid: %v", err)
	}

	s := Subscription{
		ID:              id.String(),
		TripID:          ns.TripID,
		StopTimeID:      ns.StopTimeID,
		Archived:        false,
		Created:         time.Now().Round(time.Second),
		NotificationIDs: ns.NotificationIDs,
		UserID:          ns.UserID,
	}

	// Setup subscribed days, days not present will remain false as per golang default
	for _, d := range ns.Days {
		switch d {
		case MONDAY:
			s.Monday = true
		case TUESDAY:
			s.Tuesday = true
		case WEDNESDAY:
			s.Wednesday = true
		case THURSDAY:
			s.Thursday = true
		case FRIDAY:
			s.Friday = true
		case SATURDAY:
			s.Saturday = true
		case SUNDAY:
			s.Sunday = true
		}
	}

	_, err = ss.db.Exec("INSERT INTO subscription (sub_id, trip_id, stoptime_id, user_id, archived, date_created, monday, tuesday, wednesday, thursday, friday, saturday, sunday) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		s.ID,
		s.TripID,
		s.StopTimeID,
		s.UserID,
		s.Archived,
		s.Created,
		s.Monday,
		s.Tuesday,
		s.Wednesday,
		s.Thursday,
		s.Friday,
		s.Saturday,
		s.Sunday,
	)

	if err != nil {
		return Subscription{}, fmt.Errorf("users - Subscription: Failed to insert subscription into db: %v", err)
	}

	for _, sn := range s.NotificationIDs {
		_, err = ss.db.Exec("INSERT INTO sub_notification (sub_id, notification_id) VALUES ($1, $2)",
			s.ID,
			sn,
		)

		if err != nil {
			return Subscription{}, fmt.Errorf("users - New: failed to link notification methods and subscription: %v", err)
		}
	}

	return s, nil
}

// Get returns a single subscription by id
func (ss *SubscriptionService) Get(id string) (Subscription, error) {
	s := Subscription{}

	row := ss.db.QueryRow("SELECT sub_id, trip_id, stoptime_id, user_id, archived, date_created, monday, tuesday, wednesday, thursday, friday, saturday, sunday FROM subscription WHERE sub_id = $1", id)

	err := row.Scan(&s.ID, &s.TripID, &s.StopTimeID, &s.UserID, &s.Archived, &s.Created, &s.Monday, &s.Tuesday, &s.Wednesday, &s.Thursday, &s.Friday, &s.Saturday, &s.Sunday)

	if err != nil {
		if err == sql.ErrNoRows {
			return Subscription{}, fmt.Errorf("subscription - Get: No subscription with id: %v", id)
		}

		return Subscription{}, fmt.Errorf("subscription - Get: Failed to query db for subscription: %v", err)
	}

	s.Created = s.Created.Local()

	rows, err := ss.db.Query("SELECT notification_id from sub_notification WHERE sub_id = $1", s.ID)

	if err != nil {
		return Subscription{}, fmt.Errorf("subscription - Get: Failed get notification ids: %v", err)
	}

	for rows.Next() {
		var id string

		err := rows.Scan(&id)

		if err != nil {
			return Subscription{}, fmt.Errorf("subscription - Get: Failed to read individual notification id: %v", err)
		}

		s.NotificationIDs = append(s.NotificationIDs, id)
	}

	return s, nil
}
