package services

import (
	"errors"

	helper "superapps/helpers"
	"superapps/models"
)

func FormBiodata(f *models.FormBiodata) (map[string]any, error) {

	query := `INSERT INTO form_biodata (place, birthdate, gender, height, weight, religion, status, user_id) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	err := db.Debug().Exec(query, f.Place, f.Birthdate, f.Gender, f.Height, f.Weight, f.Religion, f.Status, f.UserId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func FormRegion(f *models.FormRegion) (map[string]any, error) {

	query := `INSERT INTO form_regions (province_id, city_id, district_id, subdistrict_id, user_id, detail_address) 
	VALUES (?, ?, ?, ?, ?, ?)`

	err := db.Debug().Exec(query, f.ProvinceId, f.CityId, f.DistrictId, f.SubdistrictId, f.UserId, f.DetailAddress).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func FormEducation(f *models.FormEducation) (map[string]any, error) {
	query := `INSERT INTO form_educations (education_level, major, school_or_college, start_year, start_month, end_year, end_month) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	err := db.Debug().Exec(query, f.EducationLevel, f.Major, f.SchoolOrCollege, f.StartYear, f.StartMonth, f.EndYear, f.EndMonth).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func FormExercise(f *models.FormExercise) (map[string]any, error) {
	query := `INSERT INTO form_exercises (name, institution, start_year, start_month, start_year, end_year, end_month) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	err := db.Debug().Exec(query, f.Name, f.Institution, f.StartYear, f.StartMonth, f.EndYear, f.EndMonth).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func FormWork(f *models.FormWork) (map[string]any, error) {
	query := `INSERT INTO form_works (position, work, country, city, start_year, start_month, end_year, end_month) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	err := db.Debug().Exec(query, f.Position, f.Work, f.Country, f.City, f.StartYear, f.StartMonth, f.EndYear, f.EndMonth).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}
