package services

import (
	"errors"

	entities "superapps/entities"
	helper "superapps/helpers"
	models "superapps/models"
)

func DocumentStore(d *models.DocumentStore) (map[string]any, error) {
	doc := entities.DocumentStore{}

	doc.Path = d.Path
	doc.UserId = d.UserId
	doc.Type = d.Type

	errInsertDoc := db.Debug().Exec(`
	INSERT INTO user_documents (user_id, path, type) 
	VALUES (?, ?, ?)`, doc.UserId, doc.Path, doc.Type).Error

	if errInsertDoc != nil {
		helper.Logger("error", "In Server: "+errInsertDoc.Error())
		return nil, errors.New(errInsertDoc.Error())
	}

	return map[string]any{}, nil
}