package entities

type SendFcmRequest struct {
	Title string `json:"to"`
	Body  string `json:"body"`
	Token string `json:"token"`
}
