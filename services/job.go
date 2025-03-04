package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
)

func JobList() (map[string]any, error) {
	jobs := []entities.JobList{}

	query := `SELECT j.uid AS id, j.title, j.description, j.salary, 
	jc.name AS job_name, 
	p.name AS place_name
	FROM jobs j
	INNER JOIN job_categories jc ON jc.uid = j.cat_id
	INNER JOIN places p ON p.id = j.place_id`

	err := db.Debug().Raw(query).Scan(&jobs).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}

	return map[string]any{
		"data": jobs,
	}, nil
}

func JobStore() (map[string]any, error) {
	jobs := entities.JobStore{}

	query := `INSERT INTO `

	err := db.Debug().Exec(query).Scan(&jobs).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

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
