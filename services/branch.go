package services

import (
	"errors"

	entities "superapps/entities"
	helper "superapps/helpers"
)

func Branch() (map[string]any, error) {
	branches := []entities.Branch{}

	query := `SELECT id, name
	FROM branchs`

	err := db.Debug().Raw(query).Scan(&branches).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": branches,
	}, nil
}
