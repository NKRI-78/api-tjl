package services

import (
	"errors"

	entities "superapps/entities"
	helper "superapps/helpers"
	models "superapps/models"
)

func Country() (map[string]any, error) {
	countries := []entities.Country{}

	query := `SELECT id, name FROM places`

	err := db.Debug().Raw(query).Scan(&countries).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": countries,
	}, nil
}

func CountryDelete(c *models.Country) (map[string]any, error) {

	query := `DELETE FROM places WHERE id = ?`

	err := db.Debug().Exec(query, c.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func CountryStore(c *models.CountryStore) (map[string]any, error) {
	query := `INSERT INTO places SET name = ?, currency = ?, kurs = ?, info = ?`

	err := db.Debug().Exec(query, c.Name, c.Currency, c.Kurs, c.Info).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func CountryUpdate(c *models.CountryUpdate) (map[string]any, error) {

	query := `UPDATE places SET name = ?, currency = ?, kurs = ?, info = ? WHERE id = ?`

	err := db.Debug().Exec(query, c.Name, c.Currency, c.Kurs, c.Info, c.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func Province() (map[string]any, error) {
	provinces := []entities.Province{}

	query := `SELECT id, name
	FROM provinces`

	err := db.Debug().Raw(query).Scan(&provinces).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": provinces,
	}, nil
}

func City(c *models.City) (map[string]any, error) {
	cities := []entities.City{}

	query := `SELECT id, province_id, name
	FROM regencies WHERE province_id = '` + c.ProvinceId + `'`

	err := db.Debug().Raw(query).Scan(&cities).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": cities,
	}, nil
}

func District(d *models.District) (map[string]any, error) {
	districts := []entities.District{}

	query := `SELECT id, regency_id, name
	FROM districts WHERE regency_id = '` + d.RegencyId + `'`

	err := db.Debug().Raw(query).Scan(&districts).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": districts,
	}, nil
}

func Subdistrict(d *models.Subdistrict) (map[string]any, error) {
	subdistricts := []entities.Subdistrict{}

	query := `SELECT id, district_id, name
	FROM villages WHERE district_id = '` + d.DistrictId + `'`

	err := db.Debug().Raw(query).Scan(&subdistricts).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": subdistricts,
	}, nil
}
