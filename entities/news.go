package entities

import "time"

type AllNews struct {
	Id string `json:"id"`
}

type News struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	UserId    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

type NewsResponse struct {
	Id        string      `json:"id"`
	Title     string      `json:"title"`
	Caption   string      `json:"caption"`
	Media     []NewsMedia `json:"media"`
	User      NewsUser    `json:"user"`
	CreatedAt time.Time   `json:"created_at"`
}

type NewsStore struct {
	Title   string `json:"title"`
	Caption string `json:"caption"`
	UserId  string `json:"user_id"`
}

type NewsUser struct {
	Id   string `json:"id"`
	Name string `json:"fullname"`
}

type NewsMedia struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}
