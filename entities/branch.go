package entities

type Branch struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type StoreBranch struct {
	Name string `json:"name"`
}

type DeleteBranch struct {
	Id int `json:"id"`
}

type UpdateBranch struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
