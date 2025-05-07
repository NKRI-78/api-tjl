package entities

import "time"

type CompanyStore struct {
	Id      int    `json:"id"`
	Logo    string `json:"logo"`
	Name    string `json:"name"`
	PlaceId string `json:"place_id"`
}

type CompanyDelete struct {
	Id string `json:"id"`
}

type CompanyUpdate struct {
	Id      string `json:"id"`
	Logo    string `json:"logo"`
	Name    string `json:"name"`
	PlaceId int    `json:"place_id"`
}

type CompanyListQuery struct {
	Id          string    `json:"id"`
	Logo        string    `json:"logo"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	CompanyName string    `json:"company_name"`
}

type CompanyList struct {
	Id        string    `json:"id"`
	Logo      string    `json:"logo"`
	Name      string    `json:"name"`
	Origin    string    `json:"origin"`
	CreatedAt time.Time `json:"created_at"`
}
