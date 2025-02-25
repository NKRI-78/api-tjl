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

type CheckAccount struct {
	Email string `json:"email"`
}

type CheckJobs struct {
	JobId string `json:"job_id"`
}

type UserLogin struct {
	Uid         string `json:"uid"`
	EmailActive int    `json:"email_active"`
	Password    string `json:"password"`
}

type UserOtp struct {
	Uid         string    `json:"uid"`
	EmailActive int       `json:"email_active"`
	OtpDate     time.Time `json:"otp_date"`
}
