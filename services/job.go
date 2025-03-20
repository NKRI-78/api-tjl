package services

import (
	"database/sql"
	"errors"
	"strconv"
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
		js.name AS status, aj.created_at, aj.uid AS apply_job_id, aj.link, aj.schedule
		FROM apply_job_histories aj 
		INNER JOIN job_statuses js ON js.id = aj.status
		INNER JOIN profiles paa ON paa.user_id = aj.user_id
		LEFT JOIN profiles pac ON pac.user_id = aj.user_confirm_id 
		WHERE aj.user_id = ?
		ORDER BY aj.created_at DESC
		LIMIT 1
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

		defaultIfEmpty := func(value, defaultValue string) string {
			if value == "" {
				return defaultValue
			}
			return value
		}

		data = append(data, entities.ResultInfoJob{
			Id:        dataQuery.ApplyJobId,
			Status:    dataQuery.Status,
			CreatedAt: helper.FormatDate(dataQuery.CreatedAt),
			Link:      defaultIfEmpty(dataQuery.Link, "-"),
			Schedule:  defaultIfEmpty(dataQuery.Schedule, "-"),
			UserApply: entities.UserApply{
				Id:   dataQuery.ApplyUserId,
				Name: dataQuery.ApplyUserName,
			},
			UserConfirm: entities.UserConfirm{
				Id:   defaultIfEmpty(dataQuery.ConfirmUserId, "-"),
				Name: defaultIfEmpty(dataQuery.ConfirmUserName, "-"),
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
	var data []entities.ResultInfoJob

	query := `SELECT paa.user_id AS apply_user_id, paa.fullname AS apply_user_name, 
		pac.user_id AS confirm_user_id, pac.fullname AS confirm_user_name,
		js.name AS status, aj.created_at, aj.uid AS apply_job_id, aj.link, aj.schedule
		FROM apply_job_histories aj 
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

		defaultIfEmpty := func(value, defaultValue string) string {
			if value == "" {
				return defaultValue
			}
			return value
		}

		data = append(data, entities.ResultInfoJob{
			Id:        dataQuery.ApplyJobId,
			Status:    dataQuery.Status,
			CreatedAt: helper.FormatDate(dataQuery.CreatedAt),
			Link:      defaultIfEmpty(dataQuery.Link, "-"),
			Schedule:  defaultIfEmpty(dataQuery.Schedule, "-"),
			UserApply: entities.UserApply{
				Id:   dataQuery.ApplyUserId,
				Name: dataQuery.ApplyUserName,
			},
			UserConfirm: entities.UserConfirm{
				Id:   defaultIfEmpty(dataQuery.ConfirmUserId, "-"),
				Name: defaultIfEmpty(dataQuery.ConfirmUserName, "-"),
			},
		})
	}

	return map[string]any{
		"data": data,
	}, nil
}

func ApplyJob(aj *models.ApplyJob) (map[string]any, error) {

	query := `INSERT INTO apply_jobs (uid, job_id, user_id) VALUES (?, ?, ?)`

	err := db.Debug().Exec(query, aj.Id, aj.JobId, aj.UserId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	queryHistory := `INSERT INTO apply_job_histories (uid, job_id, user_id) VALUES (?, ?, ?)`

	errHistory := db.Debug().Exec(queryHistory, aj.Id, aj.JobId, aj.UserId).Error

	if errHistory != nil {
		helper.Logger("error", "In Server: "+errHistory.Error())
		return nil, errors.New(errHistory.Error())
	}

	return map[string]any{}, nil
}

func UpdateApplyJob(uaj *models.ApplyJob) (map[string]any, error) {

	var dataQuery entities.ApplyJobQuery

	query := `UPDATE apply_jobs SET user_confirm_id = ?, status = ? WHERE uid = ?`

	err := db.Debug().Exec(query, uaj.UserConfirmId, uaj.Status, uaj.ApplyJobId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	queryInfo := `SELECT * FROM apply_jobs WHERE uid = ?`

	rows, errInfo := db.Debug().Raw(queryInfo, uaj.ApplyJobId).Rows()

	if errInfo != nil {
		helper.Logger("error", "In Server: "+errInfo.Error())
		return nil, errors.New(errInfo.Error())
	}

	for rows.Next() {
		errJobRows := db.ScanRows(rows, &dataQuery)

		if errJobRows != nil {
			helper.Logger("error", "In Server: "+errJobRows.Error())
			return nil, errors.New(errJobRows.Error())
		}

		queryHistory := `INSERT INTO apply_job_histories 
		(uid, job_id, user_id, user_confirm_id, status, link, schedule) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`

		errHistory := db.Debug().Exec(queryHistory,
			dataQuery.Uid, dataQuery.JobId,
			dataQuery.UserId, dataQuery.UserConfirmId,
			uaj.Status, uaj.Link, uaj.Schedule,
		).Error

		if errHistory != nil {
			helper.Logger("error", "In Server: "+errHistory.Error())
			return nil, errors.New(errHistory.Error())
		}
	}

	return map[string]any{}, nil
}

func AdminJobList() (map[string]any, error) {
	var jobs entities.JobListAdminQuery
	var jobFavourite []entities.JobFavourite

	var dataJob = make([]entities.JobListAdmin, 0)

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
	j.created_at,
	js.id AS job_status_id,
	js.name AS job_status_name
	FROM jobs j
	INNER JOIN job_categories jc ON jc.uid = j.cat_id
	INNER JOIN apply_jobs aj ON aj.job_id = j.uid
	INNER JOIN job_statuses js ON js.id = aj.status
	INNER JOIN places p ON p.id = j.place_id
	INNER JOIN profiles up ON up.user_id = j.user_id
	`
	rows, err := db.Debug().Raw(query).Rows()

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		errJobRows := db.ScanRows(rows, &jobs)

		if errJobRows != nil {
			helper.Logger("error", "In Server: "+errJobRows.Error())
			return nil, errors.New(errJobRows.Error())
		}

		bookmarkQuery := `SELECT job_id, user_id FROM job_favourites WHERE user_id = '` + jobs.UserId + `' AND job_id = '` + jobs.Id + `'`

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

		salaryIdr := helper.FormatIDR(jobs.Salary * jobs.PlaceKurs)

		dataJob = append(dataJob, entities.JobListAdmin{
			Id:        jobs.Id,
			Title:     jobs.Title,
			Caption:   jobs.Caption,
			Salary:    int(jobs.Salary),
			SalaryIDR: salaryIdr,
			Bookmark:  bookmark,
			Status: entities.JobStatus{
				Id:   jobs.JobStatusId,
				Name: jobs.JobStatusName,
			},
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
			JobUser: entities.JobUser{
				Id:     jobs.UserId,
				Avatar: jobs.UserAvatar,
				Name:   jobs.UserName,
			},
			Created: jobs.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return map[string]any{
		"data": dataJob,
	}, nil
}

func JobList(salary, country, position string) (map[string]any, error) {
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
	j.created_at,
	j.salary,
	(j.salary * p.kurs) AS salary_idr
	FROM jobs j
	INNER JOIN job_categories jc ON jc.uid = j.cat_id
	INNER JOIN places p ON p.id = j.place_id
	INNER JOIN profiles up ON up.user_id = j.user_id
	WHERE p.name LIKE '%` + country + `%'
	AND jc.name LIKE '%` + position + `%'
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
		errJobRows := db.ScanRows(rows, &jobs)

		if errJobRows != nil {
			helper.Logger("error", "In Server: "+errJobRows.Error())
			return nil, errors.New(errJobRows.Error())
		}

		bookmarkQuery := `SELECT job_id, user_id FROM job_favourites WHERE user_id = ? AND job_id = ?`
		errBookmark := db.Debug().Raw(bookmarkQuery, jobs.UserId, jobs.Id).Scan(&jobFavourite).Error

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
			Salary:    int(jobs.Salary),
			SalaryIDR: salaryIdr,
			Bookmark:  bookmark,
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
			JobUser: entities.JobUser{
				Id:     jobs.UserId,
				Avatar: jobs.UserAvatar,
				Name:   jobs.UserName,
			},
			Created: jobs.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return map[string]any{
		"data": dataJob,
	}, nil
}

func JobDetail(f *models.Job) (map[string]any, error) {
	var jobs entities.JobListQuery
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

		salaryIdr := helper.FormatIDR(jobs.Salary * jobs.PlaceKurs)

		dataJob = append(dataJob, entities.JobList{
			Id:        jobs.Id,
			Title:     jobs.Title,
			Caption:   jobs.Caption,
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
			JobUser: entities.JobUser{
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

	query := `SELECT id, name, currency FROM places`

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
		return nil, errors.New("user not found")
	}

	checkQueryCat := `SELECT uid AS id FROM job_categories WHERE uid = '` + j.CatId + `'`

	errCheckCat := db.Debug().Raw(checkQueryCat).Scan(&categories).Error

	if errCheckCat != nil {
		helper.Logger("error", "In Server: "+errCheckCat.Error())
		return nil, errors.New(errCheckCat.Error())
	}

	isJobCategoryExist := len(categories)

	if isJobCategoryExist == 0 {
		return nil, errors.New("job not found")
	}

	checkQueryPlace := `SELECT id, name FROM places WHERE id = '` + strconv.Itoa(j.PlaceId) + `'`

	errCheckPlace := db.Debug().Raw(checkQueryPlace).Scan(&places).Error

	if errCheckPlace != nil {
		helper.Logger("error", "In Server: "+errCheckPlace.Error())
		return nil, errors.New(errCheckPlace.Error())
	}

	isPlaceExist := len(places)

	if isPlaceExist == 0 {
		return nil, errors.New("place not found")
	}

	j.Id = uuid.NewV4().String()

	query := `INSERT INTO jobs (uid, title, caption, salary, cat_id, place_id, user_id, is_draft) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	err := db.Debug().Exec(query, j.Id, j.Title, j.Caption, j.Salary, j.CatId, j.PlaceId, j.UserId, j.IsDraft).Error

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

	query := `SELECT jc.name, COUNT(*) AS total
	FROM jobs j
	INNER JOIN job_categories jc ON jc.uid = j.cat_id
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

func JobCategory() (map[string]any, error) {
	categories := []entities.JobCategory{}

	query := `SELECT uid AS id, name FROM job_categories`

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
