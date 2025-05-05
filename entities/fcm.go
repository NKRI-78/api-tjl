package entities

type SendFcmRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Token string `json:"token"`
}
