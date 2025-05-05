package entities

type SendFcmRequest struct {
	Title           string `json:"title"`
	Body            string `json:"body"`
	Token           string `json:"token"`
	JobApplcationId string `json:"job_application_id"`
	BroadcastType   string `json:"broadcast_type"`
}
