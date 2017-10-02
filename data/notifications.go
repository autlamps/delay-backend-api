package data

import (
	"database/sql"
	"fmt"

	"time"

	"encoding/json"

	"github.com/google/uuid"
)

// NotifyType is the type of notification method to be used
type NotifyType string

const (
	EMAIL NotifyType = "e"
	TXT   NotifyType = "t"
	PUSH  NotifyType = "p"
)

// NotifyInfo stores information used to notify user of someting
type NotifyInfo struct {
	ID      string     `json:"id"`
	UserID  string     `json:"user_id"`
	Type    NotifyType `json:"type"`
	Name    string     `json:"name"`
	Value   string     `json:"value"`
	Created time.Time  `json:"-"`
}

// MarshalJSON to output date created as at unix timestamp
func (ni *NotifyInfo) MarshalJSON() ([]byte, error) {
	type Notify NotifyInfo

	c := ni.Created.Unix()

	jni := struct {
		*Notify
		DateCreated int64 `json:"date_created"`
	}{
		Notify:      (*Notify)(ni),
		DateCreated: c,
	}

	return json.Marshal(jni)
}

// NotifyInfoStore is our interface for implementing a concrete service
type NotifyInfoStore interface {
	New(uid string, t NotifyType, n, v string) (NotifyInfo, error)
	Get(id string) (NotifyInfo, error)
	GetAll(uid string) ([]NotifyInfo, error)
	Delete(id string) error
}

// NotifyInfoService is our concrete psql implementation of the NotifyInfoStore
type NotifyInfoService struct {
	db *sql.DB
}

// InitNotifyInfoService returns a ready NotifyInfoService
func InitNotifyInfoService(db *sql.DB) *NotifyInfoService {
	return &NotifyInfoService{db: db}
}

// New creates a new NotifyInfo and stores it in the db
func (ns *NotifyInfoService) New(uid string, t NotifyType, n, v string) (NotifyInfo, error) {
	id, err := uuid.NewRandom()

	if err != nil {
		return NotifyInfo{}, fmt.Errorf("notifications - New: Failed to generate uuid: %v", err)
	}

	ni := NotifyInfo{
		ID:      id.String(),
		UserID:  uid,
		Type:    t,
		Name:    n,
		Value:   v,
		Created: time.Now().Round(time.Second),
	}

	_, err = ns.db.Exec("INSERT INTO notification (notification_id, user_id, type, name, value, date_created) VALUES ($1, $2, $3, $4, $5, $6)",
		ni.ID,
		ni.UserID,
		ni.Type,
		ni.Name,
		ni.Value,
		ni.Created,
	)

	if err != nil {
		return NotifyInfo{}, fmt.Errorf("notifications - New: failed to insert into db: %v", err)
	}

	return ni, nil
}

