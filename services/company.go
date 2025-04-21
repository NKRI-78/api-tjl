package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"

	uuid "github.com/satori/go.uuid"
)

func CompanyList() (map[string]any, error) {
	var company entities.CompanyListQuery

	var dataCompany = make([]entities.CompanyList, 0)

	query := `SELECT 
	c.uid AS id, 
	c.logo, c.name,
	p.name AS company_name
	FROM companies c
	INNER JOIN places p ON p.id = c.place_id
	`

	rows, err := db.Debug().Raw(query).Rows()

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		errJobRows := db.ScanRows(rows, &company)

		if errJobRows != nil {
			helper.Logger("error", "In Server: "+errJobRows.Error())
			return nil, errors.New(errJobRows.Error())
		}

		dataCompany = append(dataCompany, entities.CompanyList{
			Id:     company.Id,
			Logo:   company.Logo,
			Name:   company.Name,
			Origin: company.CompanyName,
		})
	}

	return map[string]any{
		"data": dataCompany,
	}, nil
}

func CompanyStore(c *entities.CompanyStore) (map[string]any, error) {
	Id := uuid.NewV4().String()

	query := `INSERT INTO companies (uid, logo, name, place_id) VALUES (?, ?, ?, ?)`

	err := db.Debug().Exec(query, Id, c.Logo, c.Name, c.PlaceId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func CompanyUpdate(c *entities.CompanyUpdate) (map[string]any, error) {

	query := `UPDATE companies SET logo = ?, name = ?, place_id = ? WHERE uid = ?`

	err := db.Debug().Exec(query, c.Logo, c.Name, c.PlaceId, c.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func CompanyDelete(c *entities.CompanyDelete) (map[string]any, error) {

	err := db.Debug().Exec(`
		DELETE FROM companies WHERE uid = ?
	`, c.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}
