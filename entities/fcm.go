package entities

type SendFcmRequest struct {
	Title         string `json:"title"`
	Body          string `json:"body"`
	Token         string `json:"token"`
	JobAppId      string `json:"job_app_id"`
	BroadcastType string `json:"broadcast_type"`
}
