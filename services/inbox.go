package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
)

func InboxList(userId string) (map[string]any, error) {
	inboxQuery := entities.InboxListQuery{}
	var data []entities.InboxListResult

	query := `SELECT i.uid AS id, i.title, i.caption, i.is_read, i.field1, i.field2, i.field3, i.field4, i.field5, i.type, i.created_at, p.user_id, p.fullname AS user_fullname, p.avatar AS user_avatar
		FROM inboxes i 
		INNER JOIN profiles p ON p.user_id = i.user_id
		WHERE p.user_id = ?
	`

	rows, err := db.Debug().Raw(query, userId).Rows()

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		errDepartureRows := db.ScanRows(rows, &inboxQuery)

		if errDepartureRows != nil {
			helper.Logger("error", "In Server: "+errDepartureRows.Error())
			return nil, errors.New(errDepartureRows.Error())
		}

		data = append(data, entities.InboxListResult{
			Id:      inboxQuery.Id,
			Title:   inboxQuery.Title,
			Caption: inboxQuery.Caption,
			IsRead:  inboxQuery.IsRead,
			Field1:  inboxQuery.Field1,
			Field2:  inboxQuery.Field2,
			Field3:  inboxQuery.Field3,
			Field4:  inboxQuery.Field4,
			Field5:  inboxQuery.Field5,
			Type:    inboxQuery.Type,
			User: entities.InboxUser{
				Id:       inboxQuery.UserId,
				Fullname: inboxQuery.UserFullname,
				Avatar:   inboxQuery.UserAvatar,
			},
			CreatedAt: inboxQuery.CreatedAt,
		})
	}

	if data == nil {
		data = []entities.InboxListResult{}
	}

	return map[string]any{
		"data": data,
	}, nil
}

func InboxDetail(userId, inboxId string) (map[string]any, error) {
	inboxQuery := entities.InboxListQuery{}

	query := `SELECT i.uid AS id, i.title, i.caption, i.is_read, i.field1, i.field2, i.field3, i.field4, i.field5, i.type, i.created_at, p.user_id, p.fullname AS user_fullname, p.avatar AS user_avatar
		FROM inboxes i 
		INNER JOIN profiles p ON p.user_id = i.user_id
		WHERE i.uid = ? AND p.user_id = ?
		LIMIT 1
	`

	row := db.Debug().Raw(query, inboxId, userId).Row()
	err := row.Scan(
		&inboxQuery.Id,
		&inboxQuery.Title,
		&inboxQuery.Caption,
		&inboxQuery.IsRead,
		&inboxQuery.Field1,
		&inboxQuery.Field2,
		&inboxQuery.Field3,
		&inboxQuery.Field4,
		&inboxQuery.Field5,
		&inboxQuery.Type,
		&inboxQuery.CreatedAt,
		&inboxQuery.UserId,
		&inboxQuery.UserFullname,
		&inboxQuery.UserAvatar,
	)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New("inbox not found")
	}

	updateQuery := `UPDATE inboxes SET is_read = true WHERE uid = ? AND user_id = ?`
	if err := db.Debug().Exec(updateQuery, inboxId, userId).Error; err != nil {
		helper.Logger("error", "Failed to update is_read: "+err.Error())
	}

	detail := entities.InboxListResult{
		Id:      inboxQuery.Id,
		Title:   inboxQuery.Title,
		Caption: inboxQuery.Caption,
		Field1:  inboxQuery.Field1,
		Field2:  inboxQuery.Field2,
		Field3:  inboxQuery.Field3,
		Field4:  inboxQuery.Field4,
		Field5:  inboxQuery.Field5,
		Type:    inboxQuery.Type,
		IsRead:  true,
		User: entities.InboxUser{
			Id:       inboxQuery.UserId,
			Fullname: inboxQuery.UserFullname,
			Avatar:   inboxQuery.UserAvatar,
		},
		CreatedAt: inboxQuery.CreatedAt,
	}

	return map[string]any{
		"data": detail,
	}, nil
}

func InboxBadge(userId string) (map[string]any, error) {
	var count int

	query := `SELECT COUNT(*) FROM inboxes WHERE user_id = ? AND is_read = false`
	err := db.Debug().Raw(query, userId).Row().Scan(&count)

	if err != nil {
		helper.Logger("error", "Failed to get unread inbox count: "+err.Error())
		return nil, errors.New("failed to count unread inboxes")
	}

	return map[string]any{
		"data": count,
	}, nil
}
