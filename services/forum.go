package services

import (
	"errors"
	"strconv"
	entities "superapps/entities"
	helper "superapps/helpers"
	models "superapps/models"

	uuid "github.com/satori/go.uuid"
)

func ForumCategory() (map[string]interface{}, error) {
	forumCategory := []entities.ForumCategory{}

	query := `SELECT id, name FROM forum_types`

	err := db.Debug().Raw(query).Scan(&forumCategory).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	isForumCategoryExist := len(forumCategory)

	if isForumCategoryExist == 0 {
		return nil, errors.New("forum type not found")
	}

	return map[string]any{}, nil
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
