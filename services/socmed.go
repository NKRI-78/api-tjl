package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
)

func SocmedStore(n *entities.SocialMediaStore) (map[string]any, error) {
	query := `INSERT INTO medias (icon, link, name) VALUES (?, ?, ?)`

	sqlDB := db.DB()

	result, err := sqlDB.Exec(query, n.Icon, n.Link, n.Name)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return map[string]any{"id": lastID}, nil
}

func SocmedList() (map[string]any, error) {
	socmeds := []entities.SocialMediaList{}

	query := `SELECT id, icon, link, name FROM medias`

	err := db.Debug().Raw(query).Scan(&socmeds).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": socmeds,
	}, nil
}
