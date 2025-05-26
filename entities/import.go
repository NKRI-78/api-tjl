package entities

type ImportUserStore struct {
	UserId              string `json:"user_id"`
	Email               string `json:"email"`
	Fullname            string `json:"fullname"`
	Phone               string `json:"phone"`
	Password            string `json:"password"`
	Gender              string `json:"gender"`
	Height              string `json:"height"`
	Weight              string `json:"weight"`
	MaritalStatus       string `json:"marital_status"`
	Religion            string `json:"religion"`
	Place               string `json:"place"`
	Birthdate           string `json:"birthdate"`
	EducationLevel      string `json:"education_level"`
	Major               string `json:"major"`
	SchoolOrCollege     string `json:"school_or_college"`
	StartMonthEducation string `json:"start_month_education"`
	StartYearEducation  string `json:"start_year_education"`
	EndYearEducation    string `json:"end_year_education"`
	EndMonthEducation   string `json:"end_month_education"`
	NameInstitution     string `json:"name_institution"`
	Institution         string `json:"institution"`
	StartYearExercise   string `json:"start_year_exercise"`
	StartMonthExercise  string `json:"start_month_exercise"`
	EndYearExercise     string `json:"end_year_exercise"`
	EndMonthExercise    string `json:"end_month_exercise"`
	PositionWork        string `json:"position_work"`
	InstitutionWork     string `json:"institution_work"`
	Work                string `json:"work"`
	CountryWork         string `json:"country_work"`
	CityWork            string `json:"city_work"`
	StartMonthWork      string `json:"start_month_work"`
	StartYearWork       string `json:"start_year_work"`
	EndMonthWork        string `json:"end_month_work"`
	EndYearWork         string `json:"end_year_work"`
	Level               string `json:"level_language"`
	Language            string `json:"language"`
}
