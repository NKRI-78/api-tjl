package entities

import (
	"time"
)

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
	BioExStartYear     string `json:"ex_start_year"`
	BioExEndYear       string `json:"ex_end_year"`
	BioExStartMonth    string `json:"ex_start_month"`
	BioExEndMonth      string `json:"ex_end_month"`
}

type ProfileResponse struct {
	Id            string               `json:"id"`
	Avatar        string               `json:"avatar"`
	Phone         string               `json:"phone"`
	Fullname      string               `json:"fullname"`
	IsEnabled     bool                 `json:"enabled"`
	Email         string               `json:"email"`
	Job           ProfileJobResponse   `json:"job"`
	FormBiodata   ProfileFormBiodata   `json:"form_biodata"`
	FormEducation ProfileFormEducation `json:"form_education"`
	FormExercise  ProfileFormExercise  `json:"form_exercise"`
}

type ProfileFormBiodata struct {
	Birthdate string `json:"birthdate"`
	Gender    string `json:"gender"`
	Height    string `json:"height"`
	Weight    string `json:"weight"`
	Religion  string `json:"religion"`
	Place     string `json:"place"`
	Status    string `json:"status"`
}

type ProfileFormEducation struct {
	EducationLevel  string `json:"education_level"`
	Major           string `json:"major"`
	SchoolOrCollege string `json:"school_or_college"`
	StartYear       string `json:"start_year"`
	EndYear         string `json:"end_year"`
	StartMonth      string `json:"start_month"`
	EndMonth        string `json:"end_month"`
}

type ProfileFormExercise struct {
	Name        string `json:"name"`
	Institution string `json:"institution"`
	StartYear   string `json:"start_year"`
	EndYear     string `json:"end_year"`
	StartMonth  string `json:"start_month"`
	EndMonth    string `json:"end_month"`
}

type ProfileJobResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CheckAccount struct {
	Email string `json:"email"`
}

type CheckJobs struct {
	JobId string `json:"job_id"`
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
