package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
	models "superapps/models"
)

func GetProfile(p *models.Profile) (map[string]interface{}, error) {
	profiles := []entities.Profile{}
	education := entities.ProfileFormEducation{}
	exercise := entities.ProfileFormExercise{}
	educationMedia := entities.FormExerciseCertificate{}
	exerciseMedia := entities.FormExerciseCertificate{}
	work := entities.ProfileFormWorkQuery{}
	language := entities.ProfileFormLanguage{}

	query := `SELECT p.user_id AS id, p.fullname, p.avatar, u.phone, u.email, u.enabled, 
	jc.uid AS job_id,
	jc.name AS job_name,
	ur.name AS role,
	fb.id AS bio_id,
	fb.birthdate AS bio_birthdate,
	fb.gender AS bio_gender,
	fb.weight AS bio_weight,
	fb.height AS bio_height,
	fb.religion AS bio_religion,
	fb.status AS bio_status,
	fb.place AS bio_place,
	fp.detail_address AS bio_detail_address,
	fp.id AS bio_address_id,
	pro.id AS bio_province_id,
	pro.name AS bio_province,
	reg.id AS bio_city_id,
	reg.name AS bio_city,
	dis.id AS bio_district_id,
	dis.name AS bio_district, 
	vil.id AS bio_subdistrict_id,
	vil.name AS bio_subdistrict 
	FROM profiles p 
	INNER JOIN users u ON u.uid = p.user_id
	INNER JOIN user_roles ur ON ur.id = u.role
	INNER JOIN user_pick_jobs upj ON upj.user_id = u.uid
	INNER JOIN job_categories jc ON jc.uid = upj.job_id
	LEFT JOIN form_biodatas fb ON fb.user_id = p.user_id
	LEFT JOIN form_places fp ON fp.user_id = p.user_id
	LEFT JOIN provinces pro ON pro.id = fp.province_id
	LEFT JOIN regencies reg ON reg.id = fp.city_id
	LEFT JOIN districts dis ON dis.id = fp.district_id
	LEFT JOIN villages vil ON vil.id = fp.subdistrict_id
	WHERE u.uid = '` + p.Id + `'`

	err := db.Debug().Raw(query).Scan(&profiles).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	isProfileExist := len(profiles)

	if isProfileExist == 0 {
		return nil, errors.New("profile not found")
	}

	profile := entities.ProfileResponse{}

	var enabled bool

	if profiles[0].Enabled == 1 {
		enabled = true
	} else {
		enabled = false
	}

	// Education

	var dataEdu = make([]entities.ProfileFormEducation, 0)

	queryEdu := `SELECT id, education_level, major, school_or_college, start_year, start_month, end_month, end_year, user_id 
	FROM form_educations WHERE user_id  = '` + profiles[0].Id + `'`

	rows, errEdu := db.Debug().Raw(queryEdu).Scan(&education).Rows()

	if errEdu != nil {
		helper.Logger("error", "In Server: "+errEdu.Error())
		return nil, errors.New(errEdu.Error())
	}

	defer rows.Close()

	for rows.Next() {
		errScanRows := db.ScanRows(rows, &education)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		queryFormEducationMedia := `SELECT id, path
		FROM form_exercise_medias WHERE exercise_id  = ?`

		rowsFormEducationMedia, errFormExerciseMedia := db.Debug().Raw(queryFormEducationMedia, education.Id).Scan(&educationMedia).Rows()

		if errFormExerciseMedia != nil {
			helper.Logger("error", "In Server: "+errFormExerciseMedia.Error())
			return nil, errors.New(errFormExerciseMedia.Error())
		}

		var dataFormEducationLetter = make([]entities.FormEducationLetter, 0)

		for rowsFormEducationMedia.Next() {
			errScanFormEducationMedia := db.ScanRows(rowsFormEducationMedia, &exerciseMedia)

			if errScanFormEducationMedia != nil {
				helper.Logger("error", "In Server: "+errScanFormEducationMedia.Error())
				return nil, errors.New(errScanFormEducationMedia.Error())
			}

			dataFormEducationLetter = append(dataFormEducationLetter, entities.FormEducationLetter{
				Id:   educationMedia.Id,
				Path: educationMedia.Path,
			})
		}

		dataEdu = append(dataEdu, entities.ProfileFormEducation{
			Id:              education.Id,
			EducationLevel:  education.EducationLevel,
			Major:           education.Major,
			SchoolOrCollege: education.SchoolOrCollege,
			StartYear:       education.StartYear,
			EndYear:         education.EndYear,
			StartMonth:      education.StartMonth,
			EndMonth:        education.EndMonth,
			Letters:         dataFormEducationLetter,
		})
	}

	// Trainning

	var dataTraining = make([]entities.ProfileFormExercise, 0)

	queryTraining := `SELECT id, name, institution, start_month, start_year, end_month, end_year, user_id 
	FROM form_exercises WHERE user_id  = ?`

	rows, errTrainning := db.Debug().Raw(queryTraining, profiles[0].Id).Scan(&exercise).Rows()

	if errTrainning != nil {
		helper.Logger("error", "In Server: "+errTrainning.Error())
		return nil, errors.New(errTrainning.Error())
	}

	defer rows.Close()

	for rows.Next() {
		errScanRows := db.ScanRows(rows, &exercise)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		queryFormExerciseMedia := `SELECT id, path
		FROM form_exercise_medias WHERE exercise_id  = ?`

		rowsFormExerciseMedia, errFormExerciseMedia := db.Debug().Raw(queryFormExerciseMedia, exercise.Id).Scan(&exerciseMedia).Rows()

		if errFormExerciseMedia != nil {
			helper.Logger("error", "In Server: "+errFormExerciseMedia.Error())
			return nil, errors.New(errFormExerciseMedia.Error())
		}

		defer rowsFormExerciseMedia.Close()

		var dataFormExerciseCertificate = make([]entities.FormExerciseCertificate, 0)

		for rowsFormExerciseMedia.Next() {
			errScanFormExerciseMedia := db.ScanRows(rowsFormExerciseMedia, &exerciseMedia)

			if errScanFormExerciseMedia != nil {
				helper.Logger("error", "In Server: "+errScanFormExerciseMedia.Error())
				return nil, errors.New(errScanFormExerciseMedia.Error())
			}

			dataFormExerciseCertificate = append(dataFormExerciseCertificate, entities.FormExerciseCertificate{
				Id:   exerciseMedia.Id,
				Path: exerciseMedia.Path,
			})
		}

		dataTraining = append(dataTraining, entities.ProfileFormExercise{
			Id:           exercise.Id,
			Name:         exercise.Name,
			Institution:  exercise.Institution,
			StartYear:    exercise.StartYear,
			StartMonth:   exercise.StartMonth,
			EndYear:      exercise.EndYear,
			EndMonth:     exercise.EndMonth,
			Certificates: dataFormExerciseCertificate,
		})
	}

	// Work

	var dataWork = make([]entities.ProfileFormWork, 0)

	queryWork := `SELECT id, work, position, institution, is_work, country, city, start_month, start_year, end_month, end_year, user_id 
	FROM form_works WHERE user_id  = '` + profiles[0].Id + `'`

	rows, errWork := db.Debug().Raw(queryWork).Scan(&work).Rows()

	if errWork != nil {
		helper.Logger("error", "In Server: "+errWork.Error())
		return nil, errors.New(errWork.Error())
	}

	defer rows.Close()

	for rows.Next() {
		errScanRows := db.ScanRows(rows, &work)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		var isWork bool

		if work.IsWork == 1 {
			isWork = true
		} else {
			isWork = false
		}

		var dataWork = make([]entities.ProfileFormWork, 0)

		dataWork = append(dataWork, entities.ProfileFormWork{
			Id:          work.Id,
			Position:    work.Position,
			Institution: work.Institution,
			Work:        work.Work,
			IsWork:      isWork,
			City:        work.City,
			Country:     work.Country,
			StartMonth:  work.StartMonth,
			EndMonth:    work.EndMonth,
			StartYear:   work.StartYear,
			EndYear:     work.EndYear,
		})
	}

	// Language

	var dataLanguage = make([]entities.ProfileFormLanguage, 0)

	queryLanguage := `SELECT id, language, level
	FROM form_languages WHERE user_id  = '` + profiles[0].Id + `'`

	rows, errLanguage := db.Debug().Raw(queryLanguage).Scan(&language).Rows()

	if errLanguage != nil {
		helper.Logger("error", "In Server: "+errLanguage.Error())
		return nil, errors.New(errLanguage.Error())
	}

	defer rows.Close()

	for rows.Next() {
		errScanRows := db.ScanRows(rows, &language)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		dataLanguage = append(dataLanguage, language)
	}

	profile.Id = profiles[0].Id
	profile.Avatar = profiles[0].Avatar
	profile.Phone = profiles[0].Phone
	profile.Email = profiles[0].Email
	profile.Fullname = profiles[0].Fullname
	profile.Role = profiles[0].Role
	profile.IsEnabled = enabled
	profile.Job = entities.ProfileJobResponse{
		Id:   profiles[0].JobId,
		Name: profiles[0].JobName,
	}
	profile.Biodata = entities.Biodata{
		Personal: entities.ProfileFormBiodata{
			Id:        profiles[0].BioId,
			Birthdate: profiles[0].BioBirthdate,
			Gender:    profiles[0].BioGender,
			Height:    profiles[0].BioHeight,
			Weight:    profiles[0].BioWeight,
			Religion:  profiles[0].BioReligion,
			Place:     profiles[0].BioPlace,
			Status:    profiles[0].BioStatus,
		},
		Address: entities.ProfileFormPlace{
			Id:            profiles[0].BioAddressId,
			DetailAddress: profiles[0].BioDetailAddress,
			Province: entities.ProfileFormPlaceData{
				Id:   profiles[0].BioProvinceId,
				Name: profiles[0].BioProvince,
			},
			City: entities.ProfileCityPlaceData{
				Id:   profiles[0].BioCityId,
				Name: profiles[0].BioCity,
			},
			District: entities.ProfileDistrictPlaceData{
				Id:   profiles[0].BioDistrictId,
				Name: profiles[0].BioDistrict,
			},
			Subdistrict: entities.ProfileSubdistrictPlaceData{
				Id:   profiles[0].BioSubdistrictId,
				Name: profiles[0].BioSubdistrict,
			},
		},
		Educations:  dataEdu,
		Trainings:   dataTraining,
		Experiences: dataWork,
		Languages:   dataLanguage,
	}

	return map[string]any{
		"data": profile,
	}, nil
}

func UpdateProfile(p *models.Profile) (map[string]interface{}, error) {

	errUpdateProfile := db.Debug().Exec(`
	UPDATE profiles SET fullname = '` + p.Fullname + `', 
	avatar = '` + p.Avatar + `' WHERE user_id = '` + p.Id + `'`).Error

	if errUpdateProfile != nil {
		helper.Logger("error", "In Server: "+errUpdateProfile.Error())
		return nil, errors.New(errUpdateProfile.Error())
	}

	return map[string]any{}, nil
}
