package entities

type FormBiodata struct {
	Id        string `json:"id"`
	Place     string `json:"place"`
	Height    string `json:"height"`
	Birthdate string `json:"birthdate"`
	Weight    string `json:"weight"`
	Gender    string `json:"gender"`
	Status    string `json:"status"`
	UserId    string `json:"user_id"`
	Religion  string `json:"religion"`
}

type FormRegion struct {
	ProvinceId    string `json:"province_id"`
	CityId        string `json:"city_id"`
	DistrictId    string `json:"district_id"`
	SubdistrictId string `json:"subdistrict_id"`
	UserId        string `json:"user_id"`
	DetailAddress string `json:"detail_address"`
}

type FormPlace struct {
	Id            string `json:"id"`
	ProvinceId    string `json:"province_id"`
	CityId        string `json:"city_id"`
	DistrictId    string `json:"district_id"`
	SubdistrictId string `json:"subdistrict_id"`
	UserId        string `json:"user_id"`
	DetailAddress string `json:"detail_address"`
}

type FormEducation struct {
	Id              string `json:"id"`
	EducationLevel  string `json:"education_level"`
	Major           string `json:"major"`
	SchoolOrCollege string `json:"school_or_college"`
	StartMonth      string `json:"start_month"`
	StartYear       string `json:"start_year"`
	EndMonth        string `json:"end_month"`
	EndYear         string `json:"end_year"`
	DetailAddress   string `json:"detail_address"`
	UserId          string `json:"user_id"`
}

type FormExercise struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Institution string `json:"institution"`
	StartMonth  string `json:"start_month"`
	StartYear   string `json:"start_year"`
	EndMonth    string `json:"end_month"`
	EndYear     string `json:"end_year"`
	UserId      string `json:"user_id"`
}

type FormWork struct {
	Id          string `json:"id"`
	Position    string `json:"position"`
	Institution string `json:"institution"`
	Work        string `json:"work"`
	IsWork      int    `json:"is_work"`
	StillWork   bool   `json:"still_work"`
	Country     string `json:"country"`
	City        string `json:"city"`
	StartMonth  string `json:"start_month"`
	StartYear   string `json:"start_year"`
	EndMonth    string `json:"end_month"`
	EndYear     string `json:"end_year"`
	UserId      string `json:"user_id"`
}
type FormLanguage struct {
	Id       string `json:"id"`
	Language string `json:"language"`
	Level    string `json:"level"`
	UserId   string `json:"user_id"`
}
