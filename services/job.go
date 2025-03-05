package services

import (
	"errors"
	"strconv"
	entities "superapps/entities"
	helper "superapps/helpers"
	models "superapps/models"

	uuid "github.com/satori/go.uuid"
)

func JobList() (map[string]any, error) {
	var jobs entities.JobListQuery
	var dataJob = make([]entities.JobList, 0)

	query := `SELECT j.uid AS id, j.title, j.caption, j.salary, 
	jc.uid as cat_id,
	jc.name AS cat_name, 
	p.id AS place_id,
	p.name AS place_name,
	p.currency AS place_currency,
	up.user_id,
	up.avatar AS user_avatar,
	up.fullname AS user_name,
	j.created_at
	FROM jobs j
	INNER JOIN job_categories jc ON jc.uid = j.cat_id
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

		dataJob = append(dataJob, entities.JobList{
			Id:      jobs.Id,
			Title:   jobs.Title,
			Caption: jobs.Caption,
			Salary:  jobs.Salary,
			JobCategory: entities.JobCategory{
				Id:   jobs.CatId,
				Name: jobs.CatName,
			},
			JobPlace: entities.JobPlace{
				Id:       jobs.PlaceId,
				Name:     jobs.PlaceName,
				Currency: jobs.PlaceCurrency,
			},
			JobUser: entities.JobUser{
				Id:     jobs.UserId,
				Avatar: jobs.UserAvatar,
				Name:   jobs.UserName,
			},
			Created: jobs.CreatedAt,
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

		dataJob = append(dataJob, entities.JobList{
			Id:      jobs.Id,
			Title:   jobs.Title,
			Caption: jobs.Caption,
			Salary:  jobs.Salary,
			JobCategory: entities.JobCategory{
				Id:   jobs.CatId,
				Name: jobs.CatName,
			},
			JobPlace: entities.JobPlace{
				Id:       jobs.PlaceId,
				Name:     jobs.PlaceName,
				Currency: jobs.PlaceCurrency,
			},
			JobUser: entities.JobUser{
				Id:     jobs.UserId,
				Avatar: jobs.UserAvatar,
				Name:   jobs.UserName,
			},
			Created: jobs.CreatedAt,
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
