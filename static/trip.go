package static

import (
	"database/sql"
	"fmt"
)

// Trip represents a trip as stored in the database
type Trip struct {
	ID        string
	RouteID   string
	ServiceID string
	GTFSID    string
	Headsign  string
}

type Trips []Trip

// TripStore defines methods that a concrete trip service should implement
type TripStore interface {
	GetTripByGTFSID(id string) (Trip, error)
	GetTripByRouteID(id string) (Trips, error)
}

// TripService implements TripStore for psql
type TripService struct {
	DB *sql.DB
}

// TripServiceInit initializes and returns a TripService with a given sql db connector
func TripServiceInit(db *sql.DB) *TripService {
	return &TripService{DB: db}
}

// GetTripByGTFSID returns a trip with the given realtime trip id or an error
func (ts *TripService) GetTripByGTFSID(id string) (Trip, error) {
	t := Trip{}

	row := ts.DB.QueryRow("SELECT trip_id, route_id, service_id, gtfs_trip_id, trip_headsign FROM trips WHERE gtfs_trip_id = $1", id)
	err := row.Scan(&t.ID, &t.RouteID, &t.ServiceID, &t.GTFSID, &t.Headsign)

	if err != nil {
		return t, err
	}

	return t, nil
}

func (ts *TripService) GetTripByRouteID(id string) (Trips, error) {
	tps := Trips{}
	var ri string

	row := ts.DB.QueryRow("SELECT route_id FROM routes WHERE gtfs_route_id = $1", id)
	err := row.Scan(&ri)

	if err != nil {
		return tps, err
	}

	rows, err := ts.DB.Query("SELECT trip_id, route_id, service_id, gtfs_trip_id, trip_headsign FROM trips WHERE route_id = $1", ri)
	if err != nil {
		return tps, fmt.Errorf("trip - GetTripFromRouteID: %v", err)
	}

	for rows.Next() {
		t := Trip{}
		err = row.Scan(&t.ID, &t.RouteID, &t.ServiceID, &t.GTFSID, &t.Headsign)
		if err != nil {
			return t, fmt.Errorf("trip - GetTripFromRouteID: Failed to scan: %v", err)
		}
		tps = append(tps, t)

	}

	return tps, nil
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
