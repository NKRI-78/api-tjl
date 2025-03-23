package entities

import (
	"time"
)

type Biodata struct {
	Personal    ProfileFormBiodata     `json:"personal"`
	Address     ProfileFormPlace       `json:"address"`
	Educations  []ProfileFormEducation `json:"educations"`
	Trainings   []ProfileFormExercise  `json:"trainings"`
	Experiences []ProfileFormWork      `json:"experiences"`
	Languages   []ProfileFormLanguage  `json:"languages"`
}

type CheckAccount struct {
	Email string `json:"email"`
}

type CheckJobs struct {
	JobId string `json:"job_id"`
}

type User struct {
	Uid      string `json:"uid"`
	Id       string `json:"id"`
	Avatar   string `json:"avatar"`
	Val      string `json:"val"`
	JobId    string `json:"job_id"`
	BranchId string `json:"branch_id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	AppName  string `json:"app_name"`
	Otp      string `json:"otp"`
	Password string `json:"password"`
}

type Profile struct {
	Id                 string `json:"id"`
	Avatar             string `json:"avatar"`
	Phone              string `json:"phone"`
	Fullname           string `json:"fullname"`
	Enabled            int    `json:"enabled"`
	Email              string `json:"email"`
	JobId              string `json:"job_id"`
	JobName            string `json:"job_name"`
	BioId              int    `json:"bio_id"`
	BioBirthdate       string `json:"bio_birthdate"`
	BioGender          string `json:"bio_gender"`
	BioHeight          string `json:"bio_height"`
	BioWeight          string `json:"bio_weight"`
	BioStatus          string `json:"bio_status"`
	BioReligion        string `json:"bio_religion"`
	BioPlace           string `json:"bio_place"`
	BioEducationLevel  string `json:"edu_education_level"`
	BioMajor           string `json:"edu_major"`
	BioSchoolOrCollege string `json:"edu_school_or_college"`
	BioStartMonth      string `json:"edu_start_month"`
	BioEndMonth        string `json:"edu_end_month"`
	BioStartYear       string `json:"edu_start_year"`
	BioEndYear         string `json:"edu_end_year"`
	BioName            string `json:"ex_name"`
	BioInstitution     string `json:"ex_institution"`
	BioAddressId       int    `json:"bio_address_id"`
	BioProvinceId      int    `json:"bio_province_id"`
	BioProvince        string `json:"bio_province"`
	BioCityId          int    `json:"bio_city_id"`
	BioCity            string `json:"bio_city"`
	BioDistrictId      int    `json:"bio_district_id"`
	BioDistrict        string `json:"bio_district"`
	BioSubdistrictId   int    `json:"bio_subdistrict_id"`
	BioSubdistrict     string `json:"bio_subdistrict"`
	BioDetailAddress   string `json:"bio_detail_address"`
}

type ProfileResponse struct {
	Id        string             `json:"id"`
	Avatar    string             `json:"avatar"`
	Phone     string             `json:"phone"`
	Fullname  string             `json:"fullname"`
	IsEnabled bool               `json:"enabled"`
	Email     string             `json:"email"`
	Job       ProfileJobResponse `json:"job"`
	Biodata   Biodata            `json:"biodata"`
}

type ProfileFormBiodata struct {
	Id        int    `json:"id"`
	Birthdate string `json:"birthdate"`
	Gender    string `json:"gender"`
	Height    string `json:"height"`
	Weight    string `json:"weight"`
	Religion  string `json:"religion"`
	Place     string `json:"place"`
	Status    string `json:"status"`
}

type ProfileFormPlace struct {
	Id            int                         `json:"id"`
	DetailAddress string                      `json:"detail_address"`
	Province      ProfileFormPlaceData        `json:"province"`
	City          ProfileCityPlaceData        `json:"city"`
	District      ProfileDistrictPlaceData    `json:"district"`
	Subdistrict   ProfileSubdistrictPlaceData `json:"subdistrict"`
}

type ProfileFormPlaceData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ProfileCityPlaceData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ProfileDistrictPlaceData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ProfileSubdistrictPlaceData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ProfileFormEducation struct {
	Id              int    `json:"id"`
	EducationLevel  string `json:"education_level"`
	Major           string `json:"major"`
	SchoolOrCollege string `json:"school_or_college"`
	StartYear       string `json:"start_year"`
	EndYear         string `json:"end_year"`
	StartMonth      string `json:"start_month"`
	EndMonth        string `json:"end_month"`
}

type ProfileFormExercise struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Institution string `json:"institution"`
	StartYear   string `json:"start_year"`
	EndYear     string `json:"end_year"`
	StartMonth  string `json:"start_month"`
	EndMonth    string `json:"end_month"`
}

type ProfileFormWorkQuery struct {
	Id          int    `json:"id"`
	Position    string `json:"position"`
	Institution string `json:"institution"`
	Work        string `json:"work"`
	IsWork      int    `json:"is_work"`
	City        string `json:"city"`
	Country     string `json:"country"`
	StartMonth  string `json:"start_month"`
	StartYear   string `json:"start_year"`
	EndMonth    string `json:"end_month"`
	EndYear     string `json:"end_year"`
}

type ProfileFormWork struct {
	Id          int    `json:"id"`
	Position    string `json:"position"`
	Institution string `json:"institution"`
	Work        string `json:"work"`
	IsWork      bool   `json:"is_work"`
	City        string `json:"city"`
	Country     string `json:"country"`
	StartMonth  string `json:"start_month"`
	StartYear   string `json:"start_year"`
	EndMonth    string `json:"end_month"`
	EndYear     string `json:"end_year"`
}

type ProfileFormLanguage struct {
	Id       int    `json:"id"`
	Level    string `json:"level"`
	Language string `json:"language"`
}

type ProfileJobResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserLogin struct {
	Uid      string `json:"uid"`
	Enabled  int    `json:"enabled"`
	Password string `json:"password"`
}

type UserOtp struct {
	Uid       string    `json:"uid"`
	Enabled   int       `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}
