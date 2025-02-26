package entities

import (
	"time"
)

type User struct {
	Id       string `json:"id"`
	Avatar   string `json:"avatar"`
	Val      string `json:"val"`
	JobId    string `json:"job_id"`
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
	Email    string `json:"email"`
	JobId    string `json:"job_id"`
	JobName  string `json:"job_name"`
}

type ProfileResponse struct {
	Id       string             `json:"id"`
	Avatar   string             `json:"avatar"`
	Phone    string             `json:"phone"`
	Fullname string             `json:"fullname"`
	Email    string             `json:"email"`
	Job      ProfileJobResponse `json:"job"`
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
	Uid     string    `json:"uid"`
	Enabled int       `json:"enabled"`
	OtpDate time.Time `json:"otp_date"`
}
