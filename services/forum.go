package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
)

func GetForumCategory() (map[string]interface{}, error) {
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

	return map[string]any{
		"data": forumCategory,
	}, nil
}
