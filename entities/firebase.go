package entities

type InitFcm struct {
	Token    string `json:"token"`
	Fullname string `json:"fullname"`
	UserId   string `json:"user_id"`
}
