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

type GetDocumentAdditional struct {
	Id     int    `json:"id"`
	UserId string `json:"user_id"`
	Path   string `json:"path"`
	Type   string `json:"type"`
}

type GetDocumentAdditionalResponse struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
	Path string `json:"path"`
}

type DocumentAdditionalStore struct {
	Id     int    `json:"id"`
	Path   string `json:"path"`
	UserId string `json:"user_id"`
	Type   string `json:"type"`
}

type DocumentAdditionalUpdate struct {
	Id     int    `json:"id"`
	Path   string `json:"path"`
	UserId string `json:"user_id"`
	Type   string `json:"type"`
}
