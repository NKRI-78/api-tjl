package entities

type JobList struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Desc      string `json:"description"`
	Salary    string `json:"salary"`
	JobName   string `json:"job_name"`
	PlaceName string `json:"place_name"`
}

type JobStore struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Salary  string `json:"salary"`
	CatId   string `json:"cat_id"`
	PlaceId int    `json:"place_id"`
	IsDraft bool   `json:"is_draft"`
}

type JobCategory struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Job struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
