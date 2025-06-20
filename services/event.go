package services

import (
	"database/sql"
	"errors"
	"math"
	"os"
	"strconv"
	entities "superapps/entities"
	helper "superapps/helpers"
)

func EventList(page, limit, startDate string) (map[string]any, error) {
	url := os.Getenv("API_URL_PROD")

	var allEvent []entities.AllEvent
	var appendEventAssign = make([]entities.EventResponse, 0)
	var eventMedia entities.EventMedia
	var eventMediaAssign entities.EventMedia
	var event entities.Event

	pageinteger, _ := strconv.Atoi(page)
	limitinteger, _ := strconv.Atoi(limit)

	var offset = strconv.Itoa((pageinteger - 1) * limitinteger)

	// Count total events based on createdAt filter
	var countQuery string
	var countArgs []any

	if startDate != "" {
		countQuery = `SELECT id FROM events WHERE start_date = ?`
		countArgs = append(countArgs, startDate)
	} else {
		countQuery = `SELECT id FROM events`
	}

	errAllEvent := db.Debug().Raw(countQuery, countArgs...).Scan(&allEvent).Error
	if errAllEvent != nil {
		helper.Logger("error", "In Server: "+errAllEvent.Error())
	}

	var resultTotal = len(allEvent)
	var perPage = math.Ceil(float64(resultTotal) / float64(limitinteger))

	var prevPage int
	var nextPage int

	if pageinteger == 1 {
		prevPage = 1
	} else {
		prevPage = pageinteger - 1
	}

	nextPage = pageinteger + 1

	// Main query with optional createdAt filter
	var rows *sql.Rows
	var errEvent error

	if startDate != "" {
		rows, errEvent = db.Debug().Raw(`
			SELECT e.id, e.title, e.caption,
			p.fullname AS user_name,
			e.user_id, e.created_at,
			e.start_time, 
			e.end_time, 
			e.start_date, 
			e.end_date
			FROM events e
			INNER JOIN profiles p ON e.user_id = p.user_id
			INNER JOIN users u ON u.uid = p.user_id
			WHERE e.start_date = ?
			LIMIT ?, ?`, startDate, offset, limit).Rows()
	} else {
		rows, errEvent = db.Debug().Raw(`
			SELECT e.id, e.title, e.caption,
			p.fullname AS user_name,
			e.user_id, e.created_at,
			e.start_time, 
			e.end_time, 
			e.start_date, 
			e.end_date
			FROM events e
			INNER JOIN profiles p ON e.user_id = p.user_id
			INNER JOIN users u ON u.uid = p.user_id
			LIMIT ?, ?`, offset, limit).Rows()
	}

	if errEvent != nil {
		helper.Logger("error", "In Server: "+errEvent.Error())
		return nil, errors.New(errEvent.Error())
	}

	for rows.Next() {
		errEventRows := db.ScanRows(rows, &event)
		if errEventRows != nil {
			helper.Logger("error", "In Server: "+errEventRows.Error())
			return nil, errors.New(errEventRows.Error())
		}

		// Event Media
		var dataEventMedia = make([]entities.EventMedia, 0)

		rowsEventMedia, errEventMediaQuery := db.Debug().Raw(`SELECT id, path 
			FROM event_medias 
			WHERE event_id = ?`, event.Id).Rows()

		if errEventMediaQuery != nil {
			helper.Logger("error", "In Server: "+errEventMediaQuery.Error())
			return nil, errors.New(errEventMediaQuery.Error())
		}

		for rowsEventMedia.Next() {
			errScanRows := db.ScanRows(rowsEventMedia, &eventMedia)
			if errScanRows != nil {
				helper.Logger("error", "In Server: "+errScanRows.Error())
				return nil, errors.New(errScanRows.Error())
			}
			eventMediaAssign.Id = eventMedia.Id
			eventMediaAssign.Path = eventMedia.Path
			dataEventMedia = append(dataEventMedia, eventMediaAssign)
		}

		appendEventAssign = append(appendEventAssign, entities.EventResponse{
			Id:        event.Id,
			Title:     event.Title,
			Caption:   event.Caption,
			Media:     dataEventMedia,
			StartDate: event.StartDate,
			EndDate:   event.EndDate,
			StartTime: event.StartTime,
			EndTime:   event.EndTime,
			CreatedAt: event.CreatedAt,
			User: entities.EventUser{
				Id:   event.UserId,
				Name: event.UserName,
			},
		})
	}

	var nextUrl = strconv.Itoa(nextPage)
	var prevUrl = strconv.Itoa(prevPage)

	return map[string]any{
		"total":        resultTotal,
		"current_page": pageinteger,
		"per_page":     int(perPage),
		"prev_page":    prevPage,
		"next_page":    nextPage,
		"next_url":     url + "/api/v1/event?page=" + nextUrl + "&start_date=" + startDate,
		"prev_url":     url + "/api/v1/event?page=" + prevUrl + "&start_date=" + startDate,
		"data":         &appendEventAssign,
	}, nil
}

func EventStoreImage(e *entities.EventStoreImage) (map[string]any, error) {
	query := `INSERT INTO event_medias (event_id, path) VALUES (?, ?)`

	err := db.Debug().Exec(query, e.EventId, e.Path).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func EventUpdateImage(n *entities.EventUpdateImage) (map[string]any, error) {
	query := `UPDATE event_medias SET path = ? WHERE id = ?`

	err := db.Debug().Exec(query, n.Path, n.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func EventDeleteImage(n *entities.EventDelete) (map[string]any, error) {
	query := `DELETE FROM event_medias WHERE id = ?`

	err := db.Debug().Exec(query, n.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func EventDetail(id string) (map[string]any, error) {
	var event entities.Event
	var eventMedia entities.EventMedia
	var eventMediaAssign entities.EventMedia
	var eventMediaList = make([]entities.EventMedia, 0)

	row := db.Debug().Raw(`
		SELECT e.id, e.title, e.caption,
		p.fullname AS user_name,
		e.start_date, e.end_date,
		e.start_time, e.end_time,
		e.user_id, e.created_at
		FROM events e
		INNER JOIN profiles p ON e.user_id = p.user_id
		INNER JOIN users u ON u.uid = p.user_id
		WHERE e.id = ?`, id).Row()

	err := row.Scan(&event.Id, &event.Title, &event.Caption, &event.UserName, &event.StartDate, &event.EndDate, &event.StartTime, &event.EndTime, &event.UserId, &event.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return map[string]any{
				"data": nil,
			}, nil
		}
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	rowsEventMedia, errEventMediaQuery := db.Debug().Raw(`
		SELECT id, path FROM event_medias WHERE event_id = ?`, event.Id).Rows()

	if errEventMediaQuery != nil {
		helper.Logger("error", "In Server: "+errEventMediaQuery.Error())
		return nil, errEventMediaQuery
	}
	defer rowsEventMedia.Close()

	for rowsEventMedia.Next() {
		errScanRows := db.ScanRows(rowsEventMedia, &eventMedia)
		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errScanRows
		}

		eventMediaAssign.Id = eventMedia.Id
		eventMediaAssign.Path = eventMedia.Path
		eventMediaList = append(eventMediaList, eventMediaAssign)
	}

	response := entities.EventResponse{
		Id:        event.Id,
		Title:     event.Title,
		StartDate: event.StartDate,
		EndDate:   event.EndDate,
		StartTime: event.StartTime,
		EndTime:   event.EndTime,
		Caption:   event.Caption,
		Media:     eventMediaList,
		CreatedAt: event.CreatedAt,
		User: entities.EventUser{
			Id:   event.UserId,
			Name: event.UserName,
		},
	}

	return map[string]any{
		"data": response,
	}, nil
}

func EventDelete(e *entities.EventDelete) (map[string]any, error) {
	query := `DELETE FROM events WHERE id = ?`

	err := db.Debug().Exec(query, e.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func EventUpdate(n *entities.EventUpdate) (map[string]any, error) {
	query := `UPDATE events SET title = ?, caption = ?, start_date = ?, end_date = ?, start_time = ?, end_time = ? WHERE id = ?`

	err := db.Debug().Exec(query, n.Title, n.Caption, n.StartDate, n.EndDate, n.StartTime, n.EndTime, n.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func EventStore(e *entities.EventStore) (map[string]any, error) {
	query := `INSERT INTO events (title, caption, user_id, start_date, end_date, start_time, end_time) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	sqlDB := db.DB()

	result, err := sqlDB.Exec(query, e.Title, e.Caption, e.UserId, e.StartDate, e.EndDate, e.StartTime, e.EndTime)
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
