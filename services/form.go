package services

import (
	"errors"

	helper "superapps/helpers"
	"superapps/models"
)

func FormBiodata(f *models.FormBiodata) (map[string]any, error) {

	query := `INSERT INTO form_biodata (place, birthdate, gender, height, weight, religion, status) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	err := db.Debug().Exec(query, f.Place, f.Birthdate, f.Gender, f.Height, f.Weight, f.Religion, f.Status).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}
