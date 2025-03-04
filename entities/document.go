package entities

type Document struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type DocumentAssign struct {
	UserId string `json:"user_id"`
}

type DocumentStore struct {
	Id     int    `json:"id"`
	Path   string `json:"path"`
	UserId string `json:"user_id"`
	Type   int    `json:"type"`
}
