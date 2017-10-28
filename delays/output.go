package delays

import (
	"encoding/json"
	"time"
)

// Models for the json output of the collection service

// OutTrip is the final output for an individual trip running abnormally
type OutTrip struct {
	TripID         string   `json:"trip_id"`
	RouteID        string   `json:"route_id"`
	RouteLongName  string   `json:"route_long_name"`
	RouteShortName string   `json:"route_short_name"`
	NextStop       NextStop `json:"next_stop"`
	VehicleID      string   `json:"vehicle_id"`
	VehicleType    int      `json:"vehicle_type"`
	Lat            float64  `json:"lat"`
	Lon            float64  `json:"lon"`
}

// NextStop is the information for the next stop of an abnormally running service
type NextStop struct {
	StopTimeID       string    `json:"stoptime_id"`
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Lat              float64   `json:"lat"`
	Lon              float64   `json:"lon"`
	ScheduledArrival time.Time `json:"scheduled_arrival"`
	Eta              time.Time `json:"eta"`
	Delay            int       `json:"delay"`
}

// Out is the final output of 1 run of the collection service, ready to be saved into redis
type Out struct {
	Count      int       `json:"count"`
	Trips      []OutTrip `json:"trips"`
	ExecName   string    `json:"exec_name"`
	Created    int64     `json:"created"`
	ValidUntil int64     `json:"valid_until"`
}

// Custom json marshal to modify time for easier usage
func (ns *NextStop) MarshalJSON() ([]byte, error) {
	type Stop NextStop

	sa := ns.ScheduledArrival.Format("15:04:05")
	et := ns.Eta.Format("15:04:05")

	jns := struct {
		*Stop
		SA string `json:"scheduled_arrival"`
		ET string `json:"eta"`
	}{
		Stop: (*Stop)(ns),
		SA:   sa,
		ET:   et,
	}

	return json.Marshal(jns)
}
