package entities

import "time"

type JobListQuery struct {
	Id            string    `json:"id"`
	Title         string    `json:"title"`
	Caption       string    `json:"caption"`
	Salary        string    `json:"salary"`
	CatId         string    `json:"cat_id"`
	CatName       string    `json:"cat_name"`
	PlaceId       int       `json:"place_id"`
	PlaceName     string    `json:"place_name"`
	PlaceCurrency string    `json:"place_currency"`
	CreatedAt     time.Time `json:"created_at"`
}

type JobList struct {
	Id          string      `json:"id"`
	Title       string      `json:"title"`
	Caption     string      `json:"caption"`
	Salary      string      `json:"salary"`
	JobCategory JobCategory `json:"category"`
	JobPlace    JobPlace    `json:"place"`
	Created     time.Time   `json:"created"`
}

type JobStore struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Caption string `json:"caption"`
	Salary  string `json:"salary"`
	CatId   string `json:"cat_id"`
	PlaceId int    `json:"place_id"`
	UserId  string `json:"user_id"`
	IsDraft int    `json:"is_draft"`
}

type JobCategory struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Job struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type JobPlace struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
}
