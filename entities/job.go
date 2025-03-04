package entities

type JobList struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Salary string `json:"salary"`
	Place string 
}

type JobCategory struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Job struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
