package entities

type Province struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type City struct {
	Id         int    `json:"id"`
	ProvinceId string `json:"province_id"`
	Name       string `json:"name"`
}

type District struct {
	Id        int    `json:"id"`
	RegencyId string `json:"regency_id"`
	Name      string `json:"name"`
}

type Subdistrict struct {
	Id         int    `json:"id"`
	DistrictId string `json:"district_id"`
	Name       string `json:"name"`
}

type Country struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CountryStore struct {
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Kurs     int    `json:"kurs"`
	Info     string `json:"info"`
}

type CountryUpdate struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Kurs     string `json:"kurs"`
	Info     string `json:"info"`
}
