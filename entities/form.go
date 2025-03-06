package entities

type FormBiodata struct {
	Id       string `json:"id"`
	Place    string `json:"place"`
	Height   string `json:"height"`
	Weight   string `json:"weight"`
	Gender   string `json:"gender"`
	Status   string `json:"status"`
	Religion string `json:"religion"`
}

type FormRegion struct {
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	Subdistrict   string `json:"subdistrict"`
	DetailAddress string `json:"detail_address"`
}

type FormEducation struct {
	EducationLevel  string `json:"education_level"`
	Major           string `json:"major"`
	SchoolOrCollege string `json:"school_or_college"`
	StartMonth      string `json:"start_month"`
	StartYear       string `json:"start_year"`
	EndMonth        string `json:"end_month"`
	EndYear         string `json:"end_year"`
}

type FormExercise struct {
	Name        string `json:"name"`
	Institution string `json:"institution"`
	StartMonth  string `json:"start_month"`
	StartYear   string `json:"start_year"`
	EndMonth    string `json:"end_month"`
	EndYear     string `json:"end_year"`
}

type FormWork struct {
	Name        string `json:"name"`
	Institution string `json:"institution"`
	StartMonth  string `json:"start_month"`
	StartYear   string `json:"start_year"`
	EndMonth    string `json:"end_month"`
	EndYear     string `json:"end_year"`
}
