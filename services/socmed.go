package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
)

func ClientStore(n *entities.ClientStore) (map[string]any, error) {
	query := `INSERT INTO clients (icon, link, name) VALUES (?, ?, ?)`

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

func ClientList() (map[string]any, error) {
	clients := []entities.ClientList{}

	query := `SELECT id, icon, link, name FROM clients`

	err := db.Debug().Raw(query).Scan(&clients).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": clients,
	}, nil
}
