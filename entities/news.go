package entities

import "time"

type AllNews struct {
	Id string `json:"id"`
}

type News struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	UserId    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

type NewsResponse struct {
	Id        string      `json:"id"`
	Title     string      `json:"title"`
	Desc      string      `json:"desc"`
	Media     []NewsMedia `json:"media"`
	User      NewsUser    `json:"user"`
	CreatedAt time.Time   `json:"created_at"`
}

type NewsUser struct {
	Id   string `json:"id"`
	Name string `json:"fullname"`
}

type NewsMedia struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}
