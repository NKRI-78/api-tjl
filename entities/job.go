package entities

import "time"

type ApplyJob struct {
	UserId string `json:"user_id"`
	JobId  string `json:"job_id"`
	Status string `json:"status"`
}

type JobListQuery struct {
	Id            string    `json:"id"`
	Title         string    `json:"title"`
	Caption       string    `json:"caption"`
	CatId         string    `json:"cat_id"`
	CatName       string    `json:"cat_name"`
	PlaceId       int       `json:"place_id"`
	PlaceName     string    `json:"place_name"`
	PlaceCurrency string    `json:"place_currency"`
	PlaceKurs     float64   `json:"place_kurs"`
	PlaceInfo     string    `json:"place_info"`
	Salary        float64   `json:"salary"`
	SalaryIDR     float64   `json:"salary_id"`
	UserId        string    `json:"user_id"`
	UserAvatar    string    `json:"user_avatar"`
	UserName      string    `json:"user_name"`
	CreatedAt     time.Time `json:"created_at"`
}

type JobListAdminQuery struct {
	Id            string    `json:"id"`
	Title         string    `json:"title"`
	Caption       string    `json:"caption"`
	Salary        float64   `json:"salary"`
	CatId         string    `json:"cat_id"`
	CatName       string    `json:"cat_name"`
	JobStatusId   int       `json:"job_status_id"`
	JobStatusName string    `json:"job_status_name"`
	PlaceId       int       `json:"place_id"`
	PlaceName     string    `json:"place_name"`
	PlaceCurrency string    `json:"place_currency"`
	PlaceKurs     float64   `json:"place_kurs"`
	PlaceInfo     string    `json:"place_info"`
	UserId        string    `json:"user_id"`
	UserAvatar    string    `json:"user_avatar"`
	UserName      string    `json:"user_name"`
	CreatedAt     time.Time `json:"created_at"`
}

type JobListAdmin struct {
	Id          string      `json:"id"`
	Title       string      `json:"title"`
	Caption     string      `json:"caption"`
	Salary      int         `json:"salary"`
	SalaryIDR   string      `json:"salary_idr"`
	Bookmark    bool        `json:"bookmark"`
	Created     string      `json:"created"`
	Status      JobStatus   `json:"status"`
	JobCategory JobCategory `json:"category"`
	JobPlace    JobPlace    `json:"place"`
	JobUser     JobUser     `json:"user"`
}

type JobStatus struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type JobList struct {
	Id          string      `json:"id"`
	Title       string      `json:"title"`
	Caption     string      `json:"caption"`
	Salary      int         `json:"salary"`
	SalaryIDR   string      `json:"salary_idr"`
	Bookmark    bool        `json:"bookmark"`
	Created     string      `json:"created"`
	JobCategory JobCategory `json:"category"`
	JobPlace    JobPlace    `json:"place"`
	JobUser     JobUser     `json:"user"`
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

type JobFavourite struct {
	UserId string `json:"user_id"`
	JobId  string `json:"job_id"`
}

type JobUser struct {
	Id     string `json:"id"`
	Avatar string `json:"avatar"`
	Name   string `json:"fullname"`
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
	Kurs     int    `json:"kurs"`
	Info     string `json:"info"`
}
