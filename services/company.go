package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"

	uuid "github.com/satori/go.uuid"
)

func CompanyList() (map[string]any, error) {
	companies := []entities.CompanyListQuery{}

	query := `SELECT uid AS id, logo, name FROM companies`

	err := db.Debug().Raw(query).Scan(&companies).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": companies,
	}, nil
}

func CompanyStore(c *entities.CompanyStore) (map[string]any, error) {
	Id := uuid.NewV4().String()

	query := `INSERT INTO companies (uid, logo, name) VALUES (?, ?, ?)`

	err := db.Debug().Exec(query, Id, c.Logo, c.Name).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}
