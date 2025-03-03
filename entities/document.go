package entities

type Document struct {
	Id     int    `json:"id"`
	Path   string `json:"path"`
	UserId string `json:"link"`
	Type   int    `json:"type"`
}
