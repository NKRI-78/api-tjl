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

func NewsList(page, limit string) (map[string]any, error) {
	url := os.Getenv("API_URL_PROD")

	var allNews []entities.AllNews
	var appendNewsAssign = make([]entities.NewsResponse, 0)
	var newsMedia entities.NewsMedia
	var newsMediaAssign entities.NewsMedia
	var news entities.News

	pageinteger, _ := strconv.Atoi(page)
	limitinteger, _ := strconv.Atoi(limit)

	var offset = strconv.Itoa((pageinteger - 1) * limitinteger)

	errAllForum := db.Debug().Raw(`SELECT id FROM news`).Scan(&allNews).Error

	if errAllForum != nil {
		helper.Logger("error", "In Server: "+errAllForum.Error())
	}

	var resultTotal = len(allNews)

	var perPage = math.Ceil(float64(resultTotal) / float64(limitinteger))

	var prevPage int
	var nextPage int

	if pageinteger == 1 {
		prevPage = 1
	} else {
		prevPage = pageinteger - 1
	}

	nextPage = pageinteger + 1

	rows, errNews := db.Debug().Raw(`
	SELECT n.id, n.title, n.caption,
		   p.fullname AS user_name,
	       n.user_id, n.created_at
	FROM news n
	INNER JOIN profiles p ON n.user_id = p.user_id
	INNER JOIN users u ON u.uid = p.user_id
	LIMIT ?, ?`, offset, limit).Rows()

	if errNews != nil {
		helper.Logger("error", "In Server: "+errNews.Error())
		return nil, errors.New(errNews.Error())
	}

	for rows.Next() {
		errNewsRows := db.ScanRows(rows, &news)

		if errNewsRows != nil {
			helper.Logger("error", "In Server: "+errNewsRows.Error())
			return nil, errors.New(errNewsRows.Error())
		}

		// # ----- news media ----- # //

		var dataNewsMedia = make([]entities.NewsMedia, 0)

		rowsNewsMedia, errNewsMediaQuery := db.Debug().Raw(`SELECT id, path 
			FROM news_medias 
			WHERE news_id = ?`, news.Id).Rows()

		if errNewsMediaQuery != nil {
			helper.Logger("error", "In Server: "+errNewsMediaQuery.Error())
			return nil, errors.New(errNewsMediaQuery.Error())
		}

		for rowsNewsMedia.Next() {
			errScanRows := db.ScanRows(rowsNewsMedia, &newsMedia)

			if errScanRows != nil {
				helper.Logger("error", "In Server: "+errScanRows.Error())
				return nil, errors.New(errScanRows.Error())
			}

			newsMediaAssign.Id = newsMedia.Id
			newsMediaAssign.Path = newsMedia.Path

			dataNewsMedia = append(dataNewsMedia, newsMediaAssign)
		}

		// # CLOSE ----- news media ----- # //

		appendNewsAssign = append(appendNewsAssign, entities.NewsResponse{
			Id:        news.Id,
			Title:     news.Title,
			Caption:   news.Caption,
			Media:     dataNewsMedia,
			CreatedAt: news.CreatedAt,
			User: entities.NewsUser{
				Id:   news.UserId,
				Name: news.UserName,
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
		"next_url":     url + "/api/v1/news?page=" + nextUrl,
		"prev_url":     url + "/api/v1/news?page=" + prevUrl,
		"data":         &appendNewsAssign,
	}, nil
}

func NewsDetail(id string) (map[string]any, error) {
	var news entities.News
	var newsMedia entities.NewsMedia
	var newsMediaAssign entities.NewsMedia
	var newsMediaList = make([]entities.NewsMedia, 0)

	// Fetch the news
	row := db.Debug().Raw(`
		SELECT n.id, n.title, n.caption,
		       p.fullname AS user_name,
		       n.user_id, n.created_at
		FROM news n
		INNER JOIN profiles p ON n.user_id = p.user_id
		INNER JOIN users u ON u.uid = p.user_id
		WHERE n.id = ?`, id).Row()

	err := row.Scan(&news.Id, &news.Title, &news.Caption, &news.UserName, &news.UserId, &news.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return map[string]any{
				"data": nil,
			}, nil
		}
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	// Fetch media
	rowsNewsMedia, errNewsMediaQuery := db.Debug().Raw(`
		SELECT id, path FROM news_medias WHERE news_id = ?`, news.Id).Rows()

	if errNewsMediaQuery != nil {
		helper.Logger("error", "In Server: "+errNewsMediaQuery.Error())
		return nil, errNewsMediaQuery
	}
	defer rowsNewsMedia.Close()

	for rowsNewsMedia.Next() {
		errScanRows := db.ScanRows(rowsNewsMedia, &newsMedia)
		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errScanRows
		}

		newsMediaAssign.Id = newsMedia.Id
		newsMediaAssign.Path = newsMedia.Path
		newsMediaList = append(newsMediaList, newsMediaAssign)
	}

	// Prepare response
	response := entities.NewsResponse{
		Id:        news.Id,
		Title:     news.Title,
		Caption:   news.Caption,
		Media:     newsMediaList,
		CreatedAt: news.CreatedAt,
		User: entities.NewsUser{
			Id:   news.UserId,
			Name: news.UserName,
		},
	}

	return map[string]any{
		"data": response,
	}, nil
}

func NewsStore(n *entities.NewsStore) (map[string]any, error) {
	query := `INSERT INTO news (title, caption, user_id) VALUES (?, ?, ?)`

	sqlDB := db.DB()

	result, err := sqlDB.Exec(query, n.Title, n.Caption, n.UserId)
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

func NewsUpdate(n *entities.NewsUpdate) (map[string]any, error) {
	query := `UPDATE news SET title = ?, caption = ? WHERE id = ?`

	err := db.Debug().Exec(query, n.Title, n.Caption, n.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func NewsStoreImage(n *entities.NewsStoreImage) (map[string]any, error) {
	query := `INSERT INTO news_medias (news_id, path) VALUES (?, ?)`

	err := db.Debug().Exec(query, n.NewsId, n.Path).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func NewsUpdateImage(n *entities.NewsUpdateImage) (map[string]any, error) {
	query := `UPDATE news_medias SET path = ? WHERE id = ?`

	err := db.Debug().Exec(query, n.Path, n.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func NewsDelete(n *entities.NewsDelete) (map[string]any, error) {
	query := `DELETE FROM news WHERE id = ?`

	err := db.Debug().Exec(query, n.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func NewsDeleteImage(n *entities.NewsDelete) (map[string]any, error) {
	query := `DELETE FROM news_medias WHERE id = ?`

	err := db.Debug().Exec(query, n.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}
