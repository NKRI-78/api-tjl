package services

import (
	"errors"

	helper "superapps/helpers"
	"superapps/models"
)

func FormBiodata(f *models.FormBiodata) (map[string]any, error) {
	query := `INSERT INTO form_biodatas (place, birthdate, gender, height, weight, religion, status, user_id) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE 
		place = VALUES(place), 
		birthdate = VALUES(birthdate), 
		gender = VALUES(gender), 
		height = VALUES(height), 
		weight = VALUES(weight), 
		religion = VALUES(religion), 
		status = VALUES(status)`

	err := db.Debug().Exec(query, f.Place, f.Birthdate, f.Gender, f.Height, f.Weight, f.Religion, f.Status, f.UserId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func FormPlace(f *models.FormPlace) (map[string]any, error) {
	query := `INSERT INTO form_places (province_id, city_id, district_id, subdistrict_id, detail_address, user_id) 
	VALUES (?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE 
	province_id = VALUES(province_id), 
	city_id = VALUES(city_id), 
	district_id = VALUES(district_id), 
	subdistrict_id = VALUES(subdistrict_id), 
	detail_address = VALUES(detail_address)`

	err := db.Debug().Exec(query, f.ProvinceId, f.CityId, f.DistrictId, f.SubdistrictId, f.DetailAddress, f.UserId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func FormEducation(f *models.FormEducation) (map[string]any, error) {
	query := `INSERT INTO form_educations 
	(education_level, major, school_or_college, start_year, start_month, end_year, end_month, user_id) 
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

	result, err := db.DB().Exec(query, f.Name, f.Institution, f.StartYear, f.StartMonth, f.EndYear, f.EndMonth, f.UserId)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	queryExerciseMedia := `INSERT INTO form_exercise_medias (exercise_id, path) 
	VALUES (?, ?)`

	for _, v := range f.Certificates {
		errExerciseMedia := db.Debug().Exec(queryExerciseMedia, lastID, v).Error

		if errExerciseMedia != nil {
			helper.Logger("error", "In Server: "+errExerciseMedia.Error())
			return nil, errors.New(errExerciseMedia.Error())
		}
	}

	return map[string]any{}, nil
}

func FormWork(f *models.FormWork) (map[string]any, error) {
	query := `INSERT INTO form_works (position, institution, work, is_work, country, city, start_year, start_month, end_year, end_month, user_id) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	err := db.Debug().Exec(query, f.Position, f.Institution, f.Work, f.IsWork, f.Country, f.City, f.StartYear, f.StartMonth, f.EndYear, f.EndMonth, f.UserId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func DeleteFormLanguage(f *models.FormLanguage) (map[string]any, error) {
	query := `DELETE FROM form_languages WHERE id = ?`

	err := db.Debug().Exec(query, f.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func DeleteFormBiodata(f *models.FormBiodata) (map[string]any, error) {
	query := `DELETE FROM form_biodatas WHERE id = ?`

	err := db.Debug().Exec(query, f.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func DeleteFormEducation(f *models.FormEducation) (map[string]any, error) {
	query := `DELETE FROM form_educations WHERE id = ?`

	err := db.Debug().Exec(query, f.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func DeleteFormAddress(f *models.FormPlace) (map[string]any, error) {
	query := `DELETE FROM form_places WHERE id = ?`

	err := db.Debug().Exec(query, f.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func DeleteFormWork(f *models.FormWork) (map[string]any, error) {
	query := `DELETE FROM form_works WHERE id = ?`

	err := db.Debug().Exec(query, f.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func DeleteFormExercise(f *models.FormExercise) (map[string]any, error) {
	query := `DELETE FROM form_exercises WHERE id = ?`

	err := db.Debug().Exec(query, f.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func UpdateFormLanguage(f *models.FormLanguage) (map[string]any, error) {
	query := `UPDATE form_languages SET level = ?, 
	language = ? WHERE id = ?`

	err := db.Debug().Exec(query, f.Level, f.Language, f.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func UpdateFormExercise(f *models.FormExercise) (map[string]any, error) {
	query := `UPDATE form_exercises SET institution = ?, 
	name = ?, start_year = ?, start_month = ?, end_year = ?, 
	end_month = ? WHERE id = ?`

	err := db.Debug().Exec(query, f.Institution, f.Name, f.StartYear, f.StartMonth, f.EndYear, f.EndMonth, f.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	queryInsertExerciseMedia := `INSERT INTO form_exercise_medias (exercise_id, path) 
	VALUES (?, ?)`

	for _, v := range f.Certificates {
		errExerciseMedia := db.Debug().Exec(queryInsertExerciseMedia, f.Id, v).Error

		if errExerciseMedia != nil {
			helper.Logger("error", "In Server: "+errExerciseMedia.Error())
			return nil, errors.New(errExerciseMedia.Error())
		}
	}

	queryDeleteExerciseMedia := `DELETE FROM form_exercise_medias WHERE id = ?`

	for _, v := range f.CertificateIDeletes {
		errExerciseMedia := db.Debug().Exec(queryDeleteExerciseMedia, v).Error

		if errExerciseMedia != nil {
			helper.Logger("error", "In Server: "+errExerciseMedia.Error())
			return nil, errors.New(errExerciseMedia.Error())
		}
	}

	return map[string]any{}, nil
}

func UpdateFormWork(f *models.FormWork) (map[string]any, error) {
	query := `UPDATE form_works SET position = ?, institution = ?,
	city = ?, work = ?, is_work = ?, start_year = ?, start_month = ?, end_year = ?, 
	end_month = ? WHERE id = ?`

	var isWork int

	if f.StillWork {
		isWork = 1
	} else {
		isWork = 0
	}

	err := db.Debug().Exec(query, f.Position, f.Institution, f.City, f.Work, isWork, f.StartYear, f.StartMonth, f.EndYear, f.EndMonth, f.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func UpdateFormEducation(f *models.FormEducation) (map[string]any, error) {
	query := `UPDATE form_educations SET education_level = ?, 
	major = ?, school_or_college = ?, start_year = ?, start_month = ?, end_year = ?, 
	end_month = ? WHERE id = ?`

	err := db.Debug().Exec(query, f.EducationLevel, f.Major, f.SchoolOrCollege, f.StartYear, f.StartMonth, f.EndYear, f.EndMonth, f.Id).Error

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
