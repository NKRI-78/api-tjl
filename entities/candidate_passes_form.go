package entities

type CandidatePassesForm struct {
	ApplyJobId      string `json:"apply_job_id"`
	UserCandidateId string `json:"user_candidate_id"`
}

type DepartureForm struct {
	DateDeparture   string `json:"date_departure"`
	TimeDeparture   string `json:"time_departure"`
	Airplane        string `json:"airplane"`
	Location        string `json:"location"`
	Destination     string `json:"destination"`
	DepartureId     int    `json:"departure_id"`
	ApplyJobId      string `json:"apply_job_id"`
	UserCandidateId string `json:"user_candidate_id"`
}
