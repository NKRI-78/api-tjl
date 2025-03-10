package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
)

func Language() (map[string]any, error) {
	languages := []entities.Language{}

	query := `SELECT id, name FROM languages`

	err := db.Debug().Raw(query).Scan(&languages).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": languages,
	}, nil
}
