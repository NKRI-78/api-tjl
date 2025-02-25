package models

import "encoding/json"

func UnmarshalGmapsDirection(data []byte) (GmapsDirection, error) {
	var r GmapsDirection
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GmapsDirection) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GmapsDirection struct {
	GeocodedWaypoints []GeocodedWaypoint `json:"geocoded_waypoints"`
	Routes            []Route            `json:"routes"`
	Status            string             `json:"status"`
}

type GeocodedWaypoint struct {
	GeocoderStatus string   `json:"geocoder_status"`
	PlaceID        string   `json:"place_id"`
	Types          []string `json:"types"`
}

type Route struct {
	Bounds           Bounds        `json:"bounds"`
	Copyrights       string        `json:"copyrights"`
	Legs             []Leg         `json:"legs"`
	OverviewPolyline Polyline      `json:"overview_polyline"`
	Summary          string        `json:"summary"`
	Warnings         []interface{} `json:"warnings"`
	WaypointOrder    []interface{} `json:"waypoint_order"`
}

type Bounds struct {
	Northeast Northeast `json:"northeast"`
	Southwest Northeast `json:"southwest"`
}

type Northeast struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Leg struct {
	Distance          Distance      `json:"distance"`
	Duration          Distance      `json:"duration"`
	EndAddress        string        `json:"end_address"`
	EndLocation       Northeast     `json:"end_location"`
	StartAddress      string        `json:"start_address"`
	StartLocation     Northeast     `json:"start_location"`
	Steps             []Step        `json:"steps"`
	TrafficSpeedEntry []interface{} `json:"traffic_speed_entry"`
	ViaWaypoint       []interface{} `json:"via_waypoint"`
}

type Distance struct {
	Text  string `json:"text"`
	Value int64  `json:"value"`
}

type Step struct {
	Distance         Distance   `json:"distance"`
	Duration         Distance   `json:"duration"`
	EndLocation      Northeast  `json:"end_location"`
	HTMLInstructions string     `json:"html_instructions"`
	Polyline         Polyline   `json:"polyline"`
	StartLocation    Northeast  `json:"start_location"`
	TravelMode       TravelMode `json:"travel_mode"`
	Maneuver         *string    `json:"maneuver,omitempty"`
}

type Polyline struct {
	Points string `json:"points"`
}

type TravelMode string

const (
	Driving TravelMode = "DRIVING"
)