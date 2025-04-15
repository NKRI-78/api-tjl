package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
)

func IconList() (map[string]any, error) {
	icons := []entities.IconList{}

	query := `SELECT id, path FROM icons`

	err := db.Debug().Raw(query).Scan(&icons).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": icons,
	}, nil
}

func IconStore(i *entities.IconStore) (map[string]any, error) {
	icon := entities.IconStore{}

	icon.Path = i.Path

	err := db.Debug().Exec(`
	INSERT INTO icons (path) 
	VALUES (?)`, icon.Path).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func IconUpdate(i *entities.IconUpdate) (map[string]any, error) {
	icon := entities.IconUpdate{}

	icon.Id = i.Id
	icon.Path = i.Path

	err := db.Debug().Exec(`
		UPDATE icons SET path = ?
		WHERE id = ?
	`, icon.Path, icon.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func IconDelete(i *entities.IconDelete) (map[string]any, error) {
	icon := entities.IconDelete{}

	icon.Id = i.Id

	err := db.Debug().Exec(`
		DELETE FROM icons WHERE id = ?
	`, icon.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}
