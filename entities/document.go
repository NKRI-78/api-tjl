package entities

type Document struct {
	Id     int    `json:"id"`
	Path   string `json:"path"`
	UserId string `json:"user_id"`
	Type   int    `json:"type"`
}

type DocumentStore struct {
	Id     int    `json:"id"`
	Path   string `json:"path"`
	UserId string `json:"user_id"`
	Type   int    `json:"type"`
}
