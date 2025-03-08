package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
	models "superapps/models"
)

func GetProfile(p *models.Profile) (map[string]interface{}, error) {
	profiles := []entities.Profile{}

	query := `SELECT p.user_id AS id, p.fullname, p.avatar, u.phone, u.email, u.enabled, 
	jc.uid AS job_id,
	jc.name AS job_name,
	fb.birthdate AS bio_birthdate,
	fb.gender AS bio_gender,
	fb.weight AS bio_weight,
	fb.height AS bio_height,
	fb.religion AS bio_religion,
	fb.status AS bio_status,
	fb.place AS bio_place,
	fe.education_level AS edu_education_level,
	fe.major AS edu_major,
	fe.school_or_college AS edu_school_or_college, 
	fe.start_year AS edu_start_year,
	fe.end_year AS edu_end_year,
	fe.start_month AS edu_start_month,
	fe.end_month AS edu_end_month
	FROM profiles p 
	INNER JOIN users u ON u.uid = p.user_id
	INNER JOIN user_pick_jobs upj ON upj.user_id = u.uid
	INNER JOIN job_categories jc ON jc.uid = upj.job_id
	LEFT JOIN form_biodata fb ON fb.user_id = p.user_id
	LEFT JOIN form_education fe ON fe.user_id = p.user_id
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

	profile.Id = profiles[0].Id
	profile.Avatar = profiles[0].Avatar
	profile.Phone = profiles[0].Phone
	profile.Email = profiles[0].Email
	profile.Fullname = profiles[0].Fullname
	profile.IsEnabled = enabled
	profile.FormBiodata = entities.ProfileFormBiodata{
		Birthdate: profiles[0].BioBirthdate,
		Gender:    profiles[0].BioGender,
		Height:    profiles[0].BioHeight,
		Weight:    profiles[0].BioWeight,
		Religion:  profiles[0].BioReligion,
		Place:     profiles[0].BioPlace,
		Status:    profiles[0].BioStatus,
	}
	profile.FormEducation = entities.ProfileFormEducation{
		EducationLevel: profiles[0].BioEducationLevel,
		Major: profiles[0].BioMajor,
		SchoolOrCollege: profiles[0].BioSchoolOrCollege,
		StartMonth: profiles[0].BioStartMonth,
		EndMonth: profiles[0].BioEndMonth,
		StartYear: profiles[0].BioStartYear,
		EndYear: profiles[0].BioEndYear
	}
	profile.Job = entities.ProfileJobResponse{
		Id:   profiles[0].JobId,
		Name: profiles[0].JobName,
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
