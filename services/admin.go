package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
)

func AdminListUser() (map[string]any, error) {
	users := []entities.AdminListUser{}

	query := `SELECT p.user_id AS id, u.email, u.phone, p.avatar, p.fullname, ur.name AS role, u.created_at
	FROM users u 
	INNER JOIN profiles p ON p.user_id = u.uid
	INNER JOIN user_roles ur ON ur.id = u.role 
	`

	err := db.Debug().Raw(query).Scan(&users).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": users,
	}, nil
}
