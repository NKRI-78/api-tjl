package entities

import "time"

type AllEvent struct {
	Id string `json:"id"`
}

type Event struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	UserId    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

type EventDelete struct {
	Id string `json:"id"`
}

type EventUpdate struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Caption string `json:"caption"`
}

type EventResponse struct {
	Id        int          `json:"id"`
	Title     string       `json:"title"`
	Caption   string       `json:"caption"`
	Media     []EventMedia `json:"media"`
	StartDate string       `json:"start_date"`
	EndDate   string       `json:"end_date"`
	StartTime string       `json:"start_time"`
	EndTime   string       `json:"end_time"`
	User      EventUser    `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
}

type EventStore struct {
	Title     string `json:"title"`
	Caption   string `json:"caption"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	UserId    string `json:"user_id"`
}

type EventStoreImage struct {
	EventId string `json:"event_id"`
	Path    string `json:"path"`
}

type EventUpdateImage struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}

type EventUser struct {
	Id   string `json:"id"`
	Name string `json:"fullname"`
}

type EventMedia struct {
	Id   int    `json:"id"`
	Path string `json:"path"`
}
