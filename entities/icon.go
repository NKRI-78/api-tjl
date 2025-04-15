package entities

type IconList struct {
	Id   int    `json:"id"`
	Path string `json:"path"`
}

type IconStore struct {
	Path string `json:"path"`
}

type IconUpdate struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}

type IconDelete struct {
	Id string `json:"id"`
}
