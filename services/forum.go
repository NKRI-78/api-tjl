package services

import (
	"errors"
	"math"
	"os"
	"strconv"
	entities "superapps/entities"
	helper "superapps/helpers"
	models "superapps/models"

	uuid "github.com/satori/go.uuid"
)

func ForumList(search, page, limit string) (map[string]interface{}, error) {
	url := os.Getenv("API_URL")

	println(search)
	print(page)
	println(limit)

	var allForum []models.Forum
	var forum entities.Forum

	pageinteger, _ := strconv.Atoi(page)
	limitinteger, _ := strconv.Atoi(limit)

	var offset = strconv.Itoa((pageinteger - 1) * limitinteger)

	errAllForum := db.Debug().Raw(`SELECT uid FROM forums`).Scan(&allForum).Error

	if errAllForum != nil {
		helper.Logger("error", "In Server: "+errAllForum.Error())
	}

	var resultTotal = len(allForum)

	var perPage = math.Ceil(float64(resultTotal) / float64(limitinteger))

	var prevPage int
	var nextPage int

	if pageinteger == 1 {
		prevPage = 1
	} else {
		prevPage = pageinteger - 1
	}

	nextPage = pageinteger + 1

	rows, errForum := db.Debug().Raw(`SELECT f.uid AS id, f.title, f.caption,
	ft.id AS forum_type_id, 
	ft.name AS forum_type_name, 
	f.user_id, f.created_at
	FROM forums f
	INNER JOIN forum_types ft ON ft.id = f.type
	WHERE f.title LIKE '%` + search + `%'
	LIMIT ` + offset + `, ` + limit + ``).Rows()

	if errForum != nil {
		helper.Logger("error", "In Server: "+errForum.Error())
		return nil, errors.New(errForum.Error())
	}

	for rows.Next() {
		errForumRows := db.ScanRows(rows, &forum)

		if errForumRows != nil {
			helper.Logger("error", "In Server: "+errForumRows.Error())
			return nil, errors.New(errForumRows.Error())
		}

		
	}

	var nextUrl = strconv.Itoa(nextPage)
	var prevUrl = strconv.Itoa(prevPage)

	return map[string]any{
		"total":        resultTotal,
		"current_page": pageinteger,
		"per_page":     int(perPage),
		"prev_page":    prevPage,
		"next_page":    nextPage,
		"next_url":     url + "?page=" + nextUrl,
		"prev_url":     url + "?page=" + prevUrl,
	}, nil
}

func ForumCategory() (map[string]interface{}, error) {
	forumCategory := []entities.ForumCategory{}

	err := db.Debug().Raw(`SELECT id, name FROM forum_types`).Scan(&forumCategory).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	isForumCategoryExist := len(forumCategory)

	if isForumCategoryExist == 0 {
		return nil, errors.New("forum type not found")
	}

	return map[string]any{
		"data": forumCategory,
	}, nil
}

func ForumStore(f *models.ForumStore) (map[string]interface{}, error) {
	forum := entities.ForumStore{}
	forumTypes := []entities.ForumCategory{}

	forum.Id = uuid.NewV4().String()

	forum.Title = f.Title
	forum.Desc = f.Desc
	forum.Type = f.Type

	errCheckForumCategory := db.Debug().Raw(`SELECT id FROM forum_types WHERE id = '` + strconv.Itoa(forum.Type) + `'`).Scan(&forumTypes).Error

	if errCheckForumCategory != nil {
		helper.Logger("error", "In Server: "+errCheckForumCategory.Error())
		return nil, errors.New(errCheckForumCategory.Error())
	}

	isForumTypeExist := len(forumTypes)

	if isForumTypeExist == 0 {
		helper.Logger("error", "In Server: Forum type not found")
		return nil, errors.New("forum type not found")
	}

	errInsertForum := db.Debug().Exec(`
		INSERT INTO forums (uid, title, caption, user_id, type) 
		VALUES (?, ?, ?, ?, ?)`,
		forum.Id, forum.Title, forum.Desc, f.UserId, strconv.Itoa(forum.Type)).Error

	if errInsertForum != nil {
		helper.Logger("error", "In Server: "+errInsertForum.Error())
		return nil, errors.New(errInsertForum.Error())
	}

	return map[string]any{}, nil
}

func ForumDelete(f *models.Forum) (map[string]interface{}, error) {
	forum := entities.Forum{}
	forums := []entities.Forum{}

	forum.Id = f.Id

	errCheckForum := db.Debug().Raw(`SELECT uid AS id FROM forums WHERE uid = '` + forum.Id + `'`).Scan(&forums).Error

	if errCheckForum != nil {
		helper.Logger("error", "In Server: "+errCheckForum.Error())
		return nil, errors.New(errCheckForum.Error())
	}

	isForumExist := len(forums)

	if isForumExist == 0 {
		helper.Logger("error", "Forum not found")
		return nil, errors.New("forum not found")
	}

	errDeleteForum := db.Debug().Exec(`DELETE FROM forums WHERE uid = '` + forum.Id + `'`).Error

	if errDeleteForum != nil {
		helper.Logger("error", "In Server: "+errDeleteForum.Error())
		return nil, errors.New(errDeleteForum.Error())
	}

	return map[string]any{}, nil
}
