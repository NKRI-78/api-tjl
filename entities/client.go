package entities

type ClientStore struct {
	Icon string `json:"icon"`
	Link string `json:"link"`
	Name string `json:"name"`
}

type ClientDelete struct {
	Id string `json:"id"`
}

type ClientUpdate struct {
	Id   string `json:"id"`
	Icon string `json:"icon"`
	Link string `json:"link"`
	Name string `json:"name"`
}

type ClientStoreResponse struct {
	Id   int64  `json:"id"`
	Icon string `json:"icon"`
	Link string `json:"link"`
	Name string `json:"name"`
}

type ClientList struct {
	Id   int64  `json:"id"`
	Icon string `json:"icon"`
	Link string `json:"link"`
	Name string `json:"name"`
}
