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

func ForumList(userId, search, page, limit string) (map[string]any, error) {
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
	p.avatar,
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
	ORDER BY f.id DESC
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

		// Check Forum is Like

		checkForumIsLike := []entities.CheckIsLike{}

		errCheckForumIsLike := db.Debug().Raw(`SELECT EXISTS (
			SELECT 1
			FROM forum_likes
			WHERE user_id = ? AND forum_id = ?
		) AS is_exist`, userId, forum.Id).Scan(&checkForumIsLike).Error

		if errCheckForumIsLike != nil {
			helper.Logger("error", "In Server: "+errCheckForumIsLike.Error())
			return nil, errors.New(errCheckForumIsLike.Error())
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

		rows, errForumComment := db.Debug().Raw(`SELECT fc.uid AS id, fc.created_at, p.avatar, fc.comment, p.user_id, p.fullname 
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

			// Check Forum Comment is Like

			checkForumCommentIsLike := []entities.CheckIsLike{}

			errCheckForumCommentIsLike := db.Debug().Raw(`SELECT EXISTS (
			SELECT 1
			FROM forum_comment_likes
			WHERE user_id = ? AND comment_id = ?
			) AS is_exist`, userId, forumComment.Id).Scan(&checkForumCommentIsLike).Error

			if errCheckForumCommentIsLike != nil {
				helper.Logger("error", "In Server: "+errCheckForumCommentIsLike.Error())
				return nil, errors.New(errCheckForumCommentIsLike.Error())
			}

			var dataForumCommentReply = make([]entities.ForumCommentReply, 0)

			rows, errForumCommentReply := db.Debug().Raw(`SELECT fcr.uid AS id, fcr.created_at, fcr.reply, p.avatar, p.user_id, p.fullname 
			FROM forum_comment_replies fcr
			INNER JOIN profiles p ON p.user_id = fcr.user_id
			WHERE fcr.comment_id = '` + forumComment.Id + `'`).Rows()

			for rows.Next() {
				errScanRows := db.ScanRows(rows, &forumCommentReply)

				if errScanRows != nil {
					helper.Logger("error", "In Server: "+errScanRows.Error())
					return nil, errors.New(errScanRows.Error())
				}

				// Check Forum Comment Reply is Like

				checkForumCommentReplyIsLike := []entities.CheckIsLike{}

				errCheckForumReplyCommentIsLike := db.Debug().Raw(`SELECT EXISTS (
						SELECT 1
						FROM forum_comment_reply_likes
						WHERE user_id = ? AND reply_id = ?
					) AS is_exist`, userId, forumCommentReply.Id).Scan(&checkForumCommentReplyIsLike).Error

				if errCheckForumReplyCommentIsLike != nil {
					helper.Logger("error", "In Server: "+errCheckForumReplyCommentIsLike.Error())
					return nil, errors.New(errCheckForumReplyCommentIsLike.Error())
				}

				dataForumCommentReply = append(dataForumCommentReply, entities.ForumCommentReply{
					Id:        forumCommentReply.Id,
					Reply:     forumCommentReply.Reply,
					IsLiked:   checkForumCommentReplyIsLike[0].IsExist,
					CreatedAt: forumCommentReply.CreatedAt,
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
			forumCommentAssign.IsLiked = checkForumCommentIsLike[0].IsExist
			forumCommentAssign.CreatedAt = forumComment.CreatedAt
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
			IsLiked:      checkForumIsLike[0].IsExist,
			ForumType: entities.ForumType{
				Id:   forum.ForumTypeId,
				Name: forum.ForumTypeName,
			},
			User: entities.ForumUser{
				Id:       forum.UserId,
				Avatar:   forum.Avatar,
				Fullname: forum.Fullname,
				Email:    forum.Email,
				Phone:    forum.Phone,
			},
			CreatedAt: forum.CreatedAt,
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

func ForumDetail(f *models.Forum) (map[string]any, error) {
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

	// Check Forum is Like

	checkForumIsLike := []entities.CheckIsLike{}

	errCheckForumIsLike := db.Debug().Raw(`SELECT EXISTS (
		SELECT 1
		FROM forum_likes
		WHERE user_id = ? AND forum_id = ?
	) AS is_exist`, f.UserId, f.Id).Scan(&checkForumIsLike).Error

	if errCheckForumIsLike != nil {
		helper.Logger("error", "In Server: "+errCheckForumIsLike.Error())
		return nil, errors.New(errCheckForumIsLike.Error())
	}

	forum := []entities.Forum{}

	errForum := db.Debug().Raw(`SELECT f.uid AS id, f.title, f.caption,
	p.avatar,
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

	rows, errForumComment := db.Debug().Raw(`SELECT fc.uid AS id, fc.created_at, p.avatar, fc.comment, p.user_id, p.fullname 
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

		// Check Forum Comment is Like

		checkForumCommentIsLike := []entities.CheckIsLike{}

		errCheckForumCommentIsLike := db.Debug().Raw(`SELECT EXISTS (
		SELECT 1
		FROM forum_comment_likes
		WHERE user_id = ? AND comment_id = ?
		) AS is_exist`, f.UserId, forumComment.Id).Scan(&checkForumCommentIsLike).Error

		if errCheckForumCommentIsLike != nil {
			helper.Logger("error", "In Server: "+errCheckForumCommentIsLike.Error())
			return nil, errors.New(errCheckForumCommentIsLike.Error())
		}

		var dataForumCommentReply = make([]entities.ForumCommentReply, 0)

		rows, errForumCommentReply := db.Debug().Raw(`SELECT fcr.uid AS id, fcr.created_at, fcr.reply, p.avatar, p.user_id, p.fullname 
		FROM forum_comment_replies fcr
		INNER JOIN profiles p ON p.user_id = fcr.user_id
		WHERE fcr.comment_id = '` + forumComment.Id + `'`).Rows()

		for rows.Next() {
			errScanRows := db.ScanRows(rows, &forumCommentReply)

			if errScanRows != nil {
				helper.Logger("error", "In Server: "+errScanRows.Error())
				return nil, errors.New(errScanRows.Error())
			}

			// Check Forum Comment Reply is Like

			checkForumCommentReplyIsLike := []entities.CheckIsLike{}

			errCheckForumReplyCommentIsLike := db.Debug().Raw(`SELECT EXISTS (
				SELECT 1
				FROM forum_comment_reply_likes
				WHERE user_id = ? AND reply_id = ?
			) AS is_exist`, f.UserId, forumCommentReply.Id).Scan(&checkForumCommentReplyIsLike).Error

			if errCheckForumReplyCommentIsLike != nil {
				helper.Logger("error", "In Server: "+errCheckForumReplyCommentIsLike.Error())
				return nil, errors.New(errCheckForumReplyCommentIsLike.Error())
			}

			dataForumCommentReply = append(dataForumCommentReply, entities.ForumCommentReply{
				Id:        forumCommentReply.Id,
				Reply:     forumCommentReply.Reply,
				IsLiked:   checkForumCommentReplyIsLike[0].IsExist,
				CreatedAt: forumCommentReply.CreatedAt,
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
		forumCommentAssign.CreatedAt = forumComment.CreatedAt
		forumCommentAssign.Reply = dataForumCommentReply
		forumCommentAssign.IsLiked = checkForumCommentIsLike[0].IsExist
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
		IsLiked:      checkForumIsLike[0].IsExist,
		LikeCount:    len(dataForumLike),
		ForumType: entities.ForumType{
			Id:   forum[0].ForumTypeId,
			Name: forum[0].ForumTypeName,
		},
		User: entities.ForumUser{
			Id:       forum[0].UserId,
			Avatar:   forum[0].Avatar,
			Fullname: forum[0].Fullname,
			Email:    forum[0].Email,
			Phone:    forum[0].Phone,
		},
		CreatedAt: forum[0].CreatedAt,
	})

	return map[string]any{
		"data": &appendForumAssign[0],
	}, nil
}

func ForumCategory() (map[string]any, error) {
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

func ForumStore(f *entities.ForumStore) (map[string]any, error) {
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

	return map[string]any{
		"data": forum.Id,
	}, nil
}

func ForumStoreLike(f *entities.ForumStoreLike) (map[string]any, error) {
	var count int64

	f.Id = uuid.NewV4().String()

	errCheck := db.Debug().Table("forum_likes").
		Where("user_id = ? AND forum_id = ?", f.UserId, f.ForumId).
		Count(&count).Error
	if errCheck != nil {
		helper.Logger("error", "In Server (Check): "+errCheck.Error())
		return nil, errors.New(errCheck.Error())
	}

	if count > 0 {
		errDelete := db.Debug().Exec(`
			DELETE FROM forum_likes WHERE user_id = ? AND forum_id = ?
		`, f.UserId, f.ForumId).Error
		if errDelete != nil {
			helper.Logger("error", "In Server (Delete): "+errDelete.Error())
			return nil, errors.New(errDelete.Error())
		}

		return map[string]any{"message": "unliked"}, nil
	}

	// If it doesn't exist, insert the like
	errInsert := db.Debug().Exec(`
		INSERT INTO forum_likes (uid, user_id, forum_id) VALUES (?, ?, ?)
	`, f.Id, f.UserId, f.ForumId).Error
	if errInsert != nil {
		helper.Logger("error", "In Server (Insert): "+errInsert.Error())
		return nil, errors.New(errInsert.Error())
	}

	return map[string]any{"message": "liked"}, nil
}

func CommentStore(c *entities.CommentStore) (map[string]any, error) {

	errInsertComment := db.Debug().Exec(`INSERT INTO forum_comments (uid, forum_id, user_id, comment) 
	VALUES (?, ?, ?, ?)`, c.Id, c.ForumId, c.UserId, c.Comment).Error

	if errInsertComment != nil {
		helper.Logger("error", "In Server: "+errInsertComment.Error())
		return nil, errors.New(errInsertComment.Error())
	}
	return map[string]any{
		"data": c.Id,
	}, nil
}

func ReplyStore(r *entities.ReplyStore) (map[string]any, error) {

	errInsertReply := db.Debug().Exec(`INSERT INTO forum_comment_replies (uid, user_id, comment_id, reply) 
	VALUES (?, ?, ?, ?)`, r.Id, r.UserId, r.CommentId, r.Reply).Error

	if errInsertReply != nil {
		helper.Logger("error", "In Server: "+errInsertReply.Error())
		return nil, errors.New(errInsertReply.Error())
	}
	return map[string]any{
		"data": r.Id,
	}, nil
}

func ForumDelete(f *models.Forum) (map[string]any, error) {
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
