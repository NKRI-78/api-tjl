package entities

type CompanyStore struct {
	Id   int    `json:"id"`
	Logo string `json:"logo"`
	Name string `json:"name"`
}

type CompanyListQuery struct {
	Id   string `json:"id"`
	Logo string `json:"logo"`
	Name string `json:"name"`
}
