package static

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// StopTime represents a specific stop on a trip. Also includes embedded Stop info
type StopTime struct {
	ID           string    `json:"id"`
	TripID       string    `json:"trip_id"`
	Arrival      time.Time `json:"arrival"`
	Departure    time.Time `json:"departure"`
	StopSequence int       `json:"stop_sequence"`
	StopInfo     Stop      `json:"stop_info"`
}

// MarshalJSON for StopTime
func (s StopTime) MarshalJSON() ([]byte, error) {
	type ST StopTime

	js := struct {
		ST
		Departure string `json:"departure"`
		Arrival   string `json:"arrival"`
	}{
		ST:        (ST)(s),
		Departure: s.Departure.Format("15:04:05"),
		Arrival:   s.Arrival.Format("15:04:05"),
	}

	return json.Marshal(js)
}

// Stop represents a physical stop. Embedded into StopTime instead of being its own service for ease of use
type Stop struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

// StopTimeArray is simply a slice of StopTime
type StopTimeArray []StopTime

// StopTimeStore defines the methods that a concrete StopTimeService should implement
type StopTimeStore interface {
	GetStopTimesByTripID(tripID string) (StopTimeArray, error)
	GetStopTimeByID(id string) (StopTime, error)
	getStopByID(id string) (Stop, error)
}

// StopTimeService implements StopTimeStore in PSQL
type StopTimeService struct {
	db *sql.DB
}

// StopTimeServiceInit initializes a new StopTimeService
func StopTimeServiceInit(db *sql.DB) *StopTimeService {
	return &StopTimeService{db: db}
}

// GetStopTimesByTripID returns all stops of the given trip
func (sts *StopTimeService) GetStopTimesByTripID(tripID string) (StopTimeArray, error) {
	var sta StopTimeArray

	rows, err := sts.db.Query("SELECT stoptime_id, trip_id, arrival_time, departure_time, stop_id, "+
		"stop_sequence from stop_times WHERE trip_id = $1 ORDER BY stop_sequence ASC", tripID)

	if err != nil {
		return sta, err
	}

	for rows.Next() {
		st := StopTime{}
		var stopID string

		if err := rows.Scan(&st.ID, &st.TripID, &st.Arrival, &st.Departure, &stopID, &st.StopSequence); err != nil {
			return sta, err // TODO: decide what to do here. Do we inject logger and log it? Stop execution?
		}

		st.StopInfo, err = sts.getStopByID(stopID)

		if err != nil {
			return sta, err
		}

		sta = append(sta, st)
	}

	return sta, nil
}

// GetStopTimesByID returns a singular stoptime with stop info
func (sts *StopTimeService) GetStopTimeByID(id string) (StopTime, error) {
	st := StopTime{}
	var stopID string

	row := sts.db.QueryRow("SELECT stoptime_id, trip_id, arrival_time, departure_time, stop_id, stop_sequence from stop_times WHERE stoptime_id = $1", id)

	err := row.Scan(&st.ID, &st.TripID, &st.Arrival, &st.Departure, &stopID, &st.StopSequence)

	if err != nil {
		return StopTime{}, fmt.Errorf("stoptime - GetStopTimeByID: Failed to get stoptime: %v", err)
	}

	stop, err := sts.getStopByID(stopID)

	if err != nil {
		return StopTime{}, fmt.Errorf("stoptime - GetStopTimesByID: Failed to get stop info: %v", err)
	}

	st.StopInfo = stop

	return st, nil
}

func (sts *StopTimeService) getStopTimesByID(id string) (StopTimeArray, error) {
	str := StopTimeArray{}

	rows, err := sts.db.Query("SELECT stoptime_id, arrival_time, departure_time, stop_id, stop_sequence FROM stop_times WHERE trip_id = $1", id)
	if err != nil {
		return str, fmt.Errorf("stoptime - GetStopTimes: %v", err)
	}

	for rows.Next() {
		st := StopTime{}
		err := rows.Scan(&st.ID, &st.Arrival, &st.Departure, &st.StopInfo, &st.StopSequence)
		if err != nil {
			return str, fmt.Errorf("stoptime - GetStoptimes: Failed to scan: %v", err)
		}
		str = append(str, st)
	}

	return str, nil
}

// getStopByID returns a single stop
func (sts *StopTimeService) getStopByID(id string) (Stop, error) {
	s := Stop{}

	row := sts.db.QueryRow("SELECT stop_id, stop_name, stop_lat, stop_lon FROM stops WHERE stop_id = $1", id)

	err := row.Scan(&s.ID, &s.Name, &s.Lat, &s.Lon)

	if err != nil {
		return s, err
	}

	return s, nil
}

// IsEqual returns true if the given StopTime is equal to the StopTime the method is run on
func (st StopTime) IsEqual(x StopTime) bool {

	if st.ID != x.ID {
		return false
	}
	if st.TripID != x.TripID {
		return false
	}
	if !st.Arrival.Equal(x.Arrival) {
		return false
	}
	if !st.Departure.Equal(x.Departure) {
		return false
	}
	if st.StopSequence != x.StopSequence {
		return false
	}
	if !st.StopInfo.IsEqual(x.StopInfo) {
		return false
	}

	return true
}

// IsEqual returns true if the given Stop is equal to the Stop the method is run on
func (s Stop) IsEqual(x Stop) bool {

	if s.ID != x.ID {
		return false
	}
	if s.Lon != x.Lon {
		return false
	}
	if s.Lat != x.Lat {
		return false
	}
	if s.Name != x.Name {
		return false
	}

	return true
}

// IsEqual returns true if the given StopTimeArray is equal to the StopTimeArray the method is run on
func (st StopTimeArray) IsEqual(x StopTimeArray) bool {

	if len(st) != len(x) {
		return false
	}

	for i := 0; i < len(st); i++ {
		if !st[i].IsEqual(x[i]) {
			return false
		}
	}

	return true
}
