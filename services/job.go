package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
)

func GetJobCategory() (map[string]interface{}, error) {
	jobs := []entities.JobCategory{}

	query := `SELECT uid AS id, name FROM job_categories`

	err := db.Debug().Raw(query).Scan(&jobs).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	isJobExist := len(jobs)

	if isJobExist == 0 {
		return nil, errors.New("job not found")
	}

	return map[string]any{
		"data": jobs,
	}, nil
}
