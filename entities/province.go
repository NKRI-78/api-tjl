package entities

type Province struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type City struct {
	Id         string `json:"id"`
	ProvinceId string `json:"province_id"`
	Name       string `json:"name"`
}

type District struct {
	Id        string `json:"id"`
	RegencyId string `json:"regency_id"`
	Name      string `json:"name"`
}

type Subdistrict struct {
	Id         string `json:"id"`
	DistrictId string `json:"district_id"`
	Name       string `json:"name"`
}
