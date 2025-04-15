package entities

import "time"

type AllNews struct {
	Id string `json:"id"`
}

type News struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	UserId    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

type NewsResponse struct {
	Id        int         `json:"id"`
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

type NewsUpdate struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Caption string `json:"caption"`
}

type NewsStoreImage struct {
	NewsId string `json:"news_id"`
	Path   string `json:"path"`
}

type NewsUpdateImage struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}

type NewsUser struct {
	Id   string `json:"id"`
	Name string `json:"fullname"`
}

type NewsMedia struct {
	Id   int    `json:"id"`
	Path string `json:"path"`
}
