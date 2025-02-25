package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
	models "superapps/models"
)

func GetProfile(p *models.Profile) (map[string]interface{}, error) {

	profiles := []entities.Profile{}

	query := `SELECT p.user_id AS id, p.fullname, p.avatar, u.phone, u.email 
	FROM profiles p 
	INNER JOIN users u ON u.uid = p.user_id
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

	profile := profiles[0]

	return map[string]any{
		"data": profile,
	}, nil
}
