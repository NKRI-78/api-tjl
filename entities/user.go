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
	Id       string `json:"id"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	Fullname string `json:"fullname"`
	Enabled  int    `json:"enabled"`
	Email    string `json:"email"`
	JobId    string `json:"job_id"`
	JobName  string `json:"job_name"`
}

type ProfileResponse struct {
	Id        string             `json:"id"`
	Avatar    string             `json:"avatar"`
	Phone     string             `json:"phone"`
	Fullname  string             `json:"fullname"`
	IsEnabled bool               `json:"enabled"`
	Email     string             `json:"email"`
	Job       ProfileJobResponse `json:"job"`
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
