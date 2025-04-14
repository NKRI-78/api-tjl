package services

import (
	"errors"

	helper "superapps/helpers"
	models "superapps/models"
)

func InitFcm(d *models.InitFcm) (map[string]any, error) {
	query := `INSERT INTO fcms (token, user_id) VALUES (?, ?) 
	ON DUPLICATE KEY UPDATE token = ?`

	err := db.Debug().Exec(query, d.Token, d.UserId, d.Token).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}
