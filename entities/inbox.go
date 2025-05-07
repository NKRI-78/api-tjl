package entities

import "time"

type InboxBadge struct {
	Total int `json:"total"`
}

type InboxListQuery struct {
	Id           string    `json:"id"`
	Title        string    `json:"title"`
	Caption      string    `json:"caption"`
	IsRead       bool      `json:"is_read"`
	UserId       string    `json:"user_id"`
	UserFullname string    `json:"user_fullname"`
	UserAvatar   string    `json:"user_avatar"`
	CreatedAt    time.Time `json:"created_at"`
}

type InboxListResult struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	IsRead    bool      `json:"is_read"`
	User      InboxUser `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

type InboxUser struct {
	Id       string `json:"id"`
	Fullname string `json:"fullname"`
	Avatar   string `json:"avatar"`
}
