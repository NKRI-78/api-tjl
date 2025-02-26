package entities

type Forum struct {
	Uid     string `json:"uid"`
	Id      string `json:"id"`
	Title   string `json:"title"`
	Caption string `json:"caption"`
}

type ForumCategory struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ForumStore struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Type   int    `json:"type"`
	UserId string `json:"user_id"`
}
