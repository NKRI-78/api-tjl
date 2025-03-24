package entities

import "time"

// MYSQL SCAN

type ApplyJobQuery struct {
	Id            int       `json:"id"`
	Uid           string    `json:"uid"`
	UserId        string    `json:"user_id"`
	UserConfirmId string    `json:"user_confirm_id"`
	JobId         string    `json:"job_id"`
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type InfoApplyJobQuery struct {
	ApplyJobId      string    `json:"apply_job_id"`
	ApplyUserId     string    `json:"apply_user_id"`
	ApplyUserName   string    `json:"apply_user_name"`
	ConfirmUserId   string    `json:"confirm_user_id"`
	ConfirmUserName string    `json:"confirm_user_name"`
	Status          string    `json:"status"`
	Link            string    `json:"link"`
	JobAvatar       string    `json:"job_avatar"`
	JobTitle        string    `json:"job_title"`
	JobCategory     string    `json:"job_category"`
	JobAuthor       string    `json:"job_author"`
	DocId           int       `json:"doc_id"`
	DocName         string    `json:"doc_name"`
	DocPath         string    `json:"doc_path"`
	Schedule        string    `json:"schedule"`
	CreatedAt       time.Time `json:"created_at"`
}

type CheckApplyJobQuery struct {
	Id string `json:"uid"`
}

type DocApplyQuery struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

//

type ResultInfoJob struct {
	Id          string      `json:"id"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	Schedule    string      `json:"schedule"`
	Link        string      `json:"link"`
	Job         JobApply    `json:"job"`
	UserApply   UserApply   `json:"user_apply"`
	UserConfirm UserConfirm `json:"user_confirm"`
}

type ResultInfoJobDetail struct {
	Id          string      `json:"id"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	Schedule    string      `json:"schedule"`
	Link        string      `json:"link"`
	Job         JobApply    `json:"job"`
	Doc         []DocApply  `json:"doc"`
	UserApply   UserApply   `json:"user_apply"`
	UserConfirm UserConfirm `json:"user_confirm"`
}

type DocApply struct {
	DocId   int    `json:"id"`
	DocName string `json:"name"`
	DocPath string `json:"path"`
}

type JobApply struct {
	JobTitle    string `json:"title"`
	JobCategory string `json:"category"`
	JobAvatar   string `json:"logo"`
	JobAuthor   string `json:"author"`
}

type UserConfirm struct {
	Id   string `json:"id"`
	Name string `json:"fullname"`
}

type UserApply struct {
	Id   string `json:"id"`
	Name string `json:"fullname"`
}

type InfoApplyJob struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
}

type ApplyJob struct {
	Id            string `json:"uid"`
	Link          string `json:"link"`
	Schedule      string `json:"schedule"`
	UserId        string `json:"user_id"`
	UserConfirmId string `json:"user_confirm_id"`
	ApplyJobId    string `json:"apply_job_id"`
	JobId         string `json:"job_id"`
	Status        int    `json:"status"`
}

type AssignDocumentApplyJob struct {
	ApplyJobId string `json:"apply_job_id"`
	DocId      int    `json:"doc_id"`
	Path       string `json:"path"`
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

type JobCategoryCount struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
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
