package services

import (
	"database/sql"
	"errors"
	"fmt"
	entities "superapps/entities"
	helper "superapps/helpers"
	models "superapps/models"

	uuid "github.com/satori/go.uuid"
)

func ListInfoApplyJob(iaj *models.InfoApplyJob) (map[string]any, error) {

	var dataQuery entities.InfoApplyJobQuery
	var data []entities.ResultInfoJob

	query := `SELECT paa.user_id AS apply_user_id, paa.fullname AS apply_user_name, 
		pac.user_id AS confirm_user_id, pac.fullname AS confirm_user_name,
		js.name AS status, aj.uid AS apply_job_id,
		j.title AS job_title,
		jc.name AS job_category,
		p.avatar AS job_avatar,
		p.fullname AS job_author,
		aj.created_at
		FROM apply_jobs aj 
		INNER JOIN jobs j ON j.uid = aj.job_id
		INNER JOIN job_categories jc ON jc.uid = j.cat_id
		INNER JOIN profiles p ON p.user_id = j.user_id
		INNER JOIN job_statuses js ON js.id = aj.status
		INNER JOIN profiles paa ON paa.user_id = aj.user_id
		LEFT JOIN profiles pac ON pac.user_id = aj.user_confirm_id 
		WHERE aj.user_id = ?
		ORDER BY aj.created_at DESC
	`
	rows, err := db.Debug().Raw(query, iaj.UserId).Rows()

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		errJobRows := db.ScanRows(rows, &dataQuery)

		if errJobRows != nil {
			helper.Logger("error", "In Server: "+errJobRows.Error())
			return nil, errors.New(errJobRows.Error())
		}

		data = append(data, entities.ResultInfoJob{
			Id:        dataQuery.ApplyJobId,
			Status:    dataQuery.Status,
			CreatedAt: dataQuery.CreatedAt,
			Job: entities.JobApply{
				JobTitle:    dataQuery.JobTitle,
				JobCategory: dataQuery.JobCategory,
				JobAvatar:   dataQuery.JobAvatar,
				JobAuthor:   dataQuery.JobAuthor,
			},
			UserApply: entities.UserApply{
				Id:   dataQuery.ApplyUserId,
				Name: dataQuery.ApplyUserName,
			},
			UserConfirm: entities.UserConfirm{
				Id:   helper.DefaultIfEmpty(dataQuery.ConfirmUserId, "-"),
				Name: helper.DefaultIfEmpty(dataQuery.ConfirmUserName, "-"),
			},
		})
	}

	if data == nil {
		data = []entities.ResultInfoJob{}
	}

	return map[string]any{
		"data": data,
	}, nil
}

func InfoApplyJob(iaj *models.InfoApplyJob) (map[string]any, error) {

	var dataQuery entities.InfoApplyJobQuery
	var dataDocQuery entities.DocApplyQuery
	var data []entities.ResultInfoJobDetail
	var dataDoc []entities.DocApply

	query := `SELECT paa.user_id AS apply_user_id, paa.fullname AS apply_user_name, 
		pac.user_id AS confirm_user_id, pac.fullname AS confirm_user_name,
		js.name AS status, aj.created_at, aj.uid AS apply_job_id, aj.link, aj.schedule,
		j.title AS job_title,
		jc.name AS job_category,
		p.avatar AS job_avatar,
		p.fullname AS job_author
		FROM apply_job_histories aj 
		INNER JOIN jobs j ON j.uid = aj.job_id
		INNER JOIN job_categories jc ON jc.uid = j.cat_id
		INNER JOIN profiles p ON p.user_id = j.user_id
		INNER JOIN job_statuses js ON js.id = aj.status
		INNER JOIN profiles paa ON paa.user_id = aj.user_id
		LEFT JOIN profiles pac ON pac.user_id = aj.user_confirm_id 
		WHERE aj.uid = ?
	`
	rows, err := db.Debug().Raw(query, iaj.Id).Rows()

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		errJobRows := db.ScanRows(rows, &dataQuery)

		if errJobRows != nil {
			helper.Logger("error", "In Server: "+errJobRows.Error())
			return nil, errors.New(errJobRows.Error())
		}

		queryDoc := `SELECT 
			d.id, 
			d.name, 
		  COALESCE(ajd.path, '-') AS path 
		FROM 
			documents d
		LEFT JOIN 
			apply_job_documents ajd 
			ON ajd.doc_id = d.id 
			AND ajd.apply_job_id = ?`

		rowsDoc, errDoc := db.Debug().Raw(queryDoc, dataQuery.ApplyJobId).Rows()

		if errDoc != nil {
			helper.Logger("error", "In Server: "+errDoc.Error())
		}
		defer rowsDoc.Close()

		dataDoc = []entities.DocApply{}

		for rowsDoc.Next() {
			errDocRows := db.ScanRows(rowsDoc, &dataDocQuery)

			if errDocRows != nil {
				helper.Logger("error", "In Server: "+errDocRows.Error())
				return nil, errors.New(errDocRows.Error())
			}

			dataDoc = append(dataDoc, entities.DocApply{
				DocId:   dataDocQuery.Id,
				DocName: dataDocQuery.Name,
				DocPath: dataDocQuery.Path,
			})
		}

		data = append(data, entities.ResultInfoJobDetail{
			Id:        dataQuery.ApplyJobId,
			Status:    dataQuery.Status,
			Doc:       dataDoc,
			CreatedAt: dataQuery.CreatedAt,
			Link:      helper.DefaultIfEmpty(dataQuery.Link, "-"),
			Schedule:  helper.DefaultIfEmpty(dataQuery.Schedule, "-"),
			Job: entities.JobApply{
				JobTitle:    dataQuery.JobTitle,
				JobCategory: dataQuery.JobCategory,
				JobAvatar:   dataQuery.JobAvatar,
				JobAuthor:   dataQuery.JobAuthor,
			},
			UserApply: entities.UserApply{
				Id:   dataQuery.ApplyUserId,
				Name: dataQuery.ApplyUserName,
			},
			UserConfirm: entities.UserConfirm{
				Id:   helper.DefaultIfEmpty(dataQuery.ConfirmUserId, "-"),
				Name: helper.DefaultIfEmpty(dataQuery.ConfirmUserName, "-"),
			},
		})
	}

	return map[string]any{
		"data": data,
	}, nil
}

func AssignDocumentApplyJob(adaj *models.AssignDocumentApplyJob) (map[string]any, error) {
	query := `
		INSERT INTO apply_job_documents (apply_job_id, doc_id, path) 
		VALUES (?, ?, ?) 
		ON DUPLICATE KEY UPDATE path = VALUES(path)
	`

	err := db.Debug().Exec(query, adaj.ApplyJobId, adaj.DocId, adaj.Path).Error
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func ApplyJob(aj *models.ApplyJob) (map[string]any, error) {
	var dataUserFcm entities.InitFcm
	var allJob []models.CheckApplyJobQuery

	queryCheck := `SELECT uid FROM apply_jobs 
	WHERE user_id = ? 
	AND job_id = ? 
	AND status = 3`

	errAllJob := db.Debug().Raw(queryCheck, aj.UserId, aj.JobId).Scan(&allJob).Error

	if errAllJob != nil {
		helper.Logger("error", "In Server: "+errAllJob.Error())
		return nil, errors.New(errAllJob.Error())
	}

	var isUserAppliedJob = len(allJob)

	if isUserAppliedJob == 1 {
		helper.Logger("error", "In Server: user already applied job")
		return nil, errors.New("USER_ALREADY_APPLIED_JOB")
	}

	queryInsert := `INSERT INTO apply_jobs (uid, job_id, user_id) VALUES (?, ?, ?)`

	errInsert := db.Debug().Exec(queryInsert, aj.Id, aj.JobId, aj.UserId).Error

	if errInsert != nil {
		helper.Logger("error", "In Server: "+errInsert.Error())
		return nil, errors.New(errInsert.Error())
	}

	queryHistory := `INSERT INTO apply_job_histories (uid, job_id, user_id) VALUES (?, ?, ?)`

	errHistory := db.Debug().Exec(queryHistory, aj.Id, aj.JobId, aj.UserId).Error

	if errHistory != nil {
		helper.Logger("error", "In Server: "+errHistory.Error())
		return nil, errors.New(errHistory.Error())
	}

	queryUserFcm := `SELECT f.token, p.fullname FROM fcms f 
	INNER JOIN profiles p ON p.user_id = f.user_id 
	WHERE f.user_id = ?`

	rowUserFcm := db.Debug().Raw(queryUserFcm, aj.UserId).Row()

	errUserFcmRow := rowUserFcm.Scan(&dataUserFcm.Token, &dataUserFcm.Fullname)

	if errUserFcmRow != nil {
		helper.Logger("error", "In Server: "+errUserFcmRow.Error())
	}

	message := fmt.Sprintf("Silahkan menunggu untuk tahap selanjutnya [%s]", dataUserFcm.Fullname)
	helper.SendFcm("Selamat Anda telah berhasil melamar", message, dataUserFcm.Token)

	return map[string]any{}, nil
}

func UpdateApplyJob(uaj *models.ApplyJob) (map[string]any, error) {
	var dataUserFcm entities.InitFcm
	var dataQuery entities.ApplyJobQuery

	// Fetch existing job application details
	queryInfo := `SELECT uid, job_id, user_id, status FROM apply_jobs WHERE uid = ?`
	row := db.Debug().Raw(queryInfo, uaj.ApplyJobId).Row()
	errJobRow := row.Scan(&dataQuery.Uid, &dataQuery.JobId, &dataQuery.UserId, &dataQuery.Status)

	if errJobRow != nil {
		helper.Logger("error", "In Server: "+errJobRow.Error())
		return nil, errors.New("data not found")
	}

	var status string

	// Validate status transition rules
	switch dataQuery.Status {
	case 1: // IN_PROGRESS -> can only move to INTERVIEW (2)
		if uaj.Status != 2 {
			helper.Logger("error", "status IN_PROGRESS can only move to INTERVIEW")
			return nil, errors.New("status IN_PROGRESS can only move to INTERVIEW")
		}
		status = "INTERVIEW"
	case 2: // INTERVIEW -> can move to ACCEPTED (3) or DECLINED (4)
		if uaj.Status != 3 && uaj.Status != 4 {
			helper.Logger("error", "status INTERVIEW can only move to ACCEPTED or DECLINED")
			return nil, errors.New("status INTERVIEW can only move to ACCEPTED or DECLINED")
		}
		if uaj.Status == 3 {
			status = "ACCEPTED"
		} else {
			status = "DECLINED"
		}
	case 3: // ACCEPTED - no further updates
		helper.Logger("error", "status [ACCEPTED] already passed")
		return nil, errors.New("status [ACCEPTED] already passed")
	case 4: // DECLINED - no further updates
		helper.Logger("error", "status [DECLINED] already passed")
		return nil, errors.New("status [DECLINED] already passed")
	default:
		helper.Logger("error", "unknown status")
		return nil, errors.New("unknown status")
	}

	queryUserFcm := `SELECT f.token, p.fullname FROM fcms f 
	INNER JOIN profiles p ON p.user_id = f.user_id 
	WHERE f.user_id = ?`

	rowUserFcm := db.Debug().Raw(queryUserFcm, uaj.UserId).Row()

	errUserFcmRow := rowUserFcm.Scan(&dataUserFcm.Token, &dataUserFcm.Fullname)

	if errUserFcmRow != nil {
		if errors.Is(errUserFcmRow, sql.ErrNoRows) {
			helper.Logger("info", "No FCM data found for user")
		}

		helper.Logger("error", "In Server: "+errUserFcmRow.Error())
	}

	if dataUserFcm.Token != "" {
		title := fmt.Sprintf("Selamat lamaran Anda sudah dalam tahap [%s]", status)
		helper.SendFcm(title, dataUserFcm.Fullname, dataUserFcm.Token)
	}

	// Perform the update
	query := `UPDATE apply_jobs SET user_confirm_id = ?, status = ? WHERE uid = ?`
	err := db.Debug().Exec(query, uaj.UserConfirmId, uaj.Status, uaj.ApplyJobId).Error
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	// Insert into history
	queryHistory := `INSERT INTO apply_job_histories 
	(uid, job_id, user_id, user_confirm_id, status, link, schedule) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`
	errHistory := db.Debug().Exec(queryHistory,
		dataQuery.Uid, dataQuery.JobId, dataQuery.UserId,
		uaj.UserConfirmId, uaj.Status, uaj.Link, uaj.Schedule,
	).Error

	if errHistory != nil {
		helper.Logger("error", "In Server: "+errHistory.Error())
		return nil, errors.New(errHistory.Error())
	}

	return map[string]any{}, nil
}

func AdminListApplyJob() (map[string]any, error) {
	var job entities.AdminListApplyJobQuery
	var candidateExercise entities.CandidateExerciseQuery
	var candidateBiodata entities.CandidateBiodataQuery
	var candidateLanguage entities.CandidateLanguageQuery
	var candidateWork entities.CandidateWorkQuery
	var candidatePlace entities.CandidatePlaceQuery
	var candidateDoc entities.CandidateDocumentQuery

	var jobFavourite []entities.JobFavourite

	var dataJob = make([]entities.AdminListApplyJob, 0)

	query := `SELECT aj.uid AS id, j.title, j.caption, j.salary,
	aj.user_id AS user_id_candidate,
	pc.fullname AS user_name_candidate,
	jc.uid as cat_id,
	jc.name AS cat_name, 
	p.id AS place_id,
	p.name AS place_name,
	p.currency AS place_currency,
	p.kurs AS place_kurs,
	p.info AS place_info,
	c.uid AS company_id,
	c.logo AS company_logo,
	c.name AS company_name,
	up.user_id,
	up.avatar AS user_avatar,
	up.fullname AS user_name,
	j.created_at,
	js.id AS job_status_id,
	js.name AS job_status_name
	FROM jobs j
	INNER JOIN job_categories jc ON jc.uid = j.cat_id
	INNER JOIN companies c ON c.uid = j.company_id 
	INNER JOIN apply_jobs aj ON aj.job_id = j.uid
	INNER JOIN job_statuses js ON js.id = aj.status
	INNER JOIN places p ON p.id = j.place_id
	INNER JOIN profiles up ON up.user_id = j.user_id
	INNER JOIN profiles pc ON pc.user_id = aj.user_id
	`
	rows, err := db.Debug().Raw(query).Rows()

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		errJobRows := db.ScanRows(rows, &job)

		if errJobRows != nil {
			helper.Logger("error", "In Server: "+errJobRows.Error())
			return nil, errors.New(errJobRows.Error())
		}

		jobFavourite = nil

		bookmarkQuery := `SELECT job_id, user_id FROM job_favourites WHERE user_id = '` + job.UserId + `' AND job_id = '` + job.Id + `'`

		errBookmark := db.Debug().Raw(bookmarkQuery).Scan(&jobFavourite).Error

		if errBookmark != nil {
			helper.Logger("error", "In Server: "+errBookmark.Error())
			return nil, errors.New(errBookmark.Error())
		}

		isJobFavouriteExist := len(jobFavourite)

		var bookmark bool

		if isJobFavouriteExist == 1 {
			bookmark = true
		} else {
			bookmark = false
		}

		salaryIdr := helper.FormatIDR(job.Salary * job.PlaceKurs)

		// Candidate Exercise

		dataCandidateExercise := make([]entities.CandidateExercise, 0)

		queryCandidateExercise := `SELECT name, institution, start_month, start_year, end_month, end_year 
		FROM form_exercises WHERE user_id = ?`

		rowsCandidateExercise, errCandidateExercise := db.Debug().Raw(queryCandidateExercise, job.UserIdCandidate).Rows()

		if errCandidateExercise != nil {
			helper.Logger("error", "In Server: "+errCandidateExercise.Error())
		}
		defer rowsCandidateExercise.Close()

		for rowsCandidateExercise.Next() {
			errCandidateExerciseRows := db.ScanRows(rowsCandidateExercise, &candidateExercise)

			if errCandidateExerciseRows != nil {
				helper.Logger("error", "In Server: "+errCandidateExerciseRows.Error())
				return nil, errors.New(errCandidateExerciseRows.Error())
			}

			dataCandidateExercise = append(dataCandidateExercise, entities.CandidateExercise{
				Name:        candidateExercise.Name,
				Institution: candidateExercise.Institution,
				StartMonth:  candidateExercise.StartMonth,
				StartYear:   candidateExercise.StartYear,
				EndMonth:    candidateExercise.EndMonth,
				EndYear:     candidateExercise.EndYear,
			})
		}

		// End Candidate Exercise

		// Candidate Biodata

		dataCandidateBiodata := make([]entities.CandidateBiodata, 0)

		queryCandidateBiodata := `SELECT birthdate, gender, weight, height, status, religion, place 
		FROM form_biodatas WHERE user_id = ?`

		rowsCandidateBiodata, errCandidateBiodata := db.Debug().Raw(queryCandidateBiodata, job.UserIdCandidate).Rows()

		if errCandidateBiodata != nil {
			helper.Logger("error", "In Server: "+errCandidateBiodata.Error())
		}
		defer rowsCandidateBiodata.Close()

		for rowsCandidateBiodata.Next() {
			errCandidateBiodataRows := db.ScanRows(rowsCandidateBiodata, &candidateBiodata)

			if errCandidateBiodataRows != nil {
				helper.Logger("error", "In Server: "+errCandidateBiodataRows.Error())
				return nil, errors.New(errCandidateBiodataRows.Error())
			}

			dataCandidateBiodata = append(dataCandidateBiodata, entities.CandidateBiodata{
				Birthdate: candidateBiodata.Birthdate,
				Gender:    candidateBiodata.Gender,
				Weight:    candidateBiodata.Weight,
				Height:    candidateBiodata.Height,
				Status:    candidateBiodata.Status,
				Religion:  candidateBiodata.Religion,
				Place:     candidateBiodata.Place,
			})
		}

		// End Candidate Biodata

		// Candidate Language

		dataCandidateLanguage := make([]entities.CandidateLanguage, 0)

		queryCandidateLanguage := `SELECT level, language 
		FROM form_languages WHERE user_id = ?`

		rowsCandidateLanguage, errCandidateLanguage := db.Debug().Raw(queryCandidateLanguage, job.UserIdCandidate).Rows()

		if errCandidateLanguage != nil {
			helper.Logger("error", "In Server: "+errCandidateLanguage.Error())
		}
		defer rowsCandidateLanguage.Close()

		for rowsCandidateLanguage.Next() {
			errCandidateLanguageRows := db.ScanRows(rowsCandidateLanguage, &candidateLanguage)

			if errCandidateLanguageRows != nil {
				helper.Logger("error", "In Server: "+errCandidateLanguageRows.Error())
				return nil, errors.New(errCandidateLanguageRows.Error())
			}

			dataCandidateLanguage = append(dataCandidateLanguage, entities.CandidateLanguage{
				Level:    candidateLanguage.Level,
				Language: candidateLanguage.Language,
			})
		}

		// End Candidate Language

		// Candidate Work

		dataCandidateWork := make([]entities.CandidateWork, 0)

		queryCandidateWork := `SELECT position, institution, work, country, city, start_month, 
		start_year, end_month, end_year, is_work 
		FROM form_works WHERE user_id = ?`

		rowsCandidateWork, errCandidateWork := db.Debug().Raw(queryCandidateWork, job.UserIdCandidate).Rows()

		if errCandidateWork != nil {
			helper.Logger("error", "In Server: "+errCandidateWork.Error())
		}
		defer rowsCandidateWork.Close()

		for rowsCandidateWork.Next() {
			errCandidateWorkRows := db.ScanRows(rowsCandidateWork, &candidateWork)

			if errCandidateWorkRows != nil {
				helper.Logger("error", "In Server: "+errCandidateWorkRows.Error())
				return nil, errors.New(errCandidateWorkRows.Error())
			}

			dataCandidateWork = append(dataCandidateWork, entities.CandidateWork{
				Position:    candidateWork.Position,
				Institution: candidateWork.Institution,
				Work:        candidateWork.Work,
				Country:     candidateWork.Country,
				City:        candidateWork.City,
				StartMonth:  candidateWork.StartMonth,
				StartYear:   candidateWork.StartYear,
				EndMonth:    candidateWork.EndMonth,
				EndYear:     candidateWork.EndYear,
				IsWork:      candidateWork.IsWork,
			})
		}

		// End Candidate Work

		// Candidate Place

		dataCandidatePlace := make([]entities.CandidatePlace, 0)

		queryCandidatePlace := `SELECT 
		p.name AS province_name, 
		r.name AS city_name, 
		d.name AS district_name, 
		s.name AS subdistrict_name
		FROM form_places fp 
		INNER JOIN provinces p ON p.id = fp.province_id
		INNER JOIN regencies r ON r.id = fp.city_id
		INNER JOIN districts d ON d.id = fp.district_id
		INNER JOIN villages s ON s.id = fp.subdistrict_id
		WHERE user_id = ?`

		rowsCandidatePlace, errCandidatePlace := db.Debug().Raw(queryCandidatePlace, job.UserIdCandidate).Rows()

		if errCandidatePlace != nil {
			helper.Logger("error", "In Server: "+errCandidatePlace.Error())
		}
		defer rowsCandidatePlace.Close()

		for rowsCandidatePlace.Next() {
			errCandidatePlaceRows := db.ScanRows(rowsCandidatePlace, &candidatePlace)

			if errCandidatePlaceRows != nil {
				helper.Logger("error", "In Server: "+errCandidatePlaceRows.Error())
				return nil, errors.New(errCandidatePlaceRows.Error())
			}

			dataCandidatePlace = append(dataCandidatePlace, entities.CandidatePlace{
				ProvinceName:    candidatePlace.ProvinceName,
				CityName:        candidatePlace.CityName,
				DistrictName:    candidatePlace.DistrictName,
				SubdistrictName: candidatePlace.SubdistrictName,
			})
		}

		// End Candidate Place

		// Candidate Document

		dataCandidateDocument := make([]entities.CandidateDocument, 0)

		queryCandidateDocument := `SELECT d.name AS document, ajd.path
		FROM documents d
		INNER JOIN apply_job_documents ajd ON ajd.doc_id = d.id
		WHERE ajd.apply_job_id = ?
		`

		rowsCandidateDocument, errCandidateDocument := db.Debug().Raw(queryCandidateDocument, job.Id).Rows()

		if errCandidateDocument != nil {
			helper.Logger("error", "In Server: "+errCandidateDocument.Error())
		}
		defer rowsCandidateDocument.Close()

		for rowsCandidateDocument.Next() {
			errCandidateDocumentRows := db.ScanRows(rowsCandidateDocument, &candidateDoc)

			if errCandidateDocumentRows != nil {
				helper.Logger("error", "In Server: "+errCandidateDocumentRows.Error())
				return nil, errors.New(errCandidateDocumentRows.Error())
			}

			dataCandidateDocument = append(dataCandidateDocument, entities.CandidateDocument{
				Document: candidateDoc.Document,
				Path:     candidateDoc.Path,
			})
		}

		// End Candidate Document

		dataJob = append(dataJob, entities.AdminListApplyJob{
			Id:        job.Id,
			Title:     job.Title,
			Caption:   job.Caption,
			Salary:    int(job.Salary),
			SalaryIDR: salaryIdr,
			Bookmark:  bookmark,
			Company: entities.JobCompany{
				Id:   job.CompanyId,
				Logo: job.CompanyLogo,
				Name: job.CompanyName,
			},
			Candidate: entities.Candidate{
				Id:                job.UserIdCandidate,
				Name:              job.UserNameCandidate,
				CandidateExercise: dataCandidateExercise,
				CandidateBiodata:  dataCandidateBiodata,
				CandidateLanguage: dataCandidateLanguage,
				CandidateWork:     dataCandidateWork,
				CandidatePlace:    dataCandidatePlace,
				CandidateDoc:      dataCandidateDocument,
			},
			Status: entities.JobStatus{
				Id:   job.JobStatusId,
				Name: job.JobStatusName,
			},
			JobCategory: entities.JobCategory{
				Id:   job.CatId,
				Name: job.CatName,
			},
			JobPlace: entities.JobPlace{
				Id:       job.PlaceId,
				Name:     job.PlaceName,
				Currency: job.PlaceCurrency,
				Kurs:     int(job.PlaceKurs),
				Info:     job.PlaceInfo,
			},
			Author: entities.AuthorJobUser{
				Id:     job.UserId,
				Avatar: job.UserAvatar,
				Name:   job.UserName,
			},
			Created: job.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return map[string]any{
		"data": dataJob,
	}, nil
}

func JobList(userId, search, salary, country, position string) (map[string]any, error) {
	var job entities.JobListQuery
	var jobFavourite []entities.JobFavourite
	var jobSkillCategory []entities.JobSkillCategory
	var dataJob = make([]entities.JobList, 0)

	query := `SELECT j.uid AS id, j.title, j.caption, j.salary, j.worker_count,
	jc.uid as cat_id,
	jc.name AS cat_name, 
	p.id AS place_id,
	p.name AS place_name,
	p.currency AS place_currency,
	p.kurs AS place_kurs,
	p.info AS place_info,
	c.uid AS company_id,
	c.logo AS company_logo,
	c.name AS company_name,
	up.user_id,
	up.avatar AS user_avatar,
	up.fullname AS user_name,
	j.created_at,
	j.salary,
	(j.salary * p.kurs) AS salary_idr
	FROM jobs j
	INNER JOIN job_categories jc ON jc.uid = j.cat_id
	INNER JOIN companies c ON c.uid = j.company_id 
	INNER JOIN places p ON p.id = j.place_id
	INNER JOIN profiles up ON up.user_id = j.user_id
	WHERE p.name LIKE '%` + country + `%'
	AND jc.name LIKE '%` + position + `%'
	AND (j.title LIKE '%` + search + `%' OR j.caption LIKE '%` + search + `%' OR jc.name LIKE '%` + search + `%')  
	`

	if salary != "" {
		query += ` AND (j.salary * p.kurs) >= '` + salary + `' `
	}

	var rows *sql.Rows
	var err error

	rows, err = db.Debug().Raw(query).Rows()

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		errJobRows := db.ScanRows(rows, &job)

		if errJobRows != nil {
			helper.Logger("error", "In Server: "+errJobRows.Error())
			return nil, errors.New(errJobRows.Error())
		}

		bookmarkQuery := `SELECT job_id, user_id FROM job_favourites WHERE user_id = ? AND job_id = ?`
		errBookmark := db.Debug().Raw(bookmarkQuery, userId, job.Id).Scan(&jobFavourite).Error

		if errBookmark != nil {
			helper.Logger("error", "In Server: "+errBookmark.Error())
			return nil, errors.New(errBookmark.Error())
		}

		isJobFavouriteExist := len(jobFavourite)
		bookmark := isJobFavouriteExist == 1

		salaryIdr := helper.FormatIDR(job.Salary * job.PlaceKurs)

		jobSkillsQuery := `SELECT jsc.uid AS id, jsc.name
		FROM job_skills js 
		INNER JOIN job_skill_categories jsc ON jsc.uid = js.cat_id
		WHERE js.job_id = ?`
		errJobSkill := db.Debug().Raw(jobSkillsQuery, job.Id).Scan(&jobSkillCategory).Error

		if errJobSkill != nil {
			helper.Logger("error", "In Server: "+errJobSkill.Error())
			return nil, errors.New(errJobSkill.Error())
		}

		dataJob = append(dataJob, entities.JobList{
			Id:      job.Id,
			Title:   job.Title,
			Caption: job.Caption,
			Skills:  jobSkillCategory,
			Company: entities.JobCompany{
				Id:   job.CompanyId,
				Logo: job.CompanyLogo,
				Name: job.CompanyName,
			},
			WorkerCount: job.WorkerCount,
			Salary:      int(job.Salary),
			SalaryIDR:   salaryIdr,
			Bookmark:    bookmark,
			JobCategory: entities.JobCategory{
				Id:   job.CatId,
				Name: job.CatName,
			},
			JobPlace: entities.JobPlace{
				Id:       job.PlaceId,
				Name:     job.PlaceName,
				Currency: job.PlaceCurrency,
				Kurs:     int(job.PlaceKurs),
				Info:     job.PlaceInfo,
			},
			Author: entities.AuthorJobUser{
				Id:     job.UserId,
				Avatar: job.UserAvatar,
				Name:   job.UserName,
			},
			Created: job.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return map[string]any{
		"data": dataJob,
	}, nil
}

func JobDetail(f *models.Job) (map[string]any, error) {
	var jobs entities.JobListQuery
	var jobFavourite []entities.JobFavourite
	var dataJob = make([]entities.JobList, 0)

	query := `SELECT j.uid AS id, j.title, j.caption, j.salary, 
		jc.uid as cat_id,
		jc.name AS cat_name, 
		p.id AS place_id,
		p.name AS place_name,
		p.currency AS place_currency,
		p.kurs AS place_kurs,
		p.info AS place_info,
		up.user_id,
		up.avatar AS user_avatar,
		up.fullname AS user_name,
		j.created_at
		FROM jobs j
		INNER JOIN job_categories jc ON jc.uid = j.cat_id
		INNER JOIN places p ON p.id = j.place_id
		INNER JOIN profiles up ON up.user_id = j.user_id
		WHERE j.uid = '` + f.Id + `'
	`

	rows, err := db.Debug().Raw(query).Rows()

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		found = true
		errJobRows := db.ScanRows(rows, &jobs)

		if errJobRows != nil {
			helper.Logger("error", "In Server: "+errJobRows.Error())
			return nil, errors.New(errJobRows.Error())
		}

		bookmarkQuery := `SELECT job_id, user_id FROM job_favourites WHERE user_id = ? AND job_id = ?`
		errBookmark := db.Debug().Raw(bookmarkQuery, f.UserId, jobs.Id).Scan(&jobFavourite).Error

		if errBookmark != nil {
			helper.Logger("error", "In Server: "+errBookmark.Error())
			return nil, errors.New(errBookmark.Error())
		}

		isJobFavouriteExist := len(jobFavourite)
		bookmark := isJobFavouriteExist == 1

		salaryIdr := helper.FormatIDR(jobs.Salary * jobs.PlaceKurs)

		dataJob = append(dataJob, entities.JobList{
			Id:        jobs.Id,
			Title:     jobs.Title,
			Caption:   jobs.Caption,
			Bookmark:  bookmark,
			Salary:    int(jobs.Salary),
			SalaryIDR: salaryIdr,
			JobCategory: entities.JobCategory{
				Id:   jobs.CatId,
				Name: jobs.CatName,
			},
			JobPlace: entities.JobPlace{
				Id:       jobs.PlaceId,
				Name:     jobs.PlaceName,
				Currency: jobs.PlaceCurrency,
				Kurs:     int(jobs.PlaceKurs),
				Info:     jobs.PlaceInfo,
			},
			Author: entities.AuthorJobUser{
				Id:     jobs.UserId,
				Avatar: jobs.UserAvatar,
				Name:   jobs.UserName,
			},
			Created: jobs.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	if !found {
		return nil, errors.New("job not found")
	}

	return map[string]any{
		"data": dataJob[0],
	}, nil
}

func JobPlace() (map[string]any, error) {
	places := []entities.JobPlace{}

	query := `SELECT id, name, currency, kurs, info FROM places`

	err := db.Debug().Raw(query).Scan(&places).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}

	return map[string]any{
		"data": places,
	}, nil
}

func JobStore(j *models.JobStore) (map[string]any, error) {

	users := []entities.User{}
	categories := []entities.JobCategory{}
	places := []entities.JobPlace{}

	checkQueryUser := `SELECT uid FROM users WHERE uid = '` + j.UserId + `'`

	errCheckUser := db.Debug().Raw(checkQueryUser).Scan(&users).Error

	if errCheckUser != nil {
		helper.Logger("error", "In Server: "+errCheckUser.Error())
		return nil, errors.New(errCheckUser.Error())
	}

	isUserExist := len(users)

	if isUserExist == 0 {
		return nil, errors.New("USER_NOT_FOUND")
	}

	checkQueryCat := `SELECT uid AS id FROM job_categories WHERE uid = '` + j.CatId + `'`

	errCheckCat := db.Debug().Raw(checkQueryCat).Scan(&categories).Error

	if errCheckCat != nil {
		helper.Logger("error", "In Server: "+errCheckCat.Error())
		return nil, errors.New(errCheckCat.Error())
	}

	isJobCategoryExist := len(categories)

	if isJobCategoryExist == 0 {
		return nil, errors.New("JOB_NOT_FOUND")
	}

	checkQueryPlace := `SELECT id, name FROM places WHERE id = ?`

	errCheckPlace := db.Debug().Raw(checkQueryPlace, j.PlaceId).Scan(&places).Error

	if errCheckPlace != nil {
		helper.Logger("error", "In Server: "+errCheckPlace.Error())
		return nil, errors.New(errCheckPlace.Error())
	}

	isPlaceExist := len(places)

	if isPlaceExist == 0 {
		return nil, errors.New("PLACE_NOT_FOUND")
	}

	j.Id = uuid.NewV4().String()

	query := `INSERT INTO jobs (uid, title, caption, salary, worker_count, cat_id, company_id, place_id, user_id, is_draft) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	err := db.Debug().Exec(query, j.Id, j.Title, j.Caption, j.Salary, j.WorkerCount, j.CatId, j.CompanyId, j.PlaceId, j.UserId, j.IsDraft).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func JobUpdate(j *models.JobUpdate) (map[string]any, error) {

	categories := []entities.JobCategory{}
	places := []entities.JobPlace{}

	checkQueryCat := `SELECT uid AS id FROM job_categories WHERE uid = '` + j.CatId + `'`

	errCheckCat := db.Debug().Raw(checkQueryCat).Scan(&categories).Error

	if errCheckCat != nil {
		helper.Logger("error", "In Server: "+errCheckCat.Error())
		return nil, errors.New(errCheckCat.Error())
	}

	isJobCategoryExist := len(categories)

	if isJobCategoryExist == 0 {
		return nil, errors.New("JOB_NOT_FOUND")
	}

	checkQueryPlace := `SELECT id, name FROM places WHERE id = ?`

	errCheckPlace := db.Debug().Raw(checkQueryPlace, j.PlaceId).Scan(&places).Error

	if errCheckPlace != nil {
		helper.Logger("error", "In Server: "+errCheckPlace.Error())
		return nil, errors.New(errCheckPlace.Error())
	}

	isPlaceExist := len(places)

	if isPlaceExist == 0 {
		return nil, errors.New("PLACE_NOT_FOUND")
	}

	query := `UPDATE jobs SET title = ?, caption = ?, salary = ?, worker_count = ?, company_id = ?, place_id = ?, cat_id = ?, place_id = ?, is_draft = ?
	WHERE uid = ?`

	err := db.Debug().Exec(query, j.Title, j.Caption, j.Salary, j.WorkerCount, j.CompanyId, j.PlaceId, j.CatId, j.PlaceId, j.IsDraft, j.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func JobFavourite(j *models.JobFavourite) (map[string]any, error) {
	jobFavourite := []entities.JobFavourite{}
	jobs := []entities.Job{}
	users := []entities.User{}

	checkQueryJob := `SELECT uid AS id FROM jobs WHERE uid = '` + j.JobId + `'`

	errCheckJob := db.Debug().Raw(checkQueryJob).Scan(&jobs).Error

	if errCheckJob != nil {
		helper.Logger("error", "In Server: "+errCheckJob.Error())
		return nil, errors.New(errCheckJob.Error())
	}

	isJobExist := len(jobs)

	if isJobExist == 0 {
		return nil, errors.New("job not found")
	}

	checkQueryUser := `SELECT uid AS id FROM users WHERE uid = '` + j.UserId + `'`

	errCheckUser := db.Debug().Raw(checkQueryUser).Scan(&users).Error

	if errCheckUser != nil {
		helper.Logger("error", "In Server: "+errCheckUser.Error())
		return nil, errors.New(errCheckUser.Error())
	}

	isUserExist := len(users)

	if isUserExist == 0 {
		return nil, errors.New("user not found")
	}

	checkQueryJobFavourite := `SELECT user_id, job_id FROM job_favourites WHERE user_id = '` + j.UserId + `' AND job_id = '` + j.JobId + `'`

	errCheckJobFavourite := db.Debug().Raw(checkQueryJobFavourite).Scan(&jobFavourite).Error

	if errCheckJobFavourite != nil {
		helper.Logger("error", "In Server: "+errCheckJobFavourite.Error())
		return nil, errors.New(errCheckJobFavourite.Error())
	}

	isJobFavouriteExist := len(jobFavourite)

	if isJobFavouriteExist == 0 {
		query := `INSERT INTO job_favourites (user_id, job_id) VALUES (?, ?)`

		err := db.Debug().Exec(query, j.UserId, j.JobId).Error

		if err != nil {
			helper.Logger("error", "In Server: "+err.Error())
			return nil, errors.New(err.Error())
		}
	} else {
		query := `DELETE FROM job_favourites WHERE user_id = ? AND job_id = ?`

		err := db.Debug().Exec(query, j.UserId, j.JobId).Error

		if err != nil {
			helper.Logger("error", "In Server: "+err.Error())
			return nil, errors.New(err.Error())
		}
	}

	return map[string]any{}, nil
}

func JobCategoryCount() (map[string]any, error) {
	categoryCounts := []entities.JobCategoryCount{}

	query := `SELECT jc.name, COUNT(j.cat_id) AS total
	FROM job_categories jc
	LEFT JOIN jobs j ON jc.uid = j.cat_id
	GROUP BY jc.name
	ORDER BY total DESC`

	err := db.Debug().Raw(query).Scan(&categoryCounts).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": categoryCounts,
	}, nil
}

func JobCategoryStore(j *models.JobCategoryStore) (map[string]any, error) {

	j.Id = uuid.NewV4().String()

	query := `INSERT INTO job_categories (uid, logo, name) VALUES (?, ?, ?)`

	err := db.Debug().Exec(query, j.Id, j.Icon, j.Name).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func JobCategoryUpdate(j *models.JobCategoryUpdate) (map[string]any, error) {

	query := `UPDATE job_categories SET name = ?, logo = ? WHERE uid = ?`

	err := db.Debug().Exec(query, j.Name, j.Icon, j.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func JobCategoryDelete(j *models.JobCategoryDelete) (map[string]any, error) {

	query := `DELETE FROM job_categories WHERE uid = ?`

	err := db.Debug().Exec(query, j.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func JobCategory() (map[string]any, error) {
	categories := []entities.JobCategory{}

	query := `SELECT uid AS id, logo AS icon, name FROM job_categories`

	err := db.Debug().Raw(query).Scan(&categories).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	isCategoryExist := len(categories)

	if isCategoryExist == 0 {
		return nil, errors.New("job not found")
	}

	return map[string]any{
		"data": categories,
	}, nil
}

func JobDelete(j *models.Job) (map[string]any, error) {
	errDeleteJob := db.Debug().Exec(`
	DELETE FROM jobs WHERE uid = ?`, j.Id).Error

	if errDeleteJob != nil {
		helper.Logger("error", "In Server: "+errDeleteJob.Error())
		return nil, errors.New(errDeleteJob.Error())
	}

	return map[string]any{}, nil
}
