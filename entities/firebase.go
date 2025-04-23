package entities

type InitFcm struct {
	Token    string `json:"token"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	UserId   string `json:"user_id"`
}
