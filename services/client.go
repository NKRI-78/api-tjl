package services

import (
	"errors"
	"fmt"
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

	return map[string]any{"data": lastID}, nil
}

func ClientUpdate(c *entities.ClientUpdate) (map[string]any, error) {
	query := `UPDATE clients SET icon = ?, link = ?, name = ? WHERE id = ?`

	sqlDB := db.DB()

	result, err := sqlDB.Exec(query, c.Icon, c.Link, c.Name, c.Id)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no client found with ID %d", c.Id)
	}

	return map[string]any{"data": "Ok"}, nil
}

func ClientDelete(c *entities.ClientDelete) (map[string]any, error) {
	query := `DELETE FROM clients WHERE id = ?`

	sqlDB := db.DB()

	result, err := sqlDB.Exec(query, c.Id)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no client found with ID %d", c.Id)
	}

	return map[string]any{"data": "Ok"}, nil
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
