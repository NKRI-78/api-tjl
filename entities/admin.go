package entities

type ApplicantPerMonthQuery struct {
	Month  string `json:"month"`
	Total  int    `json:"total"`
	Branch string `json:"branch"`
}

type ApplicantPerMonthResponse struct {
	Month  string `json:"month"`
	Branch string `json:"branch"`
	Total  int    `json:"total"`
}

type GenderQuery struct {
	Gender string `json:"gender"`
	Total  int    `json:"total"`
}

type GenderResponse struct {
	Gender string `json:"gender"`
	Total  int    `json:"total"`
}

type CountryQuery struct {
	Country string `json:"country"`
	Total   int    `json:"total"`
}

type CountryResponse struct {
	Country string `json:"country"`
	Total   int    `json:"total"`
}

type ChartSummaryResponse struct {
	ApplicantsPerMonth  []ApplicantPerMonthResponse `json:"applicants_per_month"`
	ApplicantsPerBranch []ApplicantPerMonthResponse `json:"applicants_per_branch"`
	Genders             []GenderResponse            `json:"genders"`
	Countries           []CountryResponse           `json:"countries"`
}
