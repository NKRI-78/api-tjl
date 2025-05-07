package entities

import "time"

type CandidatePassesForm struct {
	ApplyJobId      string `json:"apply_job_id"`
	UserCandidateId string `json:"user_candidate_id"`
}

type CandidatePassesFormListQuery struct {
	Id            int       `json:"id"`
	DateDeparture string    `json:"date_departure"`
	TimeDeparture string    `json:"time_departure"`
	Airplane      string    `json:"airplane"`
	Location      string    `json:"location"`
	Destination   string    `json:"destination"`
	UserId        string    `json:"user_id"`
	UserAvatar    string    `json:"avatar"`
	UserFullname  string    `json:"fullname"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CandidatePassesFormListResult struct {
	Id            int                     `json:"id"`
	DateDeparture string                  `json:"date_departure"`
	TimeDeparture string                  `json:"time_departure"`
	Airplane      string                  `json:"airplane"`
	Location      string                  `json:"location"`
	Destination   string                  `json:"destination"`
	User          CandidatePassesFormUser `json:"user"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

type CandidatePassesFormUser struct {
	Id       string `json:"id"`
	Avatar   string `json:"avatar"`
	Fullname string `json:"fullname"`
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
