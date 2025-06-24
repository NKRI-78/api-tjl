package services

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	entities "superapps/entities"
	helper "superapps/helpers"
	models "superapps/models"
)

func DocumentList(d *models.DocumentAssign) (map[string]any, error) {
	docs := []entities.Document{}

	query := `SELECT ud.id, d.name, ud.type, ud.path
	FROM documents d
	INNER JOIN user_documents ud ON ud.type = d.id
	INNER JOIN profiles p ON p.user_id = ud.user_id 
	WHERE ud.user_id = '` + d.UserId + `'`

	err := db.Debug().Raw(query).Scan(&docs).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": docs,
	}, nil
}

func DocumentStore(d *models.DocumentStore) (map[string]any, error) {
	// TYPE MAP
	typeNames := map[string]string{
		"1":  "KTP",
		"2":  "PASPOR",
		"3":  "SKCS",
		"4":  "BIRTH_CERTIFICATE",
		"5":  "GUARDIAN_APPROVAL_LETTER",
		"6":  "FAMILY_CARD",
		"7":  "LAST_DIPLOMA",
		"8":  "EXPERIENCE_CERTIFICATE_WORK",
		"9":  "JOB_SKILLS_CERTIFICATE",
		"10": "MARRIAGE_DOCUMENT",
	}

	// Validasi input
	if d.UserId == "" || d.Path == "" {
		return nil, errors.New("user_id, path, and type cannot be empty")
	}

	// Validasi type
	typeName, ok := typeNames[strconv.Itoa(d.Type)]
	if !ok {
		return nil, fmt.Errorf("invalid document type: %s", d.Type)
	}

	// Check if document exists
	var count int64
	err := db.Table("user_documents").
		Where("user_id = ? AND type = ?", d.UserId, d.Type).
		Count(&count).Error

	if err != nil {
		helper.Logger("error", "DB check failed: "+err.Error())
		return nil, fmt.Errorf("failed to check existing document: %w", err)
	}

	if count > 0 {
		// Document exists → Update
		errUpdate := db.Exec(`
			UPDATE user_documents 
			SET path = ?, updated_at = CURRENT_TIMESTAMP 
			WHERE user_id = ? AND type = ?`,
			d.Path, d.UserId, d.Type).Error

		if errUpdate != nil {
			helper.Logger("error", "Update failed: "+errUpdate.Error())
			return nil, fmt.Errorf("failed to update document: %w", errUpdate)
		}

		return map[string]any{
			"message": fmt.Sprintf("document '%s' updated successfully", typeName),
		}, nil
	}

	// Document doesn't exist → Insert
	errInsert := db.Exec(`
		INSERT INTO user_documents (user_id, path, type) 
		VALUES (?, ?, ?)`, d.UserId, d.Path, d.Type).Error

	if errInsert != nil {
		helper.Logger("error", "Insert failed: "+errInsert.Error())
		return nil, fmt.Errorf("failed to store document: %w", errInsert)
	}

	return map[string]any{
		"message": fmt.Sprintf("document '%s' stored successfully", typeName),
	}, nil
}

func GetDocumentAdditional(userId, typeParam string) (map[string]any, error) {
	var doc entities.GetDocumentAdditional

	row := db.Debug().Raw(`
		SELECT id, type, path
		FROM user_document_additionals
		WHERE user_id = ? AND type = ? 
		ORDER BY created_at DESC`, userId, typeParam).Row()

	err := row.Scan(&doc.Id, &doc.Type, &doc.Path)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return map[string]any{
				"data": map[string]any{},
			}, nil
		}
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	response := entities.GetDocumentAdditionalResponse{
		Id:   doc.Id,
		Path: doc.Path,
		Type: doc.Type,
	}

	return map[string]any{
		"data": response,
	}, nil
}

func DocumentAdditionalStore(d *models.DocumentAdditionalStore) (map[string]any, error) {
	// Check if the document already exists
	var count int64

	row := db.
		Debug().
		Raw(`
		SELECT COUNT(*)
		FROM user_document_additionals
		WHERE user_id = ? AND type = ?`, d.UserId, d.Type).
		Row()

	err := row.Scan(&count)
	if err != nil {
		helper.Logger("error", "Check Existence Failed: "+err.Error())
		return nil, errors.New("failed to check existing document")
	}

	// If exists, return early
	if count > 0 {
		return map[string]any{
			"data": "DOCUMENT_ALREADY_EXIST",
		}, nil
	}

	// Proceed with insert
	err = db.
		Debug().
		Exec(`
		INSERT INTO user_document_additionals (user_id, path, type)
		VALUES (?, ?, ?)`, d.UserId, d.Path, d.Type).
		Error

	if err != nil {
		helper.Logger("error", "Insert Failed: "+err.Error())
		return nil, errors.New("failed to store document")
	}

	return map[string]any{
		"data": "Document stored successfully",
	}, nil

}

func DocumentAdditionalUpdate(d *models.DocumentAdditionalUpdate) (map[string]any, error) {

	doc := entities.DocumentAdditionalUpdate{}

	doc.Path = d.Path
	doc.Type = d.Type

	errInsertAdditionalDoc := db.Debug().Exec(`UPDATE user_document_additionals SET path = ?, type = ? WHERE user_id = ? AND type = ?`, doc.Path, doc.Type, d.UserId, doc.Type).Error

	if errInsertAdditionalDoc != nil {
		helper.Logger("error", "In Server: "+errInsertAdditionalDoc.Error())
		return nil, errors.New(errInsertAdditionalDoc.Error())
	}

	return map[string]any{}, nil
}

func DocumentAdditionalDelete(d *models.DocumentAdditionalDelete) (map[string]any, error) {

	doc := entities.DocumentAdditionalDelete{}

	doc.Id = d.Id

	errDeleteAdditionalDoc := db.Debug().Exec(`DELETE FROM user_document_additionals WHERE id = ?`, doc.Id).Error

	if errDeleteAdditionalDoc != nil {
		helper.Logger("error", "In Server: "+errDeleteAdditionalDoc.Error())
		return nil, errors.New(errDeleteAdditionalDoc.Error())
	}

	return map[string]any{}, nil
}
