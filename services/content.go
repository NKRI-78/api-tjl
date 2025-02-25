package services

import (
	"errors"
	"math"
	"os"
	"strconv"

	// "fmt"
	entities "superapps/entities"
	helper "superapps/helpers"
	models "superapps/models"

	uuid "github.com/satori/go.uuid"
)

func GetContent(search, page, limit, appName string) (map[string]interface{}, error) {

	url := os.Getenv("API_URL")

	var appAssign entities.ContentApplicationResponse

	var contentMedia entities.ContentMedia
	var contentMediaAssign entities.ContentMediaResponse

	var content entities.Content

	var contentLike entities.ContentLike
	var contentLikeAssign entities.ContentLikeResponse
	var contentLikeUserAssign entities.ContentLikeUserResponse

	var contentUnlike entities.ContentUnlike
	var contentUnlikeAssign entities.ContentUnlikeResponse
	var contentUnlikeUserAssign entities.ContentUnlikeUserResponse

	var contentComment entities.ContentComment
	var contentCommentAssign entities.ContentCommentResponse
	var contentCommentUserAssign entities.ContentCommentUserResponse

	var user entities.ContentUser
	var userAssign entities.ContentUserResponse

	var contentAssign entities.ContentResponse
	var appendContentAssign = make([]entities.ContentResponse, 0)

	var allCountContent []models.AllCountContent

	pageinteger, _ := strconv.Atoi(page)
	limitinteger, _ := strconv.Atoi(limit)

	var offset = strconv.Itoa((pageinteger - 1) * limitinteger)

	errAllContentQuery := db.Debug().Raw(`SELECT uid FROM contents`).Scan(&allCountContent).Error

	if errAllContentQuery != nil {
		helper.Logger("error", "In Server: "+errAllContentQuery.Error())
	}

	var resultTotal = len(allCountContent)

	var perPage = math.Ceil(float64(resultTotal) / float64(limitinteger))

	var prevPage int
	var nextPage int

	if pageinteger == 1 {
		prevPage = 1
	} else {
		prevPage = pageinteger - 1
	}

	nextPage = pageinteger + 1

	rows, errContentQuery := db.Debug().Raw(`SELECT n.uid, n.title, ct.name AS type, n.description, n.user_id, n.created_at, 
	app.name AS app_name, app.uid AS app_id 
	FROM contents n 
	INNER JOIN applications app ON app.uid = n.app_id 
	INNER JOIN content_types ct ON ct.id = n.type_id
	WHERE n.title LIKE '%` + search + `%' AND app.name LIKE '%` + appName + `%'
	LIMIT ` + offset + `, ` + limit + ``).Rows()

	if errContentQuery != nil {
		helper.Logger("error", "In Server: "+errContentQuery.Error())
		return nil, errors.New(errContentQuery.Error())
	}

	for rows.Next() {
		errScanRows := db.ScanRows(rows, &content)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		rows, errUserQuery := db.Debug().Raw(`SELECT email, phone, fullname FROM users u 
		INNER JOIN user_profiles p ON p.user_id = u.uid 
		WHERE u.uid = '` + content.UserId + `'`).Rows()

		if errUserQuery != nil {
			helper.Logger("error", "In Server: "+errUserQuery.Error())
			return nil, errors.New(errUserQuery.Error())
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

		var dataContentMedia = make([]entities.ContentMediaResponse, 0)

		rows, errContentMediaQuery := db.Debug().Raw(`SELECT content_id, path, size FROM content_medias WHERE content_id = '` + content.Uid + `'`).Rows()

		if errContentMediaQuery != nil {
			helper.Logger("error", "In Server: "+errContentMediaQuery.Error())
			return nil, errors.New(errContentMediaQuery.Error())
		}

		for rows.Next() {
			errScanRows = db.ScanRows(rows, &contentMedia)

			if errScanRows != nil {
				helper.Logger("error", "In Server: "+errScanRows.Error())
				return nil, errors.New(errScanRows.Error())
			}

			if contentMedia.Path != "" {
				contentMediaAssign.ContentId = contentMedia.ContentId
				contentMediaAssign.Path = contentMedia.Path
				contentMediaAssign.Size = contentMedia.Size

				dataContentMedia = append(dataContentMedia, contentMediaAssign)
			}
		}

		var dataContentLike = make([]entities.ContentLikeResponse, 0)

		rows, errContentLikeQuery := db.Debug().Raw(`SELECT cl.uid, p.user_id, p.fullname FROM content_likes cl 
		INNER JOIN user_profiles p ON p.user_id = cl.user_id
		WHERE cl.content_id = '` + content.Uid + `'`).Rows()

		if errContentLikeQuery != nil {
			helper.Logger("error", "In Server: "+errContentLikeQuery.Error())
			return nil, errors.New(errContentLikeQuery.Error())
		}

		for rows.Next() {
			errScanRows = db.ScanRows(rows, &contentLike)

			if errScanRows != nil {
				helper.Logger("error", "In Server: "+errScanRows.Error())
				return nil, errors.New(errScanRows.Error())
			}

			contentLikeUserAssign.UserId = contentLike.UserId
			contentLikeUserAssign.Fullname = contentLike.Fullname

			contentLikeAssign.Uid = contentLike.Uid
			contentLikeAssign.User = contentLikeUserAssign

			dataContentLike = append(dataContentLike, contentLikeAssign)
		}

		var dataContentUnlike = make([]entities.ContentUnlikeResponse, 0)

		rows, errContentUnlikeQuery := db.Debug().Raw(`SELECT cu.uid, p.user_id, p.fullname FROM content_unlikes cu
		INNER JOIN user_profiles p ON p.user_id = cu.user_id
		WHERE cu.content_id = '` + content.Uid + `'`).Rows()

		if errContentUnlikeQuery != nil {
			helper.Logger("error", "In Server: "+errContentUnlikeQuery.Error())
			return nil, errors.New(errContentUnlikeQuery.Error())
		}

		for rows.Next() {
			errScanRows = db.ScanRows(rows, &contentUnlike)

			if errScanRows != nil {
				helper.Logger("error", "In Server: "+errScanRows.Error())
				return nil, errors.New(errScanRows.Error())
			}

			contentUnlikeUserAssign.UserId = contentUnlike.UserId
			contentUnlikeUserAssign.Fullname = contentUnlike.Fullname

			contentUnlikeAssign.Uid = contentUnlike.Uid
			contentUnlikeAssign.User = contentUnlikeUserAssign

			dataContentUnlike = append(dataContentUnlike, contentUnlikeAssign)
		}

		var dataContentComment = make([]entities.ContentCommentResponse, 0)

		rows, errContentCommentQuery := db.Debug().Raw(`SELECT cc.uid, cc.comment, cc.user_id, p.fullname FROM content_comments cc
		INNER JOIN user_profiles p ON p.user_id = cc.user_id
		WHERE content_id = '` + content.Uid + `'`).Rows()

		if errContentCommentQuery != nil {
			helper.Logger("error", "In Server: "+errContentCommentQuery.Error())
			return nil, errors.New(errContentCommentQuery.Error())
		}

		for rows.Next() {
			errScanRows = db.ScanRows(rows, &contentComment)

			if errScanRows != nil {
				helper.Logger("error", "In Server: "+errScanRows.Error())
				return nil, errors.New(errScanRows.Error())
			}

			if contentComment.Uid != "" {

				contentCommentUserAssign.UserId = contentComment.UserId
				contentCommentUserAssign.Fullname = contentComment.Fullname

				contentCommentAssign.Uid = contentComment.Uid
				contentCommentAssign.Comment = contentComment.Comment

				contentCommentAssign.User = contentCommentUserAssign

				dataContentComment = append(dataContentComment, contentCommentAssign)
			}
		}

		var createdAt = content.CreatedAt.Format("2006-01-02 15:04")

		appAssign.ApplicationId = content.AppId
		appAssign.ApplicationName = content.AppName

		contentAssign.Uid = content.Uid
		contentAssign.Title = content.Title
		contentAssign.Description = content.Description

		contentAssign.Media = dataContentMedia
		contentAssign.Like = dataContentLike
		contentAssign.Unlike = dataContentUnlike
		contentAssign.Comment = dataContentComment

		contentAssign.App = appAssign
		contentAssign.User = userAssign
		contentAssign.CreatedAt = createdAt
		contentAssign.Type = content.Type

		appendContentAssign = append(appendContentAssign, contentAssign)
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
		"data":         &appendContentAssign,
	}, nil
}

func CreateMediaContent(n *models.ContentMedia) (map[string]interface{}, error) {
	// Uid  := n.ContentId
	// Path := n.Path

	return map[string]interface{}{}, nil
}

func CreateContent(n *models.Content) (map[string]interface{}, error) {

	applications := []entities.Application{}
	types := []entities.ContentTypes{}

	errCheckApp := db.Debug().Raw(`SELECT uid, name FROM applications WHERE name = '` + n.AppName + `'`).Scan(&applications).Error

	if errCheckApp != nil {
		helper.Logger("error", "In Server: "+errCheckApp.Error())
		return nil, errors.New(errCheckApp.Error())
	}

	isAppExist := len(applications)

	if isAppExist == 0 {
		return nil, errors.New("app not found")
	}

	errTypes := db.Debug().Raw(`SELECT name FROM content_types WHERE id = '` + strconv.Itoa(n.TypeId) + `'`).Scan(&types).Error

	if errTypes != nil {
		helper.Logger("error", "In Server: "+errTypes.Error())
		return nil, errors.New(errTypes.Error())
	}

	isTypesExist := len(types)

	if isTypesExist == 0 {
		return nil, errors.New("types not found")
	}

	ApplicationId := applications[0].Uid

	errIns := db.Debug().Exec(`INSERT INTO contents (uid, title, description, type_id, app_id, user_id) 
	VALUES ('` + n.Uid + `', '` + n.Title + `', '` + n.Description + `', '` + strconv.Itoa(n.TypeId) + `' ,'` + ApplicationId + `', '` + n.UserId + `')`).Error

	if errIns != nil {
		helper.Logger("error", "In Server: "+errIns.Error())
		return nil, errors.New(errIns.Error())
	}

	return map[string]interface{}{}, nil
}

func CreateContentLike(l *models.ReqContentLike) (map[string]interface{}, error) {

	contentLike := []entities.ContentLike{}

	uid := uuid.NewV4().String()

	checkLike := db.Debug().Raw(`SELECT uid FROM content_likes WHERE user_id = '` + l.UserId + `'`).Scan(&contentLike).Error

	if checkLike != nil {
		helper.Logger("error", "In Server: "+checkLike.Error())
		return nil, errors.New(checkLike.Error())
	}

	isLikeExist := len(contentLike)

	if isLikeExist != 0 {
		delLike := db.Debug().Exec(`DELETE FROM content_likes WHERE user_id = '` + l.UserId + `'`).Error

		if delLike != nil {
			helper.Logger("error", "In Server: "+delLike.Error())
			return nil, errors.New(delLike.Error())
		}
	} else {
		createLike := db.Debug().Exec(`INSERT INTO content_likes (uid, content_id, user_id) 
		VALUES ('` + uid + `', '` + l.ContentId + `', '` + l.UserId + `')`).Error

		if createLike != nil {
			helper.Logger("error", "In Server: "+createLike.Error())
			return nil, errors.New(createLike.Error())
		}
	}

	return map[string]interface{}{}, nil
}

func CreateContentUnlike(l *models.ReqContentUnlike) (map[string]interface{}, error) {

	contentUnlike := []entities.ContentUnlike{}

	uid := uuid.NewV4().String()

	checkUnlike := db.Debug().Raw(`SELECT uid FROM content_unlikes WHERE user_id = '` + l.UserId + `'`).Scan(&contentUnlike).Error

	if checkUnlike != nil {
		helper.Logger("error", "In Server: "+checkUnlike.Error())
		return nil, errors.New(checkUnlike.Error())
	}

	isUnlikeExist := len(contentUnlike)

	if isUnlikeExist != 0 {
		delUnlike := db.Debug().Exec(`DELETE FROM content_unlikes WHERE user_id = '` + l.UserId + `'`).Error

		if delUnlike != nil {
			helper.Logger("error", "In Server: "+delUnlike.Error())
			return nil, errors.New(delUnlike.Error())
		}
	} else {
		createUnlike := db.Debug().Exec(`INSERT INTO content_unlikes (uid, content_id, user_id) 
		VALUES ('` + uid + `', '` + l.ContentId + `', '` + l.UserId + `')`).Error

		if createUnlike != nil {
			helper.Logger("error", "In Server: "+createUnlike.Error())
			return nil, errors.New(createUnlike.Error())
		}
	}

	return map[string]interface{}{}, nil
}

func CreateContentComment(c *models.ReqContentComment) (map[string]interface{}, error) {

	uid := uuid.NewV4().String()

	errIns := db.Debug().Exec(`INSERT INTO content_comments (uid, content_id, comment, user_id) 
	VALUES ('` + uid + `', '` + c.ContentId + `', '` + c.Comment + `', '` + c.UserId + `')`).Error

	if errIns != nil {
		helper.Logger("error", "In Server: "+errIns.Error())
		return nil, errors.New(errIns.Error())
	}

	return map[string]interface{}{}, nil
}

func DelContent(d *models.DelContent) (map[string]interface{}, error) {

	errDel := db.Debug().Exec(`DELETE FROM contents WHERE uid = '` + d.Uid + `'`).Error

	if errDel != nil {
		helper.Logger("error", "In Server: "+errDel.Error())
		return nil, errors.New(errDel.Error())
	}

	return map[string]interface{}{}, nil
}

func DelContentComment(d *models.DelContentComment) (map[string]interface{}, error) {

	errDel := db.Debug().Exec(`DELETE FROM content_comments WHERE uid = '` + d.Uid + `'`).Error

	if errDel != nil {
		helper.Logger("error", "In Server: "+errDel.Error())
		return nil, errors.New(errDel.Error())
	}

	return map[string]interface{}{}, nil
}
