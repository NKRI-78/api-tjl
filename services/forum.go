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
	url := os.Getenv("API_URL_DEV")

	var allForum []models.Forum
	var appendForumAssign = make([]entities.ForumResponse, 0)
	var forum entities.Forum
	var forumLike entities.ForumLikeQuery
	var forumLikeAssign entities.ForumLike
	var forumComment entities.ForumCommentQuery
	var forumCommentAssign entities.ForumComment
	var forumCommentReply entities.ForumCommentReplyQuery
	var forumMedia entities.ForumMedia
	var forumMediaAssign entities.ForumMedia
	var user entities.ForumUser
	var userAssign entities.ForumUser

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
	p.fullname, 
	p.user_id,
	u.email,
	u.phone,
	ft.id AS forum_type_id, 
	ft.name AS forum_type_name, 
	f.user_id, f.created_at
	FROM forums f
	INNER JOIN forum_types ft ON ft.id = f.type
	INNER JOIN profiles p ON f.user_id = p.user_id
	INNER JOIN users u ON u.uid = p.user_id
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

		rows, errUser := db.Debug().Raw(`SELECT email, phone, fullname FROM users u 
		INNER JOIN profiles p ON p.user_id = u.uid 
		WHERE u.uid = '` + forum.UserId + `'`).Rows()

		if errUser != nil {
			helper.Logger("error", "In Server: "+errUser.Error())
			return nil, errors.New(errUser.Error())
		}

		for rows.Next() {
			errScanRows := db.ScanRows(rows, &user)

			if errScanRows != nil {
				helper.Logger("error", "In Server: "+errScanRows.Error())
				return nil, errors.New(errScanRows.Error())
			}

			userAssign.Fullname = user.Fullname
			userAssign.Email = user.Email
			userAssign.Phone = user.Phone
		}

		// # ----- forum media ----- # //

		var dataForumMedia = make([]entities.ForumMedia, 0)

		rows, errForumMediaQuery := db.Debug().Raw(`SELECT id, path 
			FROM forum_medias 
			WHERE forum_id = '` + forum.Id + `'`).Rows()

		if errForumMediaQuery != nil {
			helper.Logger("error", "In Server: "+errForumMediaQuery.Error())
			return nil, errors.New(errForumMediaQuery.Error())
		}

		for rows.Next() {
			errScanRows := db.ScanRows(rows, &forumMedia)

			if errScanRows != nil {
				helper.Logger("error", "In Server: "+errScanRows.Error())
				return nil, errors.New(errScanRows.Error())
			}

			forumMediaAssign.Id = forumMedia.Id
			forumMediaAssign.Path = forumMedia.Path

			dataForumMedia = append(dataForumMedia, forumMediaAssign)
		}

		// # CLOSE ----- forum media ----- # //

		// # ----- forum like ----- # //

		var dataForumLike = make([]entities.ForumLike, 0)

		rows, errForumLike := db.Debug().Raw(`SELECT fl.uid AS id,
		p.user_id, p.avatar, p.fullname 
		FROM forum_likes fl 
		INNER JOIN profiles p ON p.user_id = fl.user_id
		WHERE fl.forum_id = '` + forum.Id + `'`).Rows()

		if errForumLike != nil {
			helper.Logger("error", "In Server: "+errForumLike.Error())
			return nil, errors.New(errForumLike.Error())
		}

		for rows.Next() {
			errScanRows := db.ScanRows(rows, &forumLike)

			if errScanRows != nil {
				helper.Logger("error", "In Server: "+errScanRows.Error())
				return nil, errors.New(errScanRows.Error())
			}

			forumLikeAssign.Id = forumLike.Id
			forumLikeAssign.User = entities.ForumLikeUser{
				Id:       forumLike.UserId,
				Avatar:   forumLike.Avatar,
				Fullname: forumLike.Fullname,
			}

			dataForumLike = append(dataForumLike, forumLikeAssign)
		}

		// # CLOSE ----- forum like ----- # //

		// # ----- forum comment ----- # //

		var dataForumComment = make([]entities.ForumComment, 0)

		rows, errForumComment := db.Debug().Raw(`SELECT fc.uid AS id, p.avatar, fc.comment, p.user_id, p.fullname 
		FROM forum_comments fc
		INNER JOIN profiles p ON p.user_id = fc.user_id
		WHERE fc.forum_id = '` + forum.Id + `'`).Rows()

		if errForumComment != nil {
			helper.Logger("error", "In Server: "+errForumComment.Error())
			return nil, errors.New(errForumComment.Error())
		}

		for rows.Next() {
			errScanRows := db.ScanRows(rows, &forumComment)

			if errScanRows != nil {
				helper.Logger("error", "In Server: "+errScanRows.Error())
				return nil, errors.New(errScanRows.Error())
			}

			var dataForumCommentReply = make([]entities.ForumCommentReply, 0)

			rows, errForumCommentReply := db.Debug().Raw(`SELECT fcr.uid AS id, fcr.reply, p.avatar, p.user_id, p.fullname 
			FROM forum_comment_replies fcr
			INNER JOIN profiles p ON p.user_id = fcr.user_id
			WHERE fcr.comment_id = '` + forumComment.Id + `'`).Rows()

			for rows.Next() {
				errScanRows := db.ScanRows(rows, &forumCommentReply)

				if errScanRows != nil {
					helper.Logger("error", "In Server: "+errScanRows.Error())
					return nil, errors.New(errScanRows.Error())
				}

				dataForumCommentReply = append(dataForumCommentReply, entities.ForumCommentReply{
					Id:    forumCommentReply.Id,
					Reply: forumCommentReply.Reply,
					User: entities.ForumCommentReplyUser{
						Id:       forumCommentReply.Id,
						Avatar:   forumCommentReply.Avatar,
						Fullname: forumCommentReply.Fullname,
					},
				})
			}

			if errForumCommentReply != nil {
				helper.Logger("error", "In Server: "+errForumComment.Error())
				return nil, errors.New(errForumComment.Error())
			}

			forumCommentAssign.Id = forumComment.Id
			forumCommentAssign.Comment = forumComment.Comment
			forumCommentAssign.Reply = dataForumCommentReply
			forumCommentAssign.ReplyCount = len(dataForumCommentReply)
			forumCommentAssign.User = entities.ForumCommentUser{
				Id:       forumComment.UserId,
				Avatar:   forumComment.Avatar,
				Fullname: forumComment.Fullname,
			}

			dataForumComment = append(dataForumComment, forumCommentAssign)
		}

		// # CLOSE ----- forum comment ----- # //

		appendForumAssign = append(appendForumAssign, entities.ForumResponse{
			Id:           forum.Id,
			Title:        forum.Title,
			Caption:      forum.Caption,
			Media:        dataForumMedia,
			Comment:      dataForumComment,
			CommentCount: len(dataForumComment),
			Like:         dataForumLike,
			LikeCount:    len(dataForumLike),
			ForumType: entities.ForumType{
				Id:   forum.ForumTypeId,
				Name: forum.ForumTypeName,
			},
			User: entities.ForumUser{
				Id:       forum.UserId,
				Fullname: forum.Fullname,
				Email:    forum.Email,
				Phone:    forum.Phone,
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
		"next_url":     url + "?page=" + nextUrl,
		"prev_url":     url + "?page=" + prevUrl,
		"data":         &appendForumAssign,
	}, nil
}

func ForumDetail(f *models.Forum) (map[string]interface{}, error) {
	var appendForumAssign = make([]entities.ForumResponse, 0)
	var forumComment entities.ForumCommentQuery
	var forumCommentReply entities.ForumCommentReplyQuery
	var forumMedia entities.ForumMedia
	var forumMediaAssign entities.ForumMedia
	var forumLike entities.ForumLikeQuery
	var forumLikeAssign entities.ForumLike
	var forumCommentAssign entities.ForumComment
	var user entities.ForumUser
	var userAssign entities.ForumUser

	forum := []entities.Forum{}

	errForum := db.Debug().Raw(`SELECT f.uid AS id, f.title, f.caption,
	p.fullname, 
	p.user_id,
	u.email,
	u.phone,
	ft.id AS forum_type_id, 
	ft.name AS forum_type_name, 
	f.user_id, f.created_at
	FROM forums f
	INNER JOIN forum_types ft ON ft.id = f.type
	INNER JOIN profiles p ON f.user_id = p.user_id
	INNER JOIN users u ON u.uid = p.user_id
	WHERE f.uid = '` + f.Id + `'`).Scan(&forum).Error

	if errForum != nil {
		helper.Logger("error", "In Server: "+errForum.Error())
		return nil, errors.New(errForum.Error())
	}

	forums := len(forum)

	if forums == 0 {
		helper.Logger("error", "In Server: Forum type not found")
		return nil, errors.New("forum not found")
	}

	rows, errUser := db.Debug().Raw(`SELECT email, phone, fullname FROM users u 
	INNER JOIN profiles p ON p.user_id = u.uid 
	WHERE u.uid = '` + forum[0].UserId + `'`).Rows()

	if errUser != nil {
		helper.Logger("error", "In Server: "+errUser.Error())
		return nil, errors.New(errUser.Error())
	}

	for rows.Next() {
		errScanRows := db.ScanRows(rows, &user)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		userAssign.Fullname = user.Fullname
		userAssign.Email = user.Email
		userAssign.Phone = user.Phone
	}

	// # ----- forum media ----- # //

	var dataForumMedia = make([]entities.ForumMedia, 0)

	rows, errForumMediaQuery := db.Debug().Raw(`SELECT id, path 
		FROM forum_medias 
		WHERE forum_id = '` + forum[0].Id + `'`).Rows()

	if errForumMediaQuery != nil {
		helper.Logger("error", "In Server: "+errForumMediaQuery.Error())
		return nil, errors.New(errForumMediaQuery.Error())
	}

	for rows.Next() {
		errScanRows := db.ScanRows(rows, &forumMedia)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		forumMediaAssign.Id = forumMedia.Id
		forumMediaAssign.Path = forumMedia.Path

		dataForumMedia = append(dataForumMedia, forumMediaAssign)
	}

	// # CLOSE ----- forum media ----- # //

	// # ----- forum like ----- # //

	var dataForumLike = make([]entities.ForumLike, 0)

	rows, errForumLike := db.Debug().Raw(`SELECT fl.uid AS id,
	p.user_id, p.avatar, p.fullname 
	FROM forum_likes fl 
	INNER JOIN profiles p ON p.user_id = fl.user_id
	WHERE fl.forum_id = '` + forum[0].Id + `'`).Rows()

	if errForumLike != nil {
		helper.Logger("error", "In Server: "+errForumLike.Error())
		return nil, errors.New(errForumLike.Error())
	}

	for rows.Next() {
		errScanRows := db.ScanRows(rows, &forumLike)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		forumLikeAssign.Id = forumLike.Id
		forumLikeAssign.User = entities.ForumLikeUser{
			Id:       forumLike.UserId,
			Avatar:   forumLike.Avatar,
			Fullname: forumLike.Fullname,
		}

		dataForumLike = append(dataForumLike, forumLikeAssign)
	}

	// # CLOSE ----- forum like ----- # //

	// # ----- forum comment ----- # //

	var dataForumComment = make([]entities.ForumComment, 0)

	rows, errForumComment := db.Debug().Raw(`SELECT fc.uid AS id, p.avatar, fc.comment, p.user_id, p.fullname 
	FROM forum_comments fc
	INNER JOIN profiles p ON p.user_id = fc.user_id
	WHERE fc.forum_id = '` + forum[0].Id + `'`).Rows()

	if errForumComment != nil {
		helper.Logger("error", "In Server: "+errForumComment.Error())
		return nil, errors.New(errForumComment.Error())
	}

	for rows.Next() {
		errScanRows := db.ScanRows(rows, &forumComment)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		var dataForumCommentReply = make([]entities.ForumCommentReply, 0)

		rows, errForumCommentReply := db.Debug().Raw(`SELECT fcr.uid AS id, fcr.reply, p.avatar, p.user_id, p.fullname 
		FROM forum_comment_replies fcr
		INNER JOIN profiles p ON p.user_id = fcr.user_id
		WHERE fcr.comment_id = '` + forumComment.Id + `'`).Rows()

		for rows.Next() {
			errScanRows := db.ScanRows(rows, &forumCommentReply)

			if errScanRows != nil {
				helper.Logger("error", "In Server: "+errScanRows.Error())
				return nil, errors.New(errScanRows.Error())
			}

			dataForumCommentReply = append(dataForumCommentReply, entities.ForumCommentReply{
				Id:    forumCommentReply.Id,
				Reply: forumCommentReply.Reply,
				User: entities.ForumCommentReplyUser{
					Id:       forumCommentReply.Id,
					Avatar:   forumCommentReply.Avatar,
					Fullname: forumCommentReply.Fullname,
				},
			})
		}

		if errForumCommentReply != nil {
			helper.Logger("error", "In Server: "+errForumComment.Error())
			return nil, errors.New(errForumComment.Error())
		}

		forumCommentAssign.Id = forumComment.Id
		forumCommentAssign.Comment = forumComment.Comment
		forumCommentAssign.Reply = dataForumCommentReply
		forumCommentAssign.ReplyCount = len(dataForumCommentReply)
		forumCommentAssign.User = entities.ForumCommentUser{
			Id:       forumComment.UserId,
			Avatar:   forumComment.Avatar,
			Fullname: forumComment.Fullname,
		}

		dataForumComment = append(dataForumComment, forumCommentAssign)
	}

	// # CLOSE ----- forum comment ----- # //

	appendForumAssign = append(appendForumAssign, entities.ForumResponse{
		Id:           forum[0].Id,
		Title:        forum[0].Title,
		Caption:      forum[0].Caption,
		Media:        dataForumMedia,
		Comment:      dataForumComment,
		CommentCount: len(dataForumComment),
		Like:         dataForumLike,
		LikeCount:    len(dataForumLike),
		ForumType: entities.ForumType{
			Id:   forum[0].ForumTypeId,
			Name: forum[0].ForumTypeName,
		},
		User: entities.ForumUser{
			Id:       forum[0].UserId,
			Fullname: forum[0].Fullname,
			Email:    forum[0].Email,
			Phone:    forum[0].Phone,
		},
	})

	return map[string]any{
		"data": &appendForumAssign[0],
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
