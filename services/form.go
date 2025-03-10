package services

import (
	"errors"

	helper "superapps/helpers"
	"superapps/models"
)

func FormBiodata(f *models.FormBiodata) (map[string]any, error) {
	query := `INSERT INTO form_biodatas (place, birthdate, gender, height, weight, religion, status, user_id) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	err := db.Debug().Exec(query, f.Place, f.Birthdate, f.Gender, f.Height, f.Weight, f.Religion, f.Status, f.UserId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func FormRegion(f *models.FormRegion) (map[string]any, error) {
	query := `INSERT INTO form_regions (province_id, city_id, district_id, subdistrict_id, detail_address, user_id) 
	VALUES (?, ?, ?, ?, ?, ?)`

	err := db.Debug().Exec(query, f.ProvinceId, f.CityId, f.DistrictId, f.SubdistrictId, f.DetailAddress, f.UserId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func FormPlace(f *models.FormPlace) (map[string]any, error) {
	query := `INSERT INTO form_places (province, city, district, subdistrict, detail_address, user_id) 
	VALUES (?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE 
	province = VALUES(province), 
	city = VALUES(city), 
	district = VALUES(district), 
	subdistrict = VALUES(subdistrict), 
	detail_address = VALUES(detail_address)`

	err := db.Debug().Exec(query, f.Province, f.City, f.District, f.Subdistrict, f.DetailAddress, f.UserId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func FormEducation(f *models.FormEducation) (map[string]any, error) {
	query := `INSERT INTO form_educations (education_level, major, school_or_college, start_year, start_month, end_year, end_month, user_id) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	err := db.Debug().Exec(query, f.EducationLevel, f.Major, f.SchoolOrCollege, f.StartYear, f.StartMonth, f.EndYear, f.EndMonth, f.UserId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func FormExercise(f *models.FormExercise) (map[string]any, error) {
	query := `INSERT INTO form_exercises (name, institution, start_year, start_month, end_year, end_month, user_id) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	err := db.Debug().Exec(query, f.Name, f.Institution, f.StartYear, f.StartMonth, f.EndYear, f.EndMonth, f.UserId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func FormWork(f *models.FormWork) (map[string]any, error) {
	query := `INSERT INTO form_works (position, work, country, city, start_year, start_month, end_year, end_month, user_id) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	err := db.Debug().Exec(query, f.Position, f.Work, f.Country, f.City, f.StartYear, f.StartMonth, f.EndYear, f.EndMonth, f.UserId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func FormLanguage(f *models.FormLanguage) (map[string]any, error) {
	query := `INSERT INTO form_languages (level, language, user_id) 
	VALUES (?, ?, ?)`

	err := db.Debug().Exec(query, f.Level, f.Language, f.UserId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}
