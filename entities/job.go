package entities

import "time"

// MYSQL SCAN

type ApplyJobQuery struct {
	Id            int       `json:"id"`
	Uid           string    `json:"uid"`
	UserId        string    `json:"user_id"`
	UserConfirmId string    `json:"user_confirm_id"`
	JobId         string    `json:"job_id"`
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type InfoCandidateApplyJobQuery struct {
	ApplyJobId      string    `json:"apply_job_id"`
	ApplyUserId     string    `json:"apply_user_id"`
	ApplyUserName   string    `json:"apply_user_name"`
	ConfirmUserId   string    `json:"confirm_user_id"`
	ConfirmUserName string    `json:"confirm_user_name"`
	CompanyId       string    `json:"company_id"`
	CompanyLogo     string    `json:"company_logo"`
	CompanyName     string    `json:"company_name"`
	Status          string    `json:"status"`
	Link            string    `json:"link"`
	JobAvatar       string    `json:"job_avatar"`
	JobTitle        string    `json:"job_title"`
	JobCategory     string    `json:"job_category"`
	JobAuthor       string    `json:"job_author"`
	DocId           int       `json:"doc_id"`
	DocName         string    `json:"doc_name"`
	DocPath         string    `json:"doc_path"`
	Schedule        string    `json:"schedule"`
	CreatedAt       time.Time `json:"created_at"`
}

type InfoApplyJobQuery struct {
	ApplyJobId          string    `json:"apply_job_id"`
	ApplyUserId         string    `json:"apply_user_id"`
	ApplyUserName       string    `json:"apply_user_name"`
	ApplyUserEmail      string    `json:"apply_user_email"`
	ApplyUserPhone      string    `json:"apply_user_phone"`
	ConfirmUserId       string    `json:"confirm_user_id"`
	ConfirmUserName     string    `json:"confirm_user_name"`
	CompanyId           string    `json:"company_id"`
	InvitationOffline   string    `json:"invitation_offline"`
	InvitationDeparture string    `json:"invitation_departure"`
	CompanyLogo         string    `json:"company_logo"`
	CompanyName         string    `json:"company_name"`
	CountryName         string    `json:"country_name"`
	Status              string    `json:"status"`
	Link                string    `json:"link"`
	JobAvatar           string    `json:"job_avatar"`
	JobTitle            string    `json:"job_title"`
	JobCategory         string    `json:"job_category"`
	JobAuthor           string    `json:"job_author"`
	DocId               int       `json:"doc_id"`
	DocName             string    `json:"doc_name"`
	DocPath             string    `json:"doc_path"`
	Schedule            string    `json:"schedule"`
	CreatedAt           time.Time `json:"created_at"`
}

type CheckApplyJobQuery struct {
	Id string `json:"uid"`
}

type DocApplyQuery struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

//

type ResultCandidateInfoApplyJob struct {
	Id                  string              `json:"id"`
	Status              string              `json:"status"`
	CreatedAt           time.Time           `json:"created_at"`
	FormFilled          bool                `json:"form_filled"`
	DocFilled           bool                `json:"doc_filled"`
	InvitationOffline   string              `json:"invitation_offline"`
	InvitationDeparture string              `json:"invitation_departure"`
	Job                 JobApply            `json:"job"`
	Company             JobCompanyCandidate `json:"company"`
	UserApply           UserApply           `json:"user_apply"`
	UserConfirm         UserConfirm         `json:"user_confirm"`
}

type ResultCandidateInfoJobDetail struct {
	Id          string      `json:"id"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	Offline     bool        `json:"offline"`
	Schedule    string      `json:"schedule"`
	Link        string      `json:"link"`
	Job         JobApply    `json:"job"`
	Company     JobCompany  `json:"company"`
	Doc         []DocApply  `json:"doc"`
	UserApply   UserApply   `json:"user_apply"`
	UserConfirm UserConfirm `json:"user_confirm"`
}

type ResultInfoApplyJob struct {
	Id             string      `json:"id"`
	Status         string      `json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	Job            JobApply    `json:"job"`
	ReadyDeparture bool        `json:"ready_departure"`
	Company        JobCompany  `json:"company"`
	UserApply      UserApply   `json:"user_apply"`
	UserConfirm    UserConfirm `json:"user_confirm"`
}

type ResultInfoJobDetail struct {
	Id          string      `json:"id"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	Offline     bool        `json:"offline"`
	Schedule    string      `json:"schedule"`
	Link        string      `json:"link"`
	Job         JobApply    `json:"job"`
	Company     JobCompany  `json:"company"`
	Doc         []DocApply  `json:"doc"`
	UserApply   UserApply   `json:"user_apply"`
	UserConfirm UserConfirm `json:"user_confirm"`
}

type DocApply struct {
	DocId   int    `json:"id"`
	DocName string `json:"name"`
	DocPath string `json:"path"`
}

type JobApply struct {
	JobTitle    string `json:"title"`
	JobCategory string `json:"category"`
	JobAvatar   string `json:"logo"`
	JobAuthor   string `json:"author"`
}

type UserConfirm struct {
	Id   string `json:"id"`
	Name string `json:"fullname"`
}

type UserApply struct {
	Id    string `json:"id"`
	Name  string `json:"fullname"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UserApplyCandidate struct {
	Id    string `json:"id"`
	Name  string `json:"fullname"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type InfoApplyJob struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
}

type AllJob struct {
	Uid string `json:"uid"`
}

type ApplyJob struct {
	Id            string `json:"uid"`
	Link          string `json:"link"`
	IsOffline     bool   `json:"is_offline"`
	Content       string `json:"content"`
	Schedule      string `json:"schedule"`
	UserId        string `json:"user_id"`
	UserConfirmId string `json:"user_confirm_id"`
	ApplyJobId    string `json:"apply_job_id"`
	JobId         string `json:"job_id"`
	Status        int    `json:"status"`
}

type AssignDocumentApplyJob struct {
	ApplyJobId string `json:"apply_job_id"`
	DocId      int    `json:"doc_id"`
	Path       string `json:"path"`
}

type AssignBranch struct {
	BranchId string `json:"branch_id"`
	UserId   string `json:"user_id"`
}

type AssignAddress struct {
	ProvinceId    string `json:"province_id"`
	CityId        string `json:"city_id"`
	DistrictId    string `json:"district_id"`
	Subdistrict   string `json:"subdistrict_id"`
	DetailAddress string `json:"detail_address"`
	UserId        string `json:"user_id"`
}

type AssignApplyJob struct {
	JobId  string `json:"job_id"`
	UserId string `json:"user_id"`
}

type JobListQuery struct {
	Id                string    `json:"id"`
	Title             string    `json:"title"`
	Caption           string    `json:"caption"`
	WorkerCount       int       `json:"worker_count"`
	CatId             string    `json:"cat_id"`
	CatIcon           string    `json:"cat_icon"`
	CatName           string    `json:"cat_name"`
	CatType           string    `json:"cat_type"`
	CompanyId         string    `json:"company_id"`
	CompanyLogo       string    `json:"company_logo"`
	CompanyName       string    `json:"company_name"`
	PlaceId           int       `json:"place_id"`
	PlaceName         string    `json:"place_name"`
	PlaceCurrency     string    `json:"place_currency"`
	PlaceKurs         float64   `json:"place_kurs"`
	PlaceInfo         string    `json:"place_info"`
	PlaceSymbol       string    `json:"place_symbol"`
	PlaceLanguageCode string    `json:"place_language_code"`
	Salary            float64   `json:"salary"`
	MinSalary         float64   `json:"min_salary"`
	MaxSalary         float64   `json:"max_salary"`
	SalaryIDR         float64   `json:"salary_id"`
	UserId            string    `json:"user_id"`
	UserAvatar        string    `json:"user_avatar"`
	UserName          string    `json:"user_name"`
	CreatedAt         time.Time `json:"created_at"`
}

type AdminListApplyJobQuery struct {
	Id                  string    `json:"id"`
	Title               string    `json:"title"`
	Caption             string    `json:"caption"`
	Salary              float64   `json:"salary"`
	CatId               string    `json:"cat_id"`
	CatName             string    `json:"cat_name"`
	BranchId            int       `json:"branch_id"`
	BranchName          string    `json:"branch_name"`
	CompanyId           string    `json:"company_id"`
	CompanyLogo         string    `json:"company_logo"`
	CompanyName         string    `json:"company_name"`
	CountryName         string    `json:"country_name"`
	JobStatusId         int       `json:"job_status_id"`
	JobStatusName       string    `json:"job_status_name"`
	PlaceId             int       `json:"place_id"`
	PlaceName           string    `json:"place_name"`
	PlaceCurrency       string    `json:"place_currency"`
	PlaceKurs           float64   `json:"place_kurs"`
	PlaceInfo           string    `json:"place_info"`
	UserIdCandidate     string    `json:"user_id_candidate"`
	UserNameCandidate   string    `json:"user_name_candidate"`
	UserPhoneCandidate  string    `json:"user_phone_candidate"`
	UserId              string    `json:"user_id"`
	UserAvatarCandidate string    `json:"user_avatar_candidate"`
	UserEmailCandidate  string    `json:"user_email_candidate"`
	UserAvatar          string    `json:"user_avatar"`
	UserName            string    `json:"user_name"`
	CreatedAt           time.Time `json:"created_at"`
}

type CandidateExerciseQuery struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Institution string `json:"institution"`
	StartMonth  int    `json:"start_month"`
	StartYear   int    `json:"start_year"`
	EndMonth    int    `json:"end_month"`
	EndYear     int    `json:"end_year"`
}

type CandidateBiodataQuery struct {
	Birthdate string `json:"birthdate"`
	Gender    string `json:"gender"`
	Weight    int    `json:"weight"`
	Height    int    `json:"height"`
	Status    string `json:"status"`
	Religion  string `json:"religion"`
	Place     string `json:"place"`
}

type CandidateLanguageQuery struct {
	Level    string `json:"level"`
	Language string `json:"language"`
}

type CandidateWorkQuery struct {
	Position    string `json:"position"`
	Institution string `json:"institution"`
	Work        string `json:"work"`
	Country     string `json:"country"`
	City        string `json:"city"`
	StartMonth  int    `json:"start_month"`
	StartYear   int    `json:"start_year"`
	EndMonth    int    `json:"end_month"`
	EndYear     int    `json:"end_year"`
	IsWork      bool   `json:"is_work"`
}

type CandidateEducationQuery struct {
	Edu             string `json:"edu"`
	Major           string `json:"major"`
	SchoolOrCollege string `json:"school_or_college"`
	StartMonth      int    `json:"start_month"`
	StartYear       int    `json:"start_year"`
	EndMonth        int    `json:"end_month"`
	EndYear         int    `json:"end_year"`
}

type CandidatePlaceQuery struct {
	ProvinceName    string `json:"province_name"`
	CityName        string `json:"city_name"`
	DistrictName    string `json:"district_name"`
	SubdistrictName string `json:"subdistrict_name"`
	DetailAddress   string `json:"detail_address"`
}

type CandidateDocumentQuery struct {
	Document string `json:"document"`
	Path     string `json:"path"`
}

type AdditionalDocQuery struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

type AdminListApplyJob struct {
	Id          string        `json:"id"`
	Title       string        `json:"title"`
	Caption     string        `json:"caption"`
	Salary      int           `json:"salary"`
	SalaryIDR   string        `json:"salary_idr"`
	Bookmark    bool          `json:"bookmark"`
	Created     string        `json:"created"`
	Company     JobCompany    `json:"company"`
	Candidate   Candidate     `json:"candidate"`
	Status      JobStatus     `json:"status"`
	JobCategory JobCategory   `json:"category"`
	JobPlace    JobPlace      `json:"place"`
	Author      AuthorJobUser `json:"author"`
	Branch      AdminBranch   `json:"branch"`
}

type AdminBranch struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Candidate struct {
	Id                string                         `json:"id"`
	Avatar            string                         `json:"avatar"`
	Email             string                         `json:"email"`
	Name              string                         `json:"name"`
	Phone             string                         `json:"phone"`
	CandidateExercise []CandidateExercise            `json:"exercises"`
	CandidateBiodata  []CandidateBiodata             `json:"biodatas"`
	CandidateLanguage []CandidateLanguage            `json:"languages"`
	CandidateWork     []CandidateWork                `json:"works"`
	CandidatePlace    []CandidatePlace               `json:"places"`
	CandidateEdu      []CandidateEducation           `json:"educations"`
	// CandidateDoc      []CandidateDocument            `json:"documents"`
	AdditionalDoc     []AdditionalDoc                `json:"introductions"`
	Document          map[string][]CandidateDocument `json:"document"`
}

type AdditionalDoc struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

type CandidateExercise struct {
	Name         string                          `json:"name"`
	Institution  string                          `json:"institution"`
	StartMonth   int                             `json:"start_month"`
	StartYear    int                             `json:"start_year"`
	EndMonth     int                             `json:"end_month"`
	EndYear      int                             `json:"end_year"`
	Certificates []CandidateExerciseCertificates `json:"certificates"`
}

type CandidateExerciseCertificates struct {
	Id   int    `json:"id"`
	Path string `json:"path"`
}

type CandidateBiodata struct {
	Birthdate string `json:"birthdate"`
	Gender    string `json:"gender"`
	Weight    int    `json:"weight"`
	Height    int    `json:"height"`
	Status    string `json:"status"`
	Religion  string `json:"religion"`
	Place     string `json:"place"`
}

type CandidateLanguage struct {
	Level    string `json:"level"`
	Language string `json:"language"`
}

type CandidateWork struct {
	Position    string `json:"position"`
	Institution string `json:"institution"`
	Work        string `json:"work"`
	Country     string `json:"country"`
	City        string `json:"city"`
	StartMonth  int    `json:"start_month"`
	StartYear   int    `json:"start_year"`
	EndMonth    int    `json:"end_month"`
	EndYear     int    `json:"end_year"`
	IsWork      bool   `json:"is_work"`
}

type CandidatePlace struct {
	ProvinceName    string `json:"province"`
	CityName        string `json:"city"`
	DistrictName    string `json:"district"`
	SubdistrictName string `json:"subdistrict"`
	DetailAddress   string `json:"detail_address"`
}

type CandidateDocument struct {
	Document string `json:"document"`
	Type     string `json:"type"`
	Path     string `json:"path"`
}

type CandidateEducation struct {
	EducationalLevel string `json:"education_level"`
	Major            string `json:"major"`
	SchoolOrCollege  string `json:"school_or_college"`
	StartMonth       int    `json:"start_month"`
	StartYear        int    `json:"start_year"`
	EndMonth         int    `json:"end_month"`
	EndYear          int    `json:"end_year"`
}

type JobSkillCategory struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type JobStatus struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type JobList struct {
	Id           string             `json:"id"`
	Title        string             `json:"title"`
	Caption      string             `json:"caption"`
	WorkerCount  int                `json:"worker_count"`
	Salary       int                `json:"salary"`
	MinSalary    int                `json:"min_salary"`
	MaxSalary    int                `json:"max_salary"`
	MinSalaryIDR string             `json:"min_salary_idr"`
	MaxSalaryIDR string             `json:"max_salary_idr"`
	SalaryIDR    string             `json:"salary_idr"`
	Bookmark     bool               `json:"bookmark"`
	Created      string             `json:"created"`
	Skills       []JobSkillCategory `json:"skills"`
	Company      JobCompany         `json:"company"`
	JobCategory  JobCategory        `json:"category"`
	JobPlace     JobPlace           `json:"place"`
	Author       AuthorJobUser      `json:"author"`
}

type JobCompany struct {
	Id      string `json:"id"`
	Logo    string `json:"logo"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

type JobCompanyCandidate struct {
	Id      string `json:"id"`
	Logo    string `json:"logo"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

type JobStore struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Skills      []string `json:"skills"`
	Caption     string   `json:"caption"`
	Salary      string   `json:"salary"`
	MinSalary   string   `json:"min_salary"`
	MaxSalary   string   `json:"max_salary"`
	WorkerCount int      `json:"worker_count"`
	CompanyId   string   `json:"company_id"`
	CatId       string   `json:"cat_id"`
	PlaceId     string   `json:"place_id"`
	UserId      string   `json:"user_id"`
	IsDraft     string   `json:"is_draft"`
}

type JobUpdate struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Skills      []string `json:"skills"`
	Caption     string   `json:"caption"`
	Salary      string   `json:"salary"`
	MinSalary   string   `json:"min_salary"`
	MaxSalary   string   `json:"max_salary"`
	WorkerCount int      `json:"worker_count"`
	CompanyId   string   `json:"company_id"`
	CatId       string   `json:"cat_id"`
	PlaceId     string   `json:"place_id"`
	IsDraft     string   `json:"is_draft"`
}

type JobFavourite struct {
	UserId string `json:"user_id"`
	JobId  string `json:"job_id"`
}

type AuthorJobUser struct {
	Id     string `json:"id"`
	Avatar string `json:"avatar"`
	Name   string `json:"fullname"`
}

type JobSkillCategoryList struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type JobSkillCategoryDelete struct {
	JobId  string   `json:"job_id"`
	Skills []string `json:"skills"`
}

type JobSkillCategoryStore struct {
	JobId string `json:"job_id"`
	CatId string `json:"cat_id"`
}

type JobCategoryCount struct {
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Total int    `json:"total"`
}

type JobCategory struct {
	Id   string `json:"id"`
	Icon string `json:"icon"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type JobCategoryStore struct {
	Id   string `json:"id"`
	Icon string `json:"icon"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type JobCategoryDelete struct {
	Id string `json:"id"`
}

type JobCategoryUpdate struct {
	Id   string `json:"id"`
	Icon string `json:"icon"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Job struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	UserId string `json:"user_id"`
}

type JobPlace struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Currency     string `json:"currency"`
	LanguageCode string `json:"language_code"`
	Symbol       string `json:"symbol"`
	Kurs         int    `json:"kurs"`
	Info         string `json:"info"`
}

type TypeJob struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TypeJobStore struct {
	Name string `json:"name"`
}

type TypeJobUpdate struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TypeJobDelete struct {
	Id string `json:"id"`
}
