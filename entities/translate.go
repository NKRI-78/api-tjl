package entities

type TranslateRequest struct {
	Text string `json:"text"`
	To   string `json:"to"`
}

type TranslateResponse struct {
	Text string `json:"text"`
}
