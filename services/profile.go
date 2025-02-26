package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
	models "superapps/models"
)

func GetProfile(p *models.Profile) (map[string]interface{}, error) {
	profiles := []entities.Profile{}

	query := `SELECT p.user_id AS id, p.fullname, p.avatar, u.phone, u.email, 
	jc.uid AS job_id,
	jc.name AS job_name
	FROM profiles p 
	INNER JOIN users u ON u.uid = p.user_id
	INNER JOIN user_pick_jobs upj ON upj.user_id = u.uid
	INNER JOIN job_categories jc ON jc.uid = upj.job_id
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

	profile.Id = profiles[0].Id
	profile.Avatar = profiles[0].Avatar
	profile.Phone = profiles[0].Phone
	profile.Email = profiles[0].Email
	profile.Fullname = profiles[0].Fullname
	profile.Job = entities.ProfileJobResponse{
		Id:   profiles[0].JobId,
		Name: profiles[0].JobName,
	}

	return map[string]any{
		"data": profile,
	}, nil
}
