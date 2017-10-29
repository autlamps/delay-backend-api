package static

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Trip represents a trip as stored in the database
type Trip struct {
	ID        string   `json:"id"`
	RouteID   string   `json:"route_id"`
	ServiceID string   `json:"service_id"`
	GTFSID    string   `json:"gtfsid"`
	Headsign  string   `json:"headsign"`
	Calendar  Calendar `json:"calendar"`
}

type Calendar struct {
	Start     time.Time `json:"-"`
	End       time.Time `json:"-"`
	Monday    bool      `json:"monday"`
	Tuesday   bool      `json:"tuesday"`
	Wednesday bool      `json:"wednesday"`
	Thursday  bool      `json:"thursday"`
	Friday    bool      `json:"friday"`
	Saturday  bool      `json:"saturday"`
	Sunday    bool      `json:"sunday"`
}

type Trips []Trip

// TripStore defines methods that a concrete trip service should implement
type TripStore interface {
	GetTripByGTFSID(id string) (Trip, error)
	GetTripsByRouteID(id string) (Trips, error)
}

// TripService implements TripStore for psql
type TripService struct {
	db *sql.DB
}

// TripServiceInit initializes and returns a TripService with a given sql db connector
func TripServiceInit(db *sql.DB) *TripService {
	return &TripService{db: db}
}

// GetTripByGTFSID returns a trip with the given realtime trip id or an error
func (ts *TripService) GetTripByGTFSID(id string) (Trip, error) {
	t := Trip{}

	row := ts.db.QueryRow("SELECT trip_id, route_id, service_id, gtfs_trip_id, trip_headsign FROM trips WHERE gtfs_trip_id = $1", id)
	err := row.Scan(&t.ID, &t.RouteID, &t.ServiceID, &t.GTFSID, &t.Headsign)

	if err != nil {
		return t, err
	}

	return t, nil
}

func (ts *TripService) GetTripsByRouteID(id string) (Trips, error) {
	tps := Trips{}

	rows, err := ts.db.Query("SELECT trip_id, route_id, service_id, gtfs_trip_id, trip_headsign FROM trips WHERE route_id = $1", id)

	if err != nil {
		return tps, fmt.Errorf("trip - GetTripFromRouteID: %v", err)
	}

	for rows.Next() {
		t := Trip{}

		err = rows.Scan(&t.ID, &t.RouteID, &t.ServiceID, &t.GTFSID, &t.Headsign)

		if err != nil {
			return Trips{}, fmt.Errorf("trip - GetTripFromRouteID: Failed to scan: %v", err)
		}

		t.Calendar, err = ts.getCalendarForTrip(t)

		if err != nil {
			return Trips{}, fmt.Errorf("trip - GetTripsFromRouteID: failed to get cal: %v", err)
		}

		tps = append(tps, t)
	}

	return tps, nil
}

func (ts *TripService) getCalendarForTrip(t Trip) (Calendar, error) {
	cal := Calendar{}

	row := ts.db.QueryRow("SELECT monday, tuesday, wednesday, thursday, friday, saturday, sunday FROM calendar WHERE service_id = $1", t.ServiceID)

	err := row.Scan(&cal.Monday, &cal.Tuesday, &cal.Wednesday, &cal.Thursday, &cal.Friday, &cal.Saturday, &cal.Sunday)

	if err != nil {
		return Calendar{}, fmt.Errorf("trip - getCalendarForTrip: Failed to read in days from db: %v", err)
	}

	arrivalTimes := []time.Time{}

	rows, err := ts.db.Query("SELECT arrival_time FROM stop_times WHERE trip_id = $1 ORDER BY stop_sequence ASC", t.ID)

	for rows.Next() {
		var arrival time.Time

		err = rows.Scan(&arrival)

		if err != nil {
			return Calendar{}, fmt.Errorf("trip - getCalendarForTrip: Failed to read in times from db: %v", err)
		}

		arrival = arrival.Local()

		arrivalTimes = append(arrivalTimes, arrival)
	}

	cal.Start = arrivalTimes[0]                 // First arrival time of the trip at a stop, the "start time"
	cal.End = arrivalTimes[len(arrivalTimes)-1] // Last arrival of a trip at a stop the "end time"

	return cal, nil
}

// IsEqual returns true if the given trip is equal to the trip the method is being called on
func (t Trip) IsEqual(x Trip) bool {

	if t.ID != x.ID {
		return false
	}

	if t.RouteID != x.RouteID {
		return false
	}

	if t.ServiceID != x.ServiceID {
		return false
	}

	if t.GTFSID != x.GTFSID {
		return false
	}

	if t.Headsign != x.Headsign {
		return false
	}

	return true
}

func (c *Calendar) MarshalJSON() ([]byte, error) {
	type Cal Calendar

	js := struct {
		*Cal
		Start string `json:"start_time"`
		End   string `json:"end"`
	}{
		Cal:   (*Cal)(c),
		Start: c.Start.Format("15:04:05"),
		End:   c.End.Format("15:04:05"),
	}

	return json.Marshal(js)
}
