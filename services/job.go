package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	entities "superapps/entities"
	helper "superapps/helpers"
	models "superapps/models"

	uuid "github.com/satori/go.uuid"
)

func CandidatePassesList(branchId string) (map[string]any, error) {
	var data []entities.ResultCandidateInfoApplyJob

	query := `
		SELECT 
			paa.user_id AS apply_user_id,
			paa.fullname AS apply_user_name, 
			pac.user_id AS confirm_user_id,
			pac.fullname AS confirm_user_name,
			u.phone AS apply_user_phone,
			u.email AS apply_user_email,
			ajo.content AS invitation_offline,
			ibs.field1 AS invitation_departure,
			js.name AS status,
			aj.uid AS apply_job_id,
			j.title AS job_title,
			jc.name AS job_category,
			p.avatar AS job_avatar,
			p.fullname AS job_author,
			c.uid AS company_id,
			c.logo AS company_logo,
			c.name AS company_name,
			pl.name AS country_name,
			aj.created_at
		FROM apply_jobs aj
			INNER JOIN jobs j ON j.uid = aj.job_id
			INNER JOIN companies c ON c.uid = j.company_id 
			INNER JOIN places pl ON pl.id = c.place_id
			INNER JOIN job_categories jc ON jc.uid = j.cat_id
			INNER JOIN profiles p ON p.user_id = j.user_id
			INNER JOIN job_statuses js ON js.id = aj.status
			INNER JOIN profiles paa ON paa.user_id = aj.user_id
			INNER JOIN users u ON u.uid = aj.user_id
			INNER JOIN user_branches ub ON ub.user_id = aj.user_id
			INNER JOIN branchs b ON b.id = ub.branch_id
			LEFT JOIN apply_job_offlines ajo ON ajo.apply_job_id = aj.uid
			LEFT JOIN (
				SELECT field2 AS apply_job_id, MAX(field1) AS field1
				FROM inboxes
				GROUP BY field2
			) ibs ON ibs.apply_job_id = aj.uid
			LEFT JOIN profiles pac ON pac.user_id = aj.user_confirm_id 
		WHERE aj.status = ?
	`

	var args []any
	args = append(args, "11")

	if branchId != "" {
		query += " AND b.id = ?"
		args = append(args, branchId)
	}

	query += " ORDER BY aj.created_at DESC"

	rows, err := db.Debug().Raw(query, args...).Rows()
	if err != nil {
		helper.Logger("error", "In Server (main query): "+err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dataQuery entities.InfoApplyJobQuery

		if err := db.ScanRows(rows, &dataQuery); err != nil {
			helper.Logger("error", "In Server (scan main row): "+err.Error())
			return nil, err
		}

		// Count form filled
		var count int
		row := db.Raw(
			`SELECT COUNT(*) FROM inboxes WHERE field2 = ? AND user_id = ?`,
			dataQuery.ApplyJobId, dataQuery.ApplyUserId,
		).Row()
		if err := row.Scan(&count); err != nil {
			helper.Logger("error", "In Server (count query): "+err.Error())
			return nil, err
		}
		formFilled := count > 0

		// Candidate documents
		dataCandidateDocument := make([]entities.CandidateDocument, 0)
		queryCandidateDocument := `
			SELECT d.name AS document, ajd.path
			FROM documents d
				INNER JOIN apply_job_documents ajd ON ajd.doc_id = d.id
			WHERE ajd.apply_job_id = ?
		`

		rowsCandidateDocument, err := db.Debug().Raw(queryCandidateDocument, dataQuery.ApplyJobId).Rows()
		if err != nil {
			helper.Logger("error", "In Server (candidate document query): "+err.Error())
			return nil, err
		}
		defer rowsCandidateDocument.Close()

		for rowsCandidateDocument.Next() {
			var candidateDoc entities.CandidateDocumentQuery
			if err := db.ScanRows(rowsCandidateDocument, &candidateDoc); err != nil {
				helper.Logger("error", "In Server (scan candidate doc): "+err.Error())
				return nil, err
			}
			dataCandidateDocument = append(dataCandidateDocument, entities.CandidateDocument{
				Document: candidateDoc.Document,
				Path:     candidateDoc.Path,
			})
		}
		docFilled := len(dataCandidateDocument) > 0

		data = append(data, entities.ResultCandidateInfoApplyJob{
			Id:                  dataQuery.ApplyJobId,
			Status:              dataQuery.Status,
			CreatedAt:           dataQuery.CreatedAt,
			FormFilled:          formFilled,
			DocFilled:           docFilled,
			InvitationOffline:   helper.DefaultIfEmpty(dataQuery.InvitationOffline, "-"),
			InvitationDeparture: helper.DefaultIfEmpty(dataQuery.InvitationDeparture, "-"),
			Job: entities.JobApply{
				JobTitle:    dataQuery.JobTitle,
				JobCategory: dataQuery.JobCategory,
				JobAvatar:   dataQuery.JobAvatar,
				JobAuthor:   dataQuery.JobAuthor,
			},
			Company: entities.JobCompanyCandidate{
				Id:      dataQuery.CompanyId,
				Logo:    dataQuery.CompanyLogo,
				Name:    dataQuery.CompanyName,
				Country: dataQuery.CountryName,
			},
			UserApply: entities.UserApply{
				Id:    dataQuery.ApplyUserId,
				Name:  dataQuery.ApplyUserName,
				Email: dataQuery.ApplyUserEmail,
				Phone: dataQuery.ApplyUserPhone,
			},
			UserConfirm: entities.UserConfirm{
				Id:   helper.DefaultIfEmpty(dataQuery.ConfirmUserId, "-"),
				Name: helper.DefaultIfEmpty(dataQuery.ConfirmUserName, "-"),
			},
		})
	}

	if data == nil {
		data = []entities.ResultCandidateInfoApplyJob{}
	}

	return map[string]any{
		"data": data,
	}, nil
}

func CandidatePassesFormList() (map[string]any, error) {

	var dataQuery entities.CandidatePassesFormListQuery
	var data []entities.CandidatePassesFormListResult

	query := `SELECT d.id, d.date_departure, d.time_departure, d.airplane, d.content, d.location, d.destination, d.created_at, d.updated_at,
	 	p.fullname AS user_fullname, p.avatar AS user_avatar, p.user_id
		FROM departures d 
		INNER JOIN candidate_passes cp 
		ON cp.departure_id = d.id
		INNER JOIN profiles p ON p.user_id = cp.user_candidate_id
	`
	rows, err := db.Debug().Raw(query).Rows()

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		errDepartureRows := db.ScanRows(rows, &dataQuery)

		if errDepartureRows != nil {
			helper.Logger("error", "In Server: "+errDepartureRows.Error())
			return nil, errors.New(errDepartureRows.Error())
		}

		data = append(data, entities.CandidatePassesFormListResult{
			Id:      dataQuery.Id,
			Content: dataQuery.Content,
			User: entities.CandidatePassesFormUser{
				Id:       dataQuery.UserId,
				Avatar:   dataQuery.UserAvatar,
				Fullname: dataQuery.UserFullname,
			},
			CreatedAt: dataQuery.CreatedAt,
			UpdatedAt: dataQuery.UpdatedAt,
		})
	}

	if data == nil {
		data = []entities.CandidatePassesFormListResult{}
	}

	return map[string]any{
		"data": data,
	}, nil
}

func CandidateInfoDeparture(userId string) (map[string]any, error) {
	var dataQuery entities.CandidatePassesFormListQuery
	var data []entities.CandidatePassesFormListResult

	query := `SELECT d.id, d.date_departure, d.time_departure, d.airplane, d.location, d.destination, d.created_at, d.updated_at,
	 	p.fullname AS user_fullname, p.avatar AS user_avatar, p.user_id
		FROM departures d 
		INNER JOIN candidate_passes cp 
		ON cp.departure_id = d.id
		INNER JOIN profiles p ON p.user_id = cp.user_candidate_id
		WHERE cp.user_candidate_id = ?
	`
	rows, err := db.Debug().Raw(query, userId).Rows()

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		errDepartureRows := db.ScanRows(rows, &dataQuery)

		if errDepartureRows != nil {
			helper.Logger("error", "In Server: "+errDepartureRows.Error())
			return nil, errors.New(errDepartureRows.Error())
		}

		data = append(data, entities.CandidatePassesFormListResult{
			Id:      dataQuery.Id,
			Content: dataQuery.Content,
			User: entities.CandidatePassesFormUser{
				Id:       dataQuery.UserId,
				Avatar:   dataQuery.UserAvatar,
				Fullname: dataQuery.UserFullname,
			},
			CreatedAt: dataQuery.CreatedAt,
			UpdatedAt: dataQuery.UpdatedAt,
		})
	}

	if data == nil {
		data = []entities.CandidatePassesFormListResult{}
	}

	return map[string]any{
		"data": data,
	}, nil
}

func ListInfoApplyJob(iaj *models.InfoApplyJob) (map[string]any, error) {

	var dataQuery entities.InfoApplyJobQuery
	var data []entities.ResultInfoApplyJob

	query := `SELECT paa.user_id AS apply_user_id, paa.fullname AS apply_user_name, 
		pac.user_id AS confirm_user_id, pac.fullname AS confirm_user_name,
		u.phone AS apply_user_phone,
		u.email AS apply_user_email,
		js.name AS status, aj.uid AS apply_job_id,
		j.title AS job_title,
		jc.name AS job_category,
		p.avatar AS job_avatar,
		p.fullname AS job_author,
		c.uid AS company_id,
		c.logo AS company_logo,
		c.name AS company_name,
		pl.name AS country_name,
		aj.created_at
		FROM apply_jobs aj 
		INNER JOIN jobs j ON j.uid = aj.job_id
		INNER JOIN companies c ON c.uid = j.company_id 
		INNER JOIN job_categories jc ON jc.uid = j.cat_id
		INNER JOIN profiles p ON p.user_id = j.user_id
		INNER JOIN places pl ON pl.id = c.place_id
		INNER JOIN job_statuses js ON js.id = aj.status
		INNER JOIN profiles paa ON paa.user_id = aj.user_id
		INNER JOIN users u ON u.uid = paa.user_id
		LEFT JOIN profiles pac ON pac.user_id = aj.user_confirm_id 
		WHERE aj.user_id = ?
		ORDER BY aj.created_at DESC
	`
	rows, err := db.Debug().Raw(query, iaj.UserId).Rows()

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		errJobRows := db.ScanRows(rows, &dataQuery)

		if errJobRows != nil {
			helper.Logger("error", "In Server: "+errJobRows.Error())
			return nil, errors.New(errJobRows.Error())
		}

		var candidateApplyJobId string
		candidatePassQuery := `SELECT apply_job_id FROM candidate_passes WHERE apply_job_id = ? LIMIT 1`
		errPass := db.Debug().Raw(candidatePassQuery, dataQuery.ApplyJobId).Row().Scan(&candidateApplyJobId)
		if errPass != nil {
			helper.Logger("error", "In Server: "+errPass.Error())
		}

		readyDeparture := candidateApplyJobId != ""

		data = append(data, entities.ResultInfoApplyJob{
			Id:             dataQuery.ApplyJobId,
			Status:         dataQuery.Status,
			CreatedAt:      dataQuery.CreatedAt,
			ReadyDeparture: readyDeparture,
			Job: entities.JobApply{
				JobTitle:    dataQuery.JobTitle,
				JobCategory: dataQuery.JobCategory,
				JobAvatar:   dataQuery.JobAvatar,
				JobAuthor:   dataQuery.JobAuthor,
			},
			Company: entities.JobCompany{
				Id:      dataQuery.CompanyId,
				Logo:    dataQuery.CompanyLogo,
				Name:    dataQuery.CompanyName,
				Country: dataQuery.CountryName,
			},
			UserApply: entities.UserApply{
				Id:    dataQuery.ApplyUserId,
				Name:  dataQuery.ApplyUserName,
				Email: dataQuery.ApplyUserEmail,
				Phone: dataQuery.ApplyUserPhone,
			},
			UserConfirm: entities.UserConfirm{
				Id:   helper.DefaultIfEmpty(dataQuery.ConfirmUserId, "-"),
				Name: helper.DefaultIfEmpty(dataQuery.ConfirmUserName, "-"),
			},
		})
	}

	if data == nil {
		data = []entities.ResultInfoApplyJob{}
	}

	return map[string]any{
		"data": data,
	}, nil
}

func ListApplyJobHistory(iaj *models.InfoApplyJob) (map[string]any, error) {
	var dataQuery entities.InfoApplyJobQuery
	var data []entities.ResultInfoApplyJob

	query := `
	SELECT 
		paa.user_id AS apply_user_id, paa.fullname AS apply_user_name, 
		pac.user_id AS confirm_user_id, pac.fullname AS confirm_user_name,
		js.name AS status, aj.uid AS apply_job_id,
		j.title AS job_title,
		jc.name AS job_category,
		p.avatar AS job_avatar,
		p.fullname AS job_author,
		c.uid AS company_id,
		c.logo AS company_logo,
		c.name AS company_name,
		aj.created_at
	FROM apply_jobs aj 
	INNER JOIN jobs j ON j.uid = aj.job_id
	INNER JOIN companies c ON c.uid = j.company_id 
	INNER JOIN job_categories jc ON jc.uid = j.cat_id
	INNER JOIN profiles p ON p.user_id = j.user_id
	INNER JOIN job_statuses js ON js.id = aj.status
	INNER JOIN profiles paa ON paa.user_id = aj.user_id
	LEFT JOIN profiles pac ON pac.user_id = aj.user_confirm_id 
	WHERE aj.user_id = ? 
	AND DATE(aj.created_at) > NOW()
	ORDER BY aj.created_at DESC
	`

	rows, err := db.Debug().Raw(query, iaj.UserId).Rows()
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		errScan := db.ScanRows(rows, &dataQuery)
		if errScan != nil {
			helper.Logger("error", "In Server: "+errScan.Error())
			return nil, errScan
		}

		var candidateApplyJobId string
		candidatePassQuery := `SELECT apply_job_id FROM candidate_passes WHERE apply_job_id = ? LIMIT 1`
		errPass := db.Debug().Raw(candidatePassQuery, dataQuery.ApplyJobId).Row().Scan(&candidateApplyJobId)
		if errPass != nil {
			helper.Logger("error", "In Server: "+errPass.Error())
		}

		readyDeparture := candidateApplyJobId != ""

		data = append(data, entities.ResultInfoApplyJob{
			Id:             dataQuery.ApplyJobId,
			Status:         dataQuery.Status,
			CreatedAt:      dataQuery.CreatedAt,
			ReadyDeparture: readyDeparture,
			Job: entities.JobApply{
				JobTitle:    dataQuery.JobTitle,
				JobCategory: dataQuery.JobCategory,
				JobAvatar:   dataQuery.JobAvatar,
				JobAuthor:   dataQuery.JobAuthor,
			},
			Company: entities.JobCompany{
				Id:   dataQuery.CompanyId,
				Logo: dataQuery.CompanyLogo,
				Name: dataQuery.CompanyName,
			},
			UserApply: entities.UserApply{
				Id:   dataQuery.ApplyUserId,
				Name: dataQuery.ApplyUserName,
			},
			UserConfirm: entities.UserConfirm{
				Id:   helper.DefaultIfEmpty(dataQuery.ConfirmUserId, "-"),
				Name: helper.DefaultIfEmpty(dataQuery.ConfirmUserName, "-"),
			},
		})
	}

	if data == nil {
		data = []entities.ResultInfoApplyJob{}
	}

	return map[string]any{
		"data": data,
	}, nil
}

func InfoApplyJob(iaj *models.InfoApplyJob) (map[string]any, error) {

	var dataQuery entities.InfoApplyJobQuery
	var dataDocQuery entities.DocApplyQuery
	var data []entities.ResultInfoJobDetail
	var dataDoc []entities.DocApply

	query := `
	SELECT 
		paa.user_id AS apply_user_id,
		paa.fullname AS apply_user_name,
		pac.user_id AS confirm_user_id,
		pac.fullname AS confirm_user_name,
		js.name AS status,
		aj.created_at,
		aj.uid AS apply_job_id,
		aj.link,
		aj.schedule,
		j.title AS job_title,
		jc.name AS job_category,
		p.avatar AS job_avatar,
		p.fullname AS job_author,
		c.uid AS company_id,
		c.logo AS company_logo,
		c.name AS company_name
	FROM apply_job_histories aj
	INNER JOIN jobs j ON j.uid = aj.job_id
	INNER JOIN companies c ON c.uid = j.company_id
	INNER JOIN job_categories jc ON jc.uid = j.cat_id
	INNER JOIN profiles p ON p.user_id = j.user_id
	INNER JOIN job_statuses js ON js.id = aj.status
	INNER JOIN profiles paa ON paa.user_id = aj.user_id
	LEFT JOIN profiles pac ON pac.user_id = aj.user_confirm_id
	WHERE aj.uid = ?
	ORDER BY aj.created_at ASC`

	rows, err := db.Debug().Raw(query, iaj.Id).Rows()
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		errJobRows := db.ScanRows(rows, &dataQuery)

		if errJobRows != nil {
			helper.Logger("error", "In Server: "+errJobRows.Error())
			return nil, errors.New(errJobRows.Error())
		}

		queryDoc := `SELECT 
			d.id, 
			d.name, 
		  COALESCE(ajd.path, '-') AS path 
		FROM 
			documents d
		LEFT JOIN 
			apply_job_documents ajd 
			ON ajd.doc_id = d.id 
			AND ajd.apply_job_id = ?`

		rowsDoc, errDoc := db.Debug().Raw(queryDoc, dataQuery.ApplyJobId).Rows()

		if errDoc != nil {
			helper.Logger("error", "In Server: "+errDoc.Error())
		}
		defer rowsDoc.Close()

		dataDoc = []entities.DocApply{}

		for rowsDoc.Next() {
			errDocRows := db.ScanRows(rowsDoc, &dataDocQuery)

			if errDocRows != nil {
				helper.Logger("error", "In Server: "+errDocRows.Error())
				return nil, errors.New(errDocRows.Error())
			}

			dataDoc = append(dataDoc, entities.DocApply{
				DocId:   dataDocQuery.Id,
				DocName: dataDocQuery.Name,
				DocPath: dataDocQuery.Path,
			})
		}

		queryJobOffline := `SELECT apply_job_id FROM apply_job_offlines WHERE apply_job_id = ?`
		var offlineFlag bool
		rowsJobOffline, errJobOffline := db.Debug().Raw(queryJobOffline, dataQuery.ApplyJobId).Rows()

		if errJobOffline != nil {
			helper.Logger("error", "In Server: "+errJobOffline.Error())
		} else {
			if rowsJobOffline.Next() {
				offlineFlag = true
			} else {
				offlineFlag = false
			}
		}
		defer rowsJobOffline.Close()

		data = append(data, entities.ResultInfoJobDetail{
			Id:        dataQuery.ApplyJobId,
			Status:    dataQuery.Status,
			Doc:       dataDoc,
			CreatedAt: dataQuery.CreatedAt,
			Offline:   offlineFlag,
			Link:      helper.DefaultIfEmpty(dataQuery.Link, "-"),
			Schedule:  helper.DefaultIfEmpty(dataQuery.Schedule, "-"),
			Job: entities.JobApply{
				JobTitle:    dataQuery.JobTitle,
				JobCategory: dataQuery.JobCategory,
				JobAvatar:   dataQuery.JobAvatar,
				JobAuthor:   dataQuery.JobAuthor,
			},
			Company: entities.JobCompany{
				Id:   dataQuery.CompanyId,
				Logo: dataQuery.CompanyLogo,
				Name: dataQuery.CompanyName,
			},
			UserApply: entities.UserApply{
				Id:   dataQuery.ApplyUserId,
				Name: dataQuery.ApplyUserName,
			},
			UserConfirm: entities.UserConfirm{
				Id:   helper.DefaultIfEmpty(dataQuery.ConfirmUserId, "-"),
				Name: helper.DefaultIfEmpty(dataQuery.ConfirmUserName, "-"),
			},
		})
	}

	return map[string]any{
		"data": data,
	}, nil
}

func ApplyJob(aj *models.ApplyJob) (map[string]any, error) {
	var dataUserFcm entities.InitFcm
	var allJob []models.CheckApplyJobQuery

	queryCheck := `SELECT uid FROM apply_jobs
	WHERE user_id = ?
	AND job_id = ?`

	errAllJob := db.Debug().Raw(queryCheck, aj.UserId, aj.JobId).Scan(&allJob).Error

	if errAllJob != nil {
		helper.Logger("error", "In Server: "+errAllJob.Error())
		return nil, errors.New(errAllJob.Error())
	}

	var isUserAppliedJob = len(allJob)

	if isUserAppliedJob == 1 {
		helper.Logger("error", "In Server: user already applied job")
		return nil, errors.New("USER_ALREADY_APPLIED_JOB")
	}

	// Insert Apply Job
	queryInsert := `INSERT INTO apply_jobs (uid, job_id, user_id) VALUES (?, ?, ?)`

	errInsert := db.Debug().Exec(queryInsert, aj.Id, aj.JobId, aj.UserId).Error

	if errInsert != nil {
		helper.Logger("error", "In Server: "+errInsert.Error())
		return nil, errors.New(errInsert.Error())
	}

	// Insert Apply Job History
	queryHistory := `INSERT INTO apply_job_histories (uid, job_id, user_id) VALUES (?, ?, ?)`

	errHistory := db.Debug().Exec(queryHistory, aj.Id, aj.JobId, aj.UserId).Error

	if errHistory != nil {
		helper.Logger("error", "In Server: "+errHistory.Error())
		return nil, errors.New(errHistory.Error())
	}

	queryUserFcm := `SELECT f.token, p.fullname FROM fcms f
	INNER JOIN profiles p ON p.user_id = f.user_id
	WHERE f.user_id = ?`

	rowUserFcm := db.Debug().Raw(queryUserFcm, aj.UserId).Row()

	errUserFcmRow := rowUserFcm.Scan(&dataUserFcm.Token, &dataUserFcm.Fullname)

	if errUserFcmRow != nil {
		helper.Logger("error", "In Server: "+errUserFcmRow.Error())
	}

	message := fmt.Sprintf("Silahkan menunggu untuk tahap selanjutnya [%s]", dataUserFcm.Fullname)
	helper.SendFcm("Selamat Anda telah berhasil melamar", message, dataUserFcm.Token, "apply-job-detail", aj.Id)

	// Insert Inbox
	queryInsertInbox := `INSERT INTO inboxes (uid, title, caption, user_id, type) VALUES (?, ?, ?, ?, ?)`

	errInsertInbox := db.Debug().Exec(queryInsertInbox, uuid.NewV4().String(), "Selamat Anda telah berhasil melamar", message, dataUserFcm.UserId, "broadcast").Error

	if errInsertInbox != nil {
		helper.Logger("error", "In Server: "+errInsertInbox.Error())
	}

	return map[string]any{
		"data": aj.Id,
	}, nil
}

func ApplyJobBadges(userId string) (map[string]any, error) {
	var dataApplyJobBadges entities.ApplyJobBadges

	query := `
		SELECT COUNT(*) AS total 
		FROM users u 
		INNER JOIN apply_jobs aj ON aj.user_id = u.uid
		WHERE aj.user_id = ? AND aj.is_finish = 0
	`

	row := db.Debug().Raw(query, userId).Row()
	err := row.Scan(&dataApplyJobBadges.Total)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	return map[string]any{
		"data": dataApplyJobBadges.Total,
	}, nil
}

func AssignApplyJob(aaj *entities.AssignApplyJob) (map[string]any, error) {
	Id := uuid.NewV4()

	query := `INSERT INTO apply_jobs (uid, job_id, user_id) 
	VALUES (?, ?, ?)`

	err := db.Debug().Exec(query, Id, aaj.JobId, aaj.UserId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func AssignDocumentApplyJob(adaj *models.AssignDocumentApplyJob) (map[string]any, error) {
	query := `
		INSERT INTO apply_job_documents (apply_job_id, doc_id, path) 
		VALUES (?, ?, ?) 
		ON DUPLICATE KEY UPDATE path = VALUES(path)
	`

	err := db.Debug().Exec(query, adaj.ApplyJobId, adaj.DocId, adaj.Path).Error
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func UpdateApplyJob(uaj *models.ApplyJob) (map[string]any, error) {
	var dataUserFcm entities.InitFcm
	var dataQuery entities.ApplyJobQuery

	queryInfo := `SELECT uid, job_id, user_id, status FROM apply_jobs WHERE uid = ?`
	row := db.Debug().Raw(queryInfo, uaj.ApplyJobId).Row()
	errJobRow := row.Scan(&dataQuery.Uid, &dataQuery.JobId, &dataQuery.UserId, &dataQuery.Status)

	if errJobRow != nil {
		helper.Logger("error", "In Server: "+errJobRow.Error())
		return nil, errors.New("data not found")
	}

	queryUserFcm := `SELECT u.email, f.token, p.fullname FROM fcms f
		INNER JOIN profiles p ON p.user_id = f.user_id
		INNER JOIN users u ON u.uid = f.user_id
		WHERE f.user_id = ?`

	rowUserFcm := db.Debug().Raw(queryUserFcm, dataQuery.UserId).Row()

	errUserFcmRow := rowUserFcm.Scan(&dataUserFcm.Email, &dataUserFcm.Token, &dataUserFcm.Fullname)

	if errUserFcmRow != nil {
		if errors.Is(errUserFcmRow, sql.ErrNoRows) {
			helper.Logger("info", "No FCM data found for user")
		}

		helper.Logger("error", "In Server: "+errUserFcmRow.Error())
	}

	var status string
	var isFinish int64 = 0

	// Daftar nama status
	statusNames := map[int]string{
		1:  "PROCESS",
		2:  "INTERVIEW",
		3:  "FINISH",
		4:  "DECLINE",
		5:  "MCU",
		6:  "LINK_SIAPKERJA",
		7:  "SKCK",
		8:  "VISA",
		9:  "TTD",
		10: "OPP",
		11: "DONE",
	}

	// Validasi status tujuan harus dikenali
	if _, ok := statusNames[uaj.Status]; !ok {
		helper.Logger("error", fmt.Sprintf("unknown target status: %d", uaj.Status))
		return nil, errors.New("unknown target status")
	}

	// Daftar status yang diizinkan sebagai transisi berikutnya
	var validNextStatuses = map[int][]int{
		1:  {2},
		2:  {3, 4}, // interview ke finish (diterima) atau decline (gagal)
		3:  {5},    // lanjut ke MCU
		5:  {6},
		6:  {7},
		7:  {8},
		8:  {9},
		9:  {10},
		10: {11}, // terakhir adalah DONE
	}

	// Validasi apakah target status valid dari status saat ini
	validNext := false
	for _, next := range validNextStatuses[dataQuery.Status] {
		if uaj.Status == next {
			validNext = true
			break
		}
	}

	if !validNext {
		helper.Logger("error", fmt.Sprintf("invalid status transition: from %d to %d", dataQuery.Status, uaj.Status))
		return nil, errors.New("invalid status transition")
	}

	status = statusNames[uaj.Status]

	// DECLINE (4) atau DONE (11)
	if uaj.Status == 4 || uaj.Status == 11 {
		isFinish = 1
	}

	if uaj.Status == 7 {
		fmt.Println("kirim email")
		helper.SendEmail(dataUserFcm.Email, "TJL", status, uaj.Content, "tjl-skck")
	}

	// DONE (11)
	if uaj.Status == 11 {

		queryDepartures := `INSERT INTO departures (content) VALUES (?)`

		resultDepartures, _ := db.DB().Exec(queryDepartures, uaj.Content)

		lastID, errDepartures := resultDepartures.LastInsertId()
		if errDepartures != nil {
			helper.Logger("error", "In Server: "+errDepartures.Error())
			return nil, errors.New(errDepartures.Error())
		}

		queryCandidatePasses := `INSERT INTO candidate_passes (departure_id, apply_job_id, user_candidate_id) VALUES (?, ?, ?)`

		errCandidatePasses := db.Debug().Exec(queryCandidatePasses, lastID, uaj.ApplyJobId, dataQuery.UserId).Error

		if errCandidatePasses != nil {
			helper.Logger("error", "In Server: "+errCandidatePasses.Error())
			return nil, errors.New(errCandidatePasses.Error())
		}

		helper.SendEmail(dataUserFcm.Email, "TJL", "Jadwal Keberangkatan", uaj.Content, "tjl-departure")

		helper.SendFcm("Jadwal Keberangkatan", dataUserFcm.Fullname, dataUserFcm.Token, "notifications", "-")
	}

	if errUserFcmRow != nil {
		if errors.Is(errUserFcmRow, sql.ErrNoRows) {
			helper.Logger("info", "No FCM data found for user")
		}

		helper.Logger("error", "In Server: "+errUserFcmRow.Error())
	}

	query := `UPDATE apply_jobs SET user_confirm_id = ?, status = ?, is_finish = ? WHERE uid = ?`
	err := db.Debug().Exec(query, uaj.UserConfirmId, uaj.Status, isFinish, uaj.ApplyJobId).Error
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	queryHistory := `INSERT INTO apply_job_histories
	(uid, job_id, user_id, user_confirm_id, status, link, schedule)
	VALUES (?, ?, ?, ?, ?, ?, ?)`
	errHistory := db.Debug().Exec(queryHistory,
		dataQuery.Uid, dataQuery.JobId, dataQuery.UserId,
		uaj.UserConfirmId, uaj.Status, uaj.Link, uaj.Schedule,
	).Error

	if errHistory != nil {
		helper.Logger("error", "In Server: "+errHistory.Error())
		return nil, errors.New(errHistory.Error())
	}

	title := fmt.Sprintf("Selamat lamaran Anda sudah dalam tahap [%s]", status)

	if status == "DECLINE" {
		title = fmt.Sprintf("Maaf lamaran Anda telah ditolak [%s]", status)
	}

	helper.SendFcm(title, dataUserFcm.Fullname, dataUserFcm.Token, "apply-job-detail", uaj.ApplyJobId)

	if uaj.IsOffline {

		queryInsertApplyJobOffline := `INSERT INTO apply_job_offlines (apply_job_id, content) VALUES (?, ?)`

		errInsertApplyJobOffline := db.Debug().Exec(queryInsertApplyJobOffline, uaj.ApplyJobId, uaj.Content).Error

		if errInsertApplyJobOffline != nil {
			helper.Logger("error", "In Server: "+errInsertApplyJobOffline.Error())
		}

		queryUserFcm := `SELECT f.token, u.email, p.fullname FROM fcms f
		INNER JOIN profiles p ON p.user_id = f.user_id
		INNER JOIN users u ON p.user_id = u.uid
		WHERE f.user_id = ?`

		rowUserFcm := db.Debug().Raw(queryUserFcm, dataQuery.UserId).Row()

		errUserFcmRow := rowUserFcm.Scan(&dataUserFcm.Token, &dataUserFcm.Email, &dataUserFcm.Fullname)

		if errUserFcmRow != nil {
			if errors.Is(errUserFcmRow, sql.ErrNoRows) {
				helper.Logger("info", "No FCM data found for user")
			}

			helper.Logger("error", "In Server: "+errUserFcmRow.Error())
		}
		message := fmt.Sprintf("Silahkan periksa Alamat E-mail [%s] Anda untuk info lebih lanjut", dataUserFcm.Email)

		helper.SendEmail(dataUserFcm.Email, "TJL", "Lolos Seleksi", uaj.Content, "tjl-apply-job-offline")

		helper.SendFcm(status, message, dataUserFcm.Token, "apply-job-detail", uaj.ApplyJobId)
	}

	queryInsertInbox := `INSERT INTO inboxes (uid, title, caption, user_id, field2, type) VALUES (?, ?, ?, ?, ?, ?)`

	errInsertInbox := db.Debug().Exec(queryInsertInbox, uuid.NewV4().String(), status, uaj.Content, uaj.UserId, uaj.ApplyJobId, "broadcast").Error

	if errInsertInbox != nil {
		helper.Logger("error", "In Server: "+errInsertInbox.Error())
	}

	return map[string]any{}, nil
}

func AdminListApplyJob(branchId, filter string) (map[string]any, error) {

	var job entities.AdminListApplyJobQuery

	var addDoc entities.AdditionalDoc

	var candidateExercise entities.CandidateExerciseQuery
	var candidateBiodata entities.CandidateBiodataQuery
	var candidateLanguage entities.CandidateLanguageQuery
	var candidateWork entities.CandidateWorkQuery
	var candidatePlace entities.CandidatePlaceQuery

	var candidateEdu entities.CandidateEducationQuery

	var exerciseMedia entities.FormExerciseCertificate

	var jobFavourite []entities.JobFavourite

	var dataJob = make([]entities.AdminListApplyJob, 0)

	query := `SELECT aj.uid AS id, j.title, j.caption, j.salary,
	aj.user_id AS user_id_candidate,
	pc.fullname AS user_name_candidate,
	pc.avatar AS user_avatar_candidate,
	upc.email AS user_email_candidate,
	upc.phone AS user_phone_candidate,
	jc.uid as cat_id,
	jc.name AS cat_name, 
	p.id AS place_id,
	p.name AS place_name,
	p.currency AS place_currency,
	p.kurs AS place_kurs,
	p.info AS place_info,
	c.uid AS company_id,
	c.logo AS company_logo,
	c.name AS company_name,
	p.name AS country_name,
	up.user_id,
	up.avatar AS user_avatar,
	up.fullname AS user_name,
	aj.created_at,
	js.id AS job_status_id,
	js.name AS job_status_name,
	b.id AS branch_id,
    b.name AS branch_name
	FROM jobs j
	INNER JOIN job_categories jc ON jc.uid = j.cat_id
	INNER JOIN companies c ON c.uid = j.company_id 
	INNER JOIN apply_jobs aj ON aj.job_id = j.uid
	INNER JOIN job_statuses js ON js.id = aj.status
	INNER JOIN places p ON p.id = j.place_id
	INNER JOIN profiles up ON up.user_id = j.user_id
	INNER JOIN profiles pc ON pc.user_id = aj.user_id
	INNER JOIN users upc ON upc.uid = pc.user_id
	INNER JOIN user_branches ub ON ub.user_id = aj.user_id
	INNER JOIN branchs b ON b.id  = ub.branch_id
	`

	var rows *sql.Rows
	var err error

	conditions := []string{}
	args := []any{}

	if branchId != "" {
		conditions = append(conditions, "ub.branch_id = ?")
		args = append(args, branchId)
	}

	if filter != "ALL" {
		conditions = append(conditions, "js.name = ?")
		args = append(args, filter)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY aj.created_at DESC"

	rows, err = db.Debug().Raw(query, args...).Rows()

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		errJobRows := db.ScanRows(rows, &job)

		if errJobRows != nil {
			helper.Logger("error", "In Server: "+errJobRows.Error())
			return nil, errors.New(errJobRows.Error())
		}

		jobFavourite = nil

		bookmarkQuery := `SELECT job_id, user_id FROM job_favourites WHERE user_id = '` + job.UserId + `' AND job_id = '` + job.Id + `'`

		errBookmark := db.Debug().Raw(bookmarkQuery).Scan(&jobFavourite).Error

		if errBookmark != nil {
			helper.Logger("error", "In Server: "+errBookmark.Error())
			return nil, errors.New(errBookmark.Error())
		}

		isJobFavouriteExist := len(jobFavourite)

		var bookmark bool

		if isJobFavouriteExist == 1 {
			bookmark = true
		} else {
			bookmark = false
		}

		salaryIdr := helper.FormatIDR(job.Salary * job.PlaceKurs)

		// Candidate Exercise

		dataCandidateExercise := make([]entities.CandidateExercise, 0)

		queryCandidateExercise := `SELECT id, name, institution, start_month, start_year, end_month, end_year 
		FROM form_exercises WHERE user_id = ?`

		rowsCandidateExercise, errCandidateExercise := db.Debug().Raw(queryCandidateExercise, job.UserIdCandidate).Rows()

		if errCandidateExercise != nil {
			helper.Logger("error", "In Server: "+errCandidateExercise.Error())
		}
		defer rowsCandidateExercise.Close()

		for rowsCandidateExercise.Next() {
			errCandidateExerciseRows := db.ScanRows(rowsCandidateExercise, &candidateExercise)

			if errCandidateExerciseRows != nil {
				helper.Logger("error", "In Server: "+errCandidateExerciseRows.Error())
				return nil, errors.New(errCandidateExerciseRows.Error())
			}

			queryFormExerciseMedia := `SELECT id, path FROM form_exercise_medias WHERE exercise_id  = ?`

			rowsFormExerciseMedia, errFormExerciseMedia := db.Debug().Raw(queryFormExerciseMedia, candidateExercise.Id).Scan(&exerciseMedia).Rows()

			if errFormExerciseMedia != nil {
				helper.Logger("error", "In Server: "+errFormExerciseMedia.Error())
				return nil, errors.New(errFormExerciseMedia.Error())
			}

			defer rowsFormExerciseMedia.Close()

			var dataFormExerciseCertificate = make([]entities.CandidateExerciseCertificates, 0)

			for rowsFormExerciseMedia.Next() {
				errScanFormExerciseMedia := db.ScanRows(rowsFormExerciseMedia, &exerciseMedia)

				if errScanFormExerciseMedia != nil {
					helper.Logger("error", "In Server: "+errScanFormExerciseMedia.Error())
					return nil, errors.New(errScanFormExerciseMedia.Error())
				}

				dataFormExerciseCertificate = append(dataFormExerciseCertificate, entities.CandidateExerciseCertificates{
					Id:   exerciseMedia.Id,
					Path: exerciseMedia.Path,
				})
			}

			dataCandidateExercise = append(dataCandidateExercise, entities.CandidateExercise{
				Name:         candidateExercise.Name,
				Institution:  candidateExercise.Institution,
				StartMonth:   candidateExercise.StartMonth,
				StartYear:    candidateExercise.StartYear,
				EndMonth:     candidateExercise.EndMonth,
				EndYear:      candidateExercise.EndYear,
				Certificates: dataFormExerciseCertificate,
			})
		}

		// End Candidate Exercise

		// Candidate Biodata

		dataCandidateBiodata := make([]entities.CandidateBiodata, 0)

		queryCandidateBiodata := `SELECT birthdate, gender, weight, height, status, religion, place 
		FROM form_biodatas WHERE user_id = ?`

		rowsCandidateBiodata, errCandidateBiodata := db.Debug().Raw(queryCandidateBiodata, job.UserIdCandidate).Rows()

		if errCandidateBiodata != nil {
			helper.Logger("error", "In Server: "+errCandidateBiodata.Error())
		}
		defer rowsCandidateBiodata.Close()

		for rowsCandidateBiodata.Next() {
			errCandidateBiodataRows := db.ScanRows(rowsCandidateBiodata, &candidateBiodata)

			if errCandidateBiodataRows != nil {
				helper.Logger("error", "In Server: "+errCandidateBiodataRows.Error())
				return nil, errors.New(errCandidateBiodataRows.Error())
			}

			dataCandidateBiodata = append(dataCandidateBiodata, entities.CandidateBiodata{
				Birthdate: candidateBiodata.Birthdate,
				Gender:    candidateBiodata.Gender,
				Weight:    candidateBiodata.Weight,
				Height:    candidateBiodata.Height,
				Status:    candidateBiodata.Status,
				Religion:  candidateBiodata.Religion,
				Place:     candidateBiodata.Place,
			})
		}

		// End Candidate Biodata

		// Candidate Language

		dataCandidateLanguage := make([]entities.CandidateLanguage, 0)

		queryCandidateLanguage := `SELECT level, language 
		FROM form_languages WHERE user_id = ?`

		rowsCandidateLanguage, errCandidateLanguage := db.Debug().Raw(queryCandidateLanguage, job.UserIdCandidate).Rows()

		if errCandidateLanguage != nil {
			helper.Logger("error", "In Server: "+errCandidateLanguage.Error())
		}
		defer rowsCandidateLanguage.Close()

		for rowsCandidateLanguage.Next() {
			errCandidateLanguageRows := db.ScanRows(rowsCandidateLanguage, &candidateLanguage)

			if errCandidateLanguageRows != nil {
				helper.Logger("error", "In Server: "+errCandidateLanguageRows.Error())
				return nil, errors.New(errCandidateLanguageRows.Error())
			}

			dataCandidateLanguage = append(dataCandidateLanguage, entities.CandidateLanguage{
				Level:    candidateLanguage.Level,
				Language: candidateLanguage.Language,
			})
		}

		// End Candidate Language

		// Candidate Work

		dataCandidateWork := make([]entities.CandidateWork, 0)

		queryCandidateWork := `SELECT position, institution, work, country, city, start_month, 
		start_year, end_month, end_year, is_work 
		FROM form_works WHERE user_id = ?`

		rowsCandidateWork, errCandidateWork := db.Debug().Raw(queryCandidateWork, job.UserIdCandidate).Rows()

		if errCandidateWork != nil {
			helper.Logger("error", "In Server: "+errCandidateWork.Error())
		}
		defer rowsCandidateWork.Close()

		for rowsCandidateWork.Next() {
			errCandidateWorkRows := db.ScanRows(rowsCandidateWork, &candidateWork)

			if errCandidateWorkRows != nil {
				helper.Logger("error", "In Server: "+errCandidateWorkRows.Error())
				return nil, errors.New(errCandidateWorkRows.Error())
			}

			dataCandidateWork = append(dataCandidateWork, entities.CandidateWork{
				Position:    candidateWork.Position,
				Institution: candidateWork.Institution,
				Work:        candidateWork.Work,
				Country:     candidateWork.Country,
				City:        candidateWork.City,
				StartMonth:  candidateWork.StartMonth,
				StartYear:   candidateWork.StartYear,
				EndMonth:    candidateWork.EndMonth,
				EndYear:     candidateWork.EndYear,
				IsWork:      candidateWork.IsWork,
			})
		}

		// End Candidate Work

		// Candidate Place

		dataCandidatePlace := make([]entities.CandidatePlace, 0)

		queryCandidatePlace := `SELECT 
		p.name AS province_name, 
		r.name AS city_name, 
		d.name AS district_name, 
		s.name AS subdistrict_name,
		fp.detail_address
		FROM form_places fp 
		INNER JOIN provinces p ON p.id = fp.province_id
		INNER JOIN regencies r ON r.id = fp.city_id
		INNER JOIN districts d ON d.id = fp.district_id
		INNER JOIN villages s ON s.id = fp.subdistrict_id
		WHERE user_id = ?`

		rowsCandidatePlace, errCandidatePlace := db.Debug().Raw(queryCandidatePlace, job.UserIdCandidate).Rows()

		if errCandidatePlace != nil {
			helper.Logger("error", "In Server: "+errCandidatePlace.Error())
		}
		defer rowsCandidatePlace.Close()

		for rowsCandidatePlace.Next() {
			errCandidatePlaceRows := db.ScanRows(rowsCandidatePlace, &candidatePlace)

			if errCandidatePlaceRows != nil {
				helper.Logger("error", "In Server: "+errCandidatePlaceRows.Error())
				return nil, errors.New(errCandidatePlaceRows.Error())
			}

			dataCandidatePlace = append(dataCandidatePlace, entities.CandidatePlace{
				ProvinceName:    candidatePlace.ProvinceName,
				CityName:        candidatePlace.CityName,
				DistrictName:    candidatePlace.DistrictName,
				SubdistrictName: candidatePlace.SubdistrictName,
				DetailAddress:   candidatePlace.DetailAddress,
			})
		}

		// End Candidate Place

		// Candidate Document

		groupedCandidateDocuments := map[string][]entities.CandidateDocument{}

		getTypeLabel := func(docType string) string {
			switch docType {
			case "regulation":
				return "regulation"
			case "preparation":
				return "preparation"
			default:
				return "other"
			}
		}

		queryCandidateDocument := `SELECT d.name AS document, d.type, ajd.path FROM documents d 
		INNER JOIN user_documents ajd ON ajd.type = d.id 
		WHERE ajd.user_id = ?`

		rowsCandidateDocument, errCandidateDocument := db.Debug().Raw(queryCandidateDocument, job.UserIdCandidate).Rows()
		if errCandidateDocument != nil {
			helper.Logger("error", "In Server: "+errCandidateDocument.Error())
		}
		defer rowsCandidateDocument.Close()

		for rowsCandidateDocument.Next() {
			var candidateDoc entities.CandidateDocument
			errCandidateDocumentRows := db.ScanRows(rowsCandidateDocument, &candidateDoc)

			if errCandidateDocumentRows != nil {
				helper.Logger("error", "In Server: "+errCandidateDocumentRows.Error())
				return nil, errors.New(errCandidateDocumentRows.Error())
			}

			label := getTypeLabel(candidateDoc.Type)
			groupedCandidateDocuments[label] = append(groupedCandidateDocuments[label], entities.CandidateDocument{
				Document: candidateDoc.Document,
				Type:     candidateDoc.Type,
				Path:     candidateDoc.Path,
			})
		}

		queryAdditionalDoc := `SELECT path, type
		FROM user_document_additionals WHERE user_id = ?`

		rowsAdditionalDoc, errAdditionalDoc := db.Debug().Raw(queryAdditionalDoc, job.UserIdCandidate).Rows()

		if errAdditionalDoc != nil {
			helper.Logger("error", "In Server: "+errAdditionalDoc.Error())
		}
		defer rowsAdditionalDoc.Close()

		hasMedicalLicense := false

		for rowsAdditionalDoc.Next() {
			var additionalDoc entities.AdditionalDoc
			errScan := db.ScanRows(rowsAdditionalDoc, &additionalDoc)
			if errScan != nil {
				helper.Logger("error", "In Server: "+errScan.Error())
				return nil, errors.New(errScan.Error())
			}

			if additionalDoc.Type == "medical-license" && !hasMedicalLicense {
				groupedCandidateDocuments["regulation"] = append(
					groupedCandidateDocuments["regulation"],
					entities.CandidateDocument{
						Document: "Medical License",
						Type:     "regulation",
						Path:     additionalDoc.Path,
					},
				)

				hasMedicalLicense = true
			}
		}

		// End Candidate Document

		// Candidate Education

		dataCandidateEducation := make([]entities.CandidateEducation, 0)

		queryCandidateEducation := `SELECT education_level AS edu, major, school_or_college, start_month, start_year, end_month, end_year 
		FROM form_educations WHERE user_id = ?`

		rowsCandidateEducation, errCandidateEducation := db.Debug().Raw(queryCandidateEducation, job.UserIdCandidate).Rows()

		if errCandidateEducation != nil {
			helper.Logger("error", "In Server: "+errCandidateEducation.Error())
		}
		defer rowsCandidateEducation.Close()

		for rowsCandidateEducation.Next() {
			errCandidateEducationRows := db.ScanRows(rowsCandidateEducation, &candidateEdu)

			if errCandidateEducationRows != nil {
				helper.Logger("error", "In Server: "+errCandidateEducationRows.Error())
				return nil, errors.New(errCandidateEducationRows.Error())
			}

			dataCandidateEducation = append(dataCandidateEducation, entities.CandidateEducation{
				EducationalLevel: candidateEdu.Edu,
				Major:            candidateEdu.Major,
				SchoolOrCollege:  candidateEdu.SchoolOrCollege,
				StartMonth:       candidateEdu.StartMonth,
				EndMonth:         candidateEdu.EndMonth,
				StartYear:        candidateEdu.StartYear,
				EndYear:          candidateEdu.EndYear,
			})
		}

		// End Candidate Education

		// Additionaldoc

		dataAddDoc := make([]entities.AdditionalDoc, 0)

		queryAddDoc := `SELECT path, type
		FROM user_document_additionals WHERE user_id = ? AND (type = ? OR type = ?)`

		rowsAddDoc, errAddDoc := db.Debug().Raw(queryAddDoc, job.UserIdCandidate, "introduction-video", "introduction-photo").Rows()

		if errAddDoc != nil {
			helper.Logger("error", "In Server: "+errAddDoc.Error())
		}
		defer rowsAddDoc.Close()

		for rowsAddDoc.Next() {
			errAddDoc := db.ScanRows(rowsAddDoc, &addDoc)

			if errAddDoc != nil {
				helper.Logger("error", "In Server: "+errAddDoc.Error())
				return nil, errors.New(errAddDoc.Error())
			}

			dataAddDoc = append(dataAddDoc, entities.AdditionalDoc{
				Path: addDoc.Path,
				Type: addDoc.Type,
			})
		}

		// End Additional Doc

		dataJob = append(dataJob, entities.AdminListApplyJob{
			Id:        job.Id,
			Title:     job.Title,
			Caption:   job.Caption,
			Salary:    int(job.Salary),
			SalaryIDR: salaryIdr,
			Bookmark:  bookmark,
			Company: entities.JobCompany{
				Id:      job.CompanyId,
				Logo:    job.CompanyLogo,
				Name:    job.CompanyName,
				Country: job.CountryName,
			},
			Candidate: entities.Candidate{
				Id:                job.UserIdCandidate,
				Email:             job.UserEmailCandidate,
				Avatar:            job.UserAvatarCandidate,
				Name:              job.UserNameCandidate,
				Phone:             job.UserPhoneCandidate,
				CandidateExercise: dataCandidateExercise,
				CandidateBiodata:  dataCandidateBiodata,
				CandidateLanguage: dataCandidateLanguage,
				CandidateWork:     dataCandidateWork,
				CandidatePlace:    dataCandidatePlace,
				CandidateEdu:      dataCandidateEducation,
				AdditionalDoc:     dataAddDoc,
				Document:          groupedCandidateDocuments,
			},
			Status: entities.JobStatus{
				Id:   job.JobStatusId,
				Name: job.JobStatusName,
			},
			JobCategory: entities.JobCategory{
				Id:   job.CatId,
				Name: job.CatName,
			},
			JobPlace: entities.JobPlace{
				Id:       job.PlaceId,
				Name:     job.PlaceName,
				Currency: job.PlaceCurrency,
				Kurs:     int(job.PlaceKurs),
				Info:     job.PlaceInfo,
			},
			Author: entities.AuthorJobUser{
				Id:     job.UserId,
				Avatar: job.UserAvatar,
				Name:   job.UserName,
			},
			Branch: entities.AdminBranch{
				Id:   job.BranchId,
				Name: job.BranchName,
			},
			Created: job.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return map[string]any{
		"data": dataJob,
	}, nil
}

func AdminDetailApplyJob(id string) (map[string]any, error) {

	var job entities.AdminListApplyJobQuery

	var addDoc entities.AdditionalDoc

	var candidateExercise entities.CandidateExerciseQuery
	var candidateBiodata entities.CandidateBiodataQuery
	var candidateLanguage entities.CandidateLanguageQuery
	var candidateWork entities.CandidateWorkQuery
	var candidatePlace entities.CandidatePlaceQuery

	var candidateEdu entities.CandidateEducationQuery

	var exerciseMedia entities.FormExerciseCertificate

	var jobFavourite []entities.JobFavourite

	var dataJob entities.AdminListApplyJob

	query := `SELECT aj.uid AS id, j.title, j.caption, j.salary,
	aj.user_id AS user_id_candidate,
	pc.fullname AS user_name_candidate,
	pc.avatar AS user_avatar_candidate,
	upc.email AS user_email_candidate,
	upc.phone AS user_phone_candidate,
	jc.uid as cat_id,
	jc.name AS cat_name, 
	p.id AS place_id,
	p.name AS place_name,
	p.currency AS place_currency,
	p.kurs AS place_kurs,
	p.info AS place_info,
	c.uid AS company_id,
	c.logo AS company_logo,
	c.name AS company_name,
	p.name AS country_name,
	up.user_id,
	up.avatar AS user_avatar,
	up.fullname AS user_name,
	aj.created_at,
	js.id AS job_status_id,
	js.name AS job_status_name,
	b.id AS branch_id,
    b.name AS branch_name
	FROM jobs j
	INNER JOIN job_categories jc ON jc.uid = j.cat_id
	INNER JOIN companies c ON c.uid = j.company_id 
	INNER JOIN apply_jobs aj ON aj.job_id = j.uid
	INNER JOIN job_statuses js ON js.id = aj.status
	INNER JOIN places p ON p.id = j.place_id
	INNER JOIN profiles up ON up.user_id = j.user_id
	INNER JOIN profiles pc ON pc.user_id = aj.user_id
	INNER JOIN users upc ON upc.uid = pc.user_id
	INNER JOIN user_branches ub ON ub.user_id = aj.user_id
	INNER JOIN branchs b ON b.id  = ub.branch_id
	WHERE aj.uid = ?
	`

	var rows *sql.Rows
	var err error

	query += " ORDER BY aj.created_at DESC"

	rows, err = db.Debug().Raw(query, id).Rows()

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		errJobRows := db.ScanRows(rows, &job)

		if errJobRows != nil {
			helper.Logger("error", "In Server: "+errJobRows.Error())
			return nil, errors.New(errJobRows.Error())
		}

		jobFavourite = nil

		bookmarkQuery := `SELECT job_id, user_id FROM job_favourites WHERE user_id = '` + job.UserId + `' AND job_id = '` + job.Id + `'`

		errBookmark := db.Debug().Raw(bookmarkQuery).Scan(&jobFavourite).Error

		if errBookmark != nil {
			helper.Logger("error", "In Server: "+errBookmark.Error())
			return nil, errors.New(errBookmark.Error())
		}

		isJobFavouriteExist := len(jobFavourite)

		var bookmark bool

		if isJobFavouriteExist == 1 {
			bookmark = true
		} else {
			bookmark = false
		}

		salaryIdr := helper.FormatIDR(job.Salary * job.PlaceKurs)

		// Candidate Exercise

		dataCandidateExercise := make([]entities.CandidateExercise, 0)

		queryCandidateExercise := `SELECT id, name, institution, start_month, start_year, end_month, end_year 
		FROM form_exercises WHERE user_id = ?`

		rowsCandidateExercise, errCandidateExercise := db.Debug().Raw(queryCandidateExercise, job.UserIdCandidate).Rows()

		if errCandidateExercise != nil {
			helper.Logger("error", "In Server: "+errCandidateExercise.Error())
		}
		defer rowsCandidateExercise.Close()

		for rowsCandidateExercise.Next() {
			errCandidateExerciseRows := db.ScanRows(rowsCandidateExercise, &candidateExercise)

			if errCandidateExerciseRows != nil {
				helper.Logger("error", "In Server: "+errCandidateExerciseRows.Error())
				return nil, errors.New(errCandidateExerciseRows.Error())
			}

			queryFormExerciseMedia := `SELECT id, path FROM form_exercise_medias WHERE exercise_id  = ?`

			rowsFormExerciseMedia, errFormExerciseMedia := db.Debug().Raw(queryFormExerciseMedia, candidateExercise.Id).Scan(&exerciseMedia).Rows()

			if errFormExerciseMedia != nil {
				helper.Logger("error", "In Server: "+errFormExerciseMedia.Error())
				return nil, errors.New(errFormExerciseMedia.Error())
			}

			defer rowsFormExerciseMedia.Close()

			var dataFormExerciseCertificate = make([]entities.CandidateExerciseCertificates, 0)

			for rowsFormExerciseMedia.Next() {
				errScanFormExerciseMedia := db.ScanRows(rowsFormExerciseMedia, &exerciseMedia)

				if errScanFormExerciseMedia != nil {
					helper.Logger("error", "In Server: "+errScanFormExerciseMedia.Error())
					return nil, errors.New(errScanFormExerciseMedia.Error())
				}

				dataFormExerciseCertificate = append(dataFormExerciseCertificate, entities.CandidateExerciseCertificates{
					Id:   exerciseMedia.Id,
					Path: exerciseMedia.Path,
				})
			}

			dataCandidateExercise = append(dataCandidateExercise, entities.CandidateExercise{
				Name:         candidateExercise.Name,
				Institution:  candidateExercise.Institution,
				StartMonth:   candidateExercise.StartMonth,
				StartYear:    candidateExercise.StartYear,
				EndMonth:     candidateExercise.EndMonth,
				EndYear:      candidateExercise.EndYear,
				Certificates: dataFormExerciseCertificate,
			})
		}

		// End Candidate Exercise

		// Candidate Biodata

		dataCandidateBiodata := make([]entities.CandidateBiodata, 0)

		queryCandidateBiodata := `SELECT birthdate, gender, weight, height, status, religion, place 
		FROM form_biodatas WHERE user_id = ?`

		rowsCandidateBiodata, errCandidateBiodata := db.Debug().Raw(queryCandidateBiodata, job.UserIdCandidate).Rows()

		if errCandidateBiodata != nil {
			helper.Logger("error", "In Server: "+errCandidateBiodata.Error())
		}
		defer rowsCandidateBiodata.Close()

		for rowsCandidateBiodata.Next() {
			errCandidateBiodataRows := db.ScanRows(rowsCandidateBiodata, &candidateBiodata)

			if errCandidateBiodataRows != nil {
				helper.Logger("error", "In Server: "+errCandidateBiodataRows.Error())
				return nil, errors.New(errCandidateBiodataRows.Error())
			}

			dataCandidateBiodata = append(dataCandidateBiodata, entities.CandidateBiodata{
				Birthdate: candidateBiodata.Birthdate,
				Gender:    candidateBiodata.Gender,
				Weight:    candidateBiodata.Weight,
				Height:    candidateBiodata.Height,
				Status:    candidateBiodata.Status,
				Religion:  candidateBiodata.Religion,
				Place:     candidateBiodata.Place,
			})
		}

		// End Candidate Biodata

		// Candidate Language

		dataCandidateLanguage := make([]entities.CandidateLanguage, 0)

		queryCandidateLanguage := `SELECT level, language 
		FROM form_languages WHERE user_id = ?`

		rowsCandidateLanguage, errCandidateLanguage := db.Debug().Raw(queryCandidateLanguage, job.UserIdCandidate).Rows()

		if errCandidateLanguage != nil {
			helper.Logger("error", "In Server: "+errCandidateLanguage.Error())
		}
		defer rowsCandidateLanguage.Close()

		for rowsCandidateLanguage.Next() {
			errCandidateLanguageRows := db.ScanRows(rowsCandidateLanguage, &candidateLanguage)

			if errCandidateLanguageRows != nil {
				helper.Logger("error", "In Server: "+errCandidateLanguageRows.Error())
				return nil, errors.New(errCandidateLanguageRows.Error())
			}

			dataCandidateLanguage = append(dataCandidateLanguage, entities.CandidateLanguage{
				Level:    candidateLanguage.Level,
				Language: candidateLanguage.Language,
			})
		}

		// End Candidate Language

		// Candidate Work

		dataCandidateWork := make([]entities.CandidateWork, 0)

		queryCandidateWork := `SELECT position, institution, work, country, city, start_month, 
		start_year, end_month, end_year, is_work 
		FROM form_works WHERE user_id = ?`

		rowsCandidateWork, errCandidateWork := db.Debug().Raw(queryCandidateWork, job.UserIdCandidate).Rows()

		if errCandidateWork != nil {
			helper.Logger("error", "In Server: "+errCandidateWork.Error())
		}
		defer rowsCandidateWork.Close()

		for rowsCandidateWork.Next() {
			errCandidateWorkRows := db.ScanRows(rowsCandidateWork, &candidateWork)

			if errCandidateWorkRows != nil {
				helper.Logger("error", "In Server: "+errCandidateWorkRows.Error())
				return nil, errors.New(errCandidateWorkRows.Error())
			}

			dataCandidateWork = append(dataCandidateWork, entities.CandidateWork{
				Position:    candidateWork.Position,
				Institution: candidateWork.Institution,
				Work:        candidateWork.Work,
				Country:     candidateWork.Country,
				City:        candidateWork.City,
				StartMonth:  candidateWork.StartMonth,
				StartYear:   candidateWork.StartYear,
				EndMonth:    candidateWork.EndMonth,
				EndYear:     candidateWork.EndYear,
				IsWork:      candidateWork.IsWork,
			})
		}

		// End Candidate Work

		// Candidate Place

		dataCandidatePlace := make([]entities.CandidatePlace, 0)

		queryCandidatePlace := `SELECT 
		p.name AS province_name, 
		r.name AS city_name, 
		d.name AS district_name, 
		s.name AS subdistrict_name,
		fp.detail_address
		FROM form_places fp 
		INNER JOIN provinces p ON p.id = fp.province_id
		INNER JOIN regencies r ON r.id = fp.city_id
		INNER JOIN districts d ON d.id = fp.district_id
		INNER JOIN villages s ON s.id = fp.subdistrict_id
		WHERE user_id = ?`

		rowsCandidatePlace, errCandidatePlace := db.Debug().Raw(queryCandidatePlace, job.UserIdCandidate).Rows()

		if errCandidatePlace != nil {
			helper.Logger("error", "In Server: "+errCandidatePlace.Error())
		}
		defer rowsCandidatePlace.Close()

		for rowsCandidatePlace.Next() {
			errCandidatePlaceRows := db.ScanRows(rowsCandidatePlace, &candidatePlace)

			if errCandidatePlaceRows != nil {
				helper.Logger("error", "In Server: "+errCandidatePlaceRows.Error())
				return nil, errors.New(errCandidatePlaceRows.Error())
			}

			dataCandidatePlace = append(dataCandidatePlace, entities.CandidatePlace{
				ProvinceName:    candidatePlace.ProvinceName,
				CityName:        candidatePlace.CityName,
				DistrictName:    candidatePlace.DistrictName,
				SubdistrictName: candidatePlace.SubdistrictName,
				DetailAddress:   candidatePlace.DetailAddress,
			})
		}

		// End Candidate Place

		// Candidate Document

		groupedCandidateDocuments := map[string][]entities.CandidateDocument{}

		getTypeLabel := func(docType string) string {
			switch docType {
			case "regulation":
				return "regulation"
			case "preparation":
				return "preparation"
			default:
				return "other"
			}
		}

		queryCandidateDocument := `SELECT d.name AS document, d.type, ajd.path FROM documents d 
		INNER JOIN user_documents ajd ON ajd.type = d.id 
		WHERE ajd.user_id = ?`

		rowsCandidateDocument, errCandidateDocument := db.Debug().Raw(queryCandidateDocument, job.UserIdCandidate).Rows()
		if errCandidateDocument != nil {
			helper.Logger("error", "In Server: "+errCandidateDocument.Error())
		}
		defer rowsCandidateDocument.Close()

		for rowsCandidateDocument.Next() {
			var candidateDoc entities.CandidateDocument
			errCandidateDocumentRows := db.ScanRows(rowsCandidateDocument, &candidateDoc)

			if errCandidateDocumentRows != nil {
				helper.Logger("error", "In Server: "+errCandidateDocumentRows.Error())
				return nil, errors.New(errCandidateDocumentRows.Error())
			}

			label := getTypeLabel(candidateDoc.Type)
			groupedCandidateDocuments[label] = append(groupedCandidateDocuments[label], entities.CandidateDocument{
				Document: candidateDoc.Document,
				Type:     candidateDoc.Type,
				Path:     candidateDoc.Path,
			})
		}

		queryAdditionalDoc := `SELECT path, type
		FROM user_document_additionals WHERE user_id = ?`

		rowsAdditionalDoc, errAdditionalDoc := db.Debug().Raw(queryAdditionalDoc, job.UserIdCandidate).Rows()

		if errAdditionalDoc != nil {
			helper.Logger("error", "In Server: "+errAdditionalDoc.Error())
		}
		defer rowsAdditionalDoc.Close()

		hasMedicalLicense := false

		for rowsAdditionalDoc.Next() {
			var additionalDoc entities.AdditionalDoc
			errScan := db.ScanRows(rowsAdditionalDoc, &additionalDoc)
			if errScan != nil {
				helper.Logger("error", "In Server: "+errScan.Error())
				return nil, errors.New(errScan.Error())
			}

			if additionalDoc.Type == "medical-license" && !hasMedicalLicense {
				groupedCandidateDocuments["regulation"] = append(
					groupedCandidateDocuments["regulation"],
					entities.CandidateDocument{
						Document: "Medical License",
						Type:     "regulation",
						Path:     additionalDoc.Path,
					},
				)

				hasMedicalLicense = true
			}
		}

		// End Candidate Document

		// Candidate Education

		dataCandidateEducation := make([]entities.CandidateEducation, 0)

		queryCandidateEducation := `SELECT education_level AS edu, major, school_or_college, start_month, start_year, end_month, end_year 
		FROM form_educations WHERE user_id = ?`

		rowsCandidateEducation, errCandidateEducation := db.Debug().Raw(queryCandidateEducation, job.UserIdCandidate).Rows()

		if errCandidateEducation != nil {
			helper.Logger("error", "In Server: "+errCandidateEducation.Error())
		}
		defer rowsCandidateEducation.Close()

		for rowsCandidateEducation.Next() {
			errCandidateEducationRows := db.ScanRows(rowsCandidateEducation, &candidateEdu)

			if errCandidateEducationRows != nil {
				helper.Logger("error", "In Server: "+errCandidateEducationRows.Error())
				return nil, errors.New(errCandidateEducationRows.Error())
			}

			dataCandidateEducation = append(dataCandidateEducation, entities.CandidateEducation{
				EducationalLevel: candidateEdu.Edu,
				Major:            candidateEdu.Major,
				SchoolOrCollege:  candidateEdu.SchoolOrCollege,
				StartMonth:       candidateEdu.StartMonth,
				EndMonth:         candidateEdu.EndMonth,
				StartYear:        candidateEdu.StartYear,
				EndYear:          candidateEdu.EndYear,
			})
		}

		// End Candidate Education

		// Additionaldoc

		dataAddDoc := make([]entities.AdditionalDoc, 0)

		queryAddDoc := `SELECT path, type
		FROM user_document_additionals WHERE user_id = ? AND (type = ? OR type = ?)`

		rowsAddDoc, errAddDoc := db.Debug().Raw(queryAddDoc, job.UserIdCandidate, "introduction-video", "introduction-photo").Rows()

		if errAddDoc != nil {
			helper.Logger("error", "In Server: "+errAddDoc.Error())
		}
		defer rowsAddDoc.Close()

		for rowsAddDoc.Next() {
			errAddDoc := db.ScanRows(rowsAddDoc, &addDoc)

			if errAddDoc != nil {
				helper.Logger("error", "In Server: "+errAddDoc.Error())
				return nil, errors.New(errAddDoc.Error())
			}

			dataAddDoc = append(dataAddDoc, entities.AdditionalDoc{
				Path: addDoc.Path,
				Type: addDoc.Type,
			})
		}

		// End Additional Doc

		dataJob = entities.AdminListApplyJob{
			Id:        job.Id,
			Title:     job.Title,
			Caption:   job.Caption,
			Salary:    int(job.Salary),
			SalaryIDR: salaryIdr,
			Bookmark:  bookmark,
			Company: entities.JobCompany{
				Id:      job.CompanyId,
				Logo:    job.CompanyLogo,
				Name:    job.CompanyName,
				Country: job.CountryName,
			},
			Candidate: entities.Candidate{
				Id:                job.UserIdCandidate,
				Email:             job.UserEmailCandidate,
				Avatar:            job.UserAvatarCandidate,
				Name:              job.UserNameCandidate,
				Phone:             job.UserPhoneCandidate,
				CandidateExercise: dataCandidateExercise,
				CandidateBiodata:  dataCandidateBiodata,
				CandidateLanguage: dataCandidateLanguage,
				CandidateWork:     dataCandidateWork,
				CandidatePlace:    dataCandidatePlace,
				CandidateEdu:      dataCandidateEducation,
				AdditionalDoc:     dataAddDoc,
				Document:          groupedCandidateDocuments,
			},
			Status: entities.JobStatus{
				Id:   job.JobStatusId,
				Name: job.JobStatusName,
			},
			JobCategory: entities.JobCategory{
				Id:   job.CatId,
				Name: job.CatName,
			},
			JobPlace: entities.JobPlace{
				Id:       job.PlaceId,
				Name:     job.PlaceName,
				Currency: job.PlaceCurrency,
				Kurs:     int(job.PlaceKurs),
				Info:     job.PlaceInfo,
			},
			Author: entities.AuthorJobUser{
				Id:     job.UserId,
				Avatar: job.UserAvatar,
				Name:   job.UserName,
			},
			Branch: entities.AdminBranch{
				Id:   job.BranchId,
				Name: job.BranchName,
			},
			Created: job.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return map[string]any{
		"data": dataJob,
	}, nil
}

func JobList(userId, search, salary, country, position, page, limit string, isRecommendation bool) (map[string]any, error) {
	url := os.Getenv("API_URL_PROD")

	var allJob []models.AllJob
	var job entities.JobListQuery
	var jobFavourite []entities.JobFavourite
	var jobSkillCategory []entities.JobSkillCategory
	var dataJob = make([]entities.JobList, 0)

	pageinteger, _ := strconv.Atoi(page)
	limitinteger, _ := strconv.Atoi(limit)

	var offset = strconv.Itoa((pageinteger - 1) * limitinteger)

	errAllJob := db.Debug().Raw(`SELECT uid FROM jobs`).Scan(&allJob).Error

	if errAllJob != nil {
		helper.Logger("error", "In Server: "+errAllJob.Error())
	}

	var resultTotal = len(allJob)

	var perPage = math.Ceil(float64(resultTotal) / float64(limitinteger))

	var prevPage int
	var nextPage int

	if pageinteger == 1 {
		prevPage = 1
	} else {
		prevPage = pageinteger - 1
	}

	nextPage = pageinteger + 1

	query := `SELECT 
	j.uid AS id,
	j.title,
	j.caption,
	j.salary,
	j.min_salary,
	j.max_salary,
	j.worker_count,
	jc.uid AS cat_id,
	jc.logo AS cat_icon,
	jc.name AS cat_name, 
	tj.name AS cat_type,
	p.id AS place_id,
	p.name AS place_name,
	p.currency AS place_currency,
	p.kurs AS place_kurs,
	p.info AS place_info,
	p.symbol as place_symbol,
	p.language_code AS place_language_code,
	c.uid AS company_id,
	c.logo AS company_logo,
	c.name AS company_name,
	up.user_id,
	up.avatar AS user_avatar,
	up.fullname AS user_name,
	j.created_at,
	j.salary,
	(j.salary * p.kurs) AS salary_idr,
	skills.skill_names
	FROM jobs j
	INNER JOIN job_categories jc ON jc.uid = j.cat_id
	INNER JOIN type_jobs tj ON tj.id = jc.type
	INNER JOIN companies c ON c.uid = j.company_id 
	INNER JOIN places p ON p.id = j.place_id
	INNER JOIN profiles up ON up.user_id = j.user_id
	LEFT JOIN (
		SELECT js.job_id, GROUP_CONCAT(jsc.name) AS skill_names
		FROM job_skills js
		JOIN job_skill_categories jsc ON jsc.uid = js.cat_id
		GROUP BY js.job_id
	) skills ON skills.job_id = j.uid
	WHERE p.name LIKE ?
	AND jc.name LIKE ?
	AND (
		jc.name LIKE ? OR 
		c.name LIKE ? OR 
		p.name LIKE ? OR 
		skills.skill_names LIKE ?
	)
	`

	params := []any{
		"%" + country + "%",
		"%" + position + "%",
		"%" + search + "%",
		"%" + search + "%",
		"%" + search + "%",
		"%" + search + "%",
	}

	if salary != "" {
		query += ` AND (j.salary * p.kurs) >= ? `
		params = append(params, salary)
	}

	if isRecommendation {
		query += ` ORDER BY j.clicked_count DESC`
	} else {
		query += ` ORDER BY j.created_at DESC`
	}

	query += ` LIMIT ?, ?`
	params = append(params, offset, limit)

	rows, err := db.Debug().Raw(query, params...).Rows()
	if err != nil {
		log.Println("Query Error:", err)
	}
	defer rows.Close()

	for rows.Next() {
		errJobRows := db.ScanRows(rows, &job)

		if errJobRows != nil {
			helper.Logger("error", "In Server: "+errJobRows.Error())
			return nil, errors.New(errJobRows.Error())
		}

		bookmarkQuery := `SELECT job_id, user_id FROM job_favourites WHERE user_id = ? AND job_id = ?`
		errBookmark := db.Debug().Raw(bookmarkQuery, userId, job.Id).Scan(&jobFavourite).Error

		if errBookmark != nil {
			helper.Logger("error", "In Server: "+errBookmark.Error())
			return nil, errors.New(errBookmark.Error())
		}

		isJobFavouriteExist := len(jobFavourite)
		bookmark := isJobFavouriteExist == 1

		salaryIdr := helper.FormatIDR(job.Salary * job.PlaceKurs)
		minSalaryIdr := helper.FormatIDR(job.MinSalary * job.PlaceKurs)
		maxSalaryIdr := helper.FormatIDR(job.MaxSalary * job.PlaceKurs)

		jobSkillsQuery := `SELECT jsc.uid AS id, jsc.name
		FROM job_skills js 
		INNER JOIN job_skill_categories jsc ON jsc.uid = js.cat_id
		WHERE js.job_id = ?`
		errJobSkill := db.Debug().Raw(jobSkillsQuery, job.Id).Scan(&jobSkillCategory).Error

		if errJobSkill != nil {
			helper.Logger("error", "In Server: "+errJobSkill.Error())
			return nil, errors.New(errJobSkill.Error())
		}

		dataJob = append(dataJob, entities.JobList{
			Id:      job.Id,
			Title:   job.Title,
			Caption: job.Caption,
			Skills:  jobSkillCategory,
			Company: entities.JobCompany{
				Id:      job.CompanyId,
				Logo:    job.CompanyLogo,
				Name:    job.CompanyName,
				Country: job.PlaceName,
			},
			WorkerCount:  job.WorkerCount,
			Salary:       int(job.Salary),
			MinSalary:    int(job.MinSalary),
			MaxSalary:    int(job.MaxSalary),
			MinSalaryIDR: minSalaryIdr,
			MaxSalaryIDR: maxSalaryIdr,
			SalaryIDR:    salaryIdr,
			Bookmark:     bookmark,
			JobCategory: entities.JobCategory{
				Id:   job.CatId,
				Icon: job.CatIcon,
				Name: job.CatName,
				Type: job.CatType,
			},
			JobPlace: entities.JobPlace{
				Id:           job.PlaceId,
				Name:         job.PlaceName,
				Currency:     job.PlaceCurrency,
				LanguageCode: job.PlaceLanguageCode,
				Symbol:       job.PlaceSymbol,
				Kurs:         int(job.PlaceKurs),
				Info:         job.PlaceInfo,
			},
			Author: entities.AuthorJobUser{
				Id:     job.UserId,
				Avatar: job.UserAvatar,
				Name:   job.UserName,
			},
			Created: job.CreatedAt.Format("2006-01-02 15:04:05"),
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
		"next_url":     url + "/api/v1/job?page=" + nextUrl + "&limit=10",
		"prev_url":     url + "/api/v1/job?page=" + prevUrl + "&limit=10",
		"data":         &dataJob,
	}, nil
}

func JobDetail(j *models.Job) (map[string]any, error) {
	var job entities.JobListQuery
	var jobFavourite []entities.JobFavourite
	var jobSkillCategory []entities.JobSkillCategory

	var dataJob = make([]entities.JobList, 0)

	query := `
		SELECT j.uid AS id, j.title, j.caption, j.salary, j.min_salary, j.max_salary, j.worker_count,
		jc.uid as cat_id,
		jc.name AS cat_name, 
		tj.name AS cat_type,
		p.id AS place_id,
		p.name AS place_name,
		p.currency AS place_currency,
		p.kurs AS place_kurs,
		p.info AS place_info,
		p.symbol AS place_symbol,
		p.language_code AS place_language_code,
		c.uid AS company_id,
		c.logo AS company_logo,
		c.name AS company_name,
		up.user_id,
		up.avatar AS user_avatar,
		up.fullname AS user_name,
		j.created_at
		FROM jobs j
		INNER JOIN job_categories jc ON jc.uid = j.cat_id
		INNER JOIN type_jobs tj ON tj.id = jc.type
		INNER JOIN companies c ON c.uid = j.company_id 
		INNER JOIN places p ON p.id = j.place_id
		INNER JOIN profiles up ON up.user_id = j.user_id
		WHERE j.uid = ?
	`

	rows, err := db.Debug().Raw(query, j.Id).Rows()

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}
	defer rows.Close()

	errUpdateClickedCount := db.Debug().Exec(`UPDATE jobs SET clicked_count = clicked_count + 1 WHERE uid = ?`, j.Id).Error

	if errUpdateClickedCount != nil {
		helper.Logger("error", "In Server: "+errUpdateClickedCount.Error())
		return nil, errors.New(errUpdateClickedCount.Error())
	}

	found := false
	for rows.Next() {
		found = true
		errJobRows := db.ScanRows(rows, &job)

		if errJobRows != nil {
			helper.Logger("error", "In Server: "+errJobRows.Error())
			return nil, errors.New(errJobRows.Error())
		}

		bookmarkQuery := `SELECT job_id, user_id FROM job_favourites WHERE user_id = ? AND job_id = ?`
		errBookmark := db.Debug().Raw(bookmarkQuery, j.UserId, job.Id).Scan(&jobFavourite).Error

		if errBookmark != nil {
			helper.Logger("error", "In Server: "+errBookmark.Error())
			return nil, errors.New(errBookmark.Error())
		}

		isJobFavouriteExist := len(jobFavourite)
		bookmark := isJobFavouriteExist == 1

		salaryIdr := helper.FormatIDR(job.Salary * job.PlaceKurs)
		minSalaryIdr := helper.FormatIDR(job.MinSalary * job.PlaceKurs)
		maxSalaryIdr := helper.FormatIDR(job.MaxSalary * job.PlaceKurs)

		jobSkillsQuery := `SELECT jsc.uid AS id, jsc.name
		FROM job_skills js 
		INNER JOIN job_skill_categories jsc ON jsc.uid = js.cat_id
		WHERE js.job_id = ?`
		errJobSkill := db.Debug().Raw(jobSkillsQuery, job.Id).Scan(&jobSkillCategory).Error

		if errJobSkill != nil {
			helper.Logger("error", "In Server: "+errJobSkill.Error())
			return nil, errors.New(errJobSkill.Error())
		}

		dataJob = append(dataJob, entities.JobList{
			Id:      job.Id,
			Title:   job.Title,
			Caption: job.Caption,
			Skills:  jobSkillCategory,
			Company: entities.JobCompany{
				Id:      job.CompanyId,
				Logo:    job.CompanyLogo,
				Name:    job.CompanyName,
				Country: job.PlaceName,
			},
			WorkerCount:  job.WorkerCount,
			Salary:       int(job.Salary),
			MinSalary:    int(job.MinSalary),
			MaxSalary:    int(job.MaxSalary),
			MinSalaryIDR: minSalaryIdr,
			MaxSalaryIDR: maxSalaryIdr,
			SalaryIDR:    salaryIdr,
			Bookmark:     bookmark,
			JobCategory: entities.JobCategory{
				Id:   job.CatId,
				Icon: job.CatIcon,
				Name: job.CatName,
				Type: job.CatType,
			},
			JobPlace: entities.JobPlace{
				Id:           job.PlaceId,
				Name:         job.PlaceName,
				Currency:     job.PlaceCurrency,
				Kurs:         int(job.PlaceKurs),
				Symbol:       job.PlaceSymbol,
				LanguageCode: job.PlaceLanguageCode,
				Info:         job.PlaceInfo,
			},
			Author: entities.AuthorJobUser{
				Id:     job.UserId,
				Avatar: job.UserAvatar,
				Name:   job.UserName,
			},
			Created: job.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	if !found {
		return nil, errors.New("job not found")
	}

	return map[string]any{
		"data": dataJob[0],
	}, nil
}

func JobPlace() (map[string]any, error) {
	places := []entities.JobPlace{}

	query := `SELECT id, name, currency, kurs, info FROM places`

	err := db.Debug().Raw(query).Scan(&places).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}

	return map[string]any{
		"data": places,
	}, nil
}

func JobStore(j *models.JobStore) (map[string]any, error) {

	j.Id = uuid.NewV4().String()

	users := []entities.User{}
	categories := []entities.JobCategory{}
	places := []entities.JobPlace{}

	checkQueryUser := `SELECT uid FROM users WHERE uid = '` + j.UserId + `'`

	errCheckUser := db.Debug().Raw(checkQueryUser).Scan(&users).Error

	if errCheckUser != nil {
		helper.Logger("error", "In Server: "+errCheckUser.Error())
		return nil, errors.New(errCheckUser.Error())
	}

	isUserExist := len(users)

	if isUserExist == 0 {
		return nil, errors.New("USER_NOT_FOUND")
	}

	checkQueryCat := `SELECT uid AS id FROM job_categories WHERE uid = '` + j.CatId + `'`

	errCheckCat := db.Debug().Raw(checkQueryCat).Scan(&categories).Error

	if errCheckCat != nil {
		helper.Logger("error", "In Server: "+errCheckCat.Error())
		return nil, errors.New(errCheckCat.Error())
	}

	isJobCategoryExist := len(categories)

	if isJobCategoryExist == 0 {
		return nil, errors.New("JOB_NOT_FOUND")
	}

	checkQueryPlace := `SELECT id, name FROM places WHERE id = ?`

	errCheckPlace := db.Debug().Raw(checkQueryPlace, j.PlaceId).Scan(&places).Error

	if errCheckPlace != nil {
		helper.Logger("error", "In Server: "+errCheckPlace.Error())
		return nil, errors.New(errCheckPlace.Error())
	}

	isPlaceExist := len(places)

	if isPlaceExist == 0 {
		return nil, errors.New("PLACE_NOT_FOUND")
	}

	query := `INSERT INTO jobs (uid, title, caption, salary, min_salary, max_salary, worker_count, cat_id, company_id, place_id, user_id, is_draft)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	err := db.Debug().Exec(query, j.Id, j.Title, j.Caption, j.Salary, j.MinSalary, j.MaxSalary, j.WorkerCount, j.CatId, j.CompanyId, j.PlaceId, j.UserId, j.IsDraft).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	for _, skill := range j.Skills {
		var JobCatId = uuid.NewV4().String()

		queryInsertJobSkillCategory := `INSERT INTO job_skill_categories (uid, name) VALUES (?, ?)`

		errInsertJobSkillCategory := db.Debug().Exec(queryInsertJobSkillCategory, JobCatId, skill).Error

		if errInsertJobSkillCategory != nil {
			helper.Logger("error", "In Server: "+errInsertJobSkillCategory.Error())
			return nil, errors.New(errInsertJobSkillCategory.Error())
		}

		queryInsertJobSkills := `INSERT INTO job_skills (job_id, cat_id)
		VALUES (?, ?)`

		errInsertJobSkills := db.Debug().Exec(queryInsertJobSkills, j.Id, JobCatId).Error

		if errInsertJobSkills != nil {
			helper.Logger("error", "In Server: "+errInsertJobSkills.Error())
			return nil, errors.New(errInsertJobSkills.Error())
		}
	}

	return map[string]any{}, nil
}

func JobUpdate(j *models.JobUpdate) (map[string]any, error) {

	categories := []entities.JobCategory{}
	places := []entities.JobPlace{}

	checkQueryCat := `SELECT uid AS id FROM job_categories WHERE uid = '` + j.CatId + `'`

	errCheckCat := db.Debug().Raw(checkQueryCat).Scan(&categories).Error

	if errCheckCat != nil {
		helper.Logger("error", "In Server: "+errCheckCat.Error())
		return nil, errors.New(errCheckCat.Error())
	}

	isJobCategoryExist := len(categories)

	if isJobCategoryExist == 0 {
		return nil, errors.New("JOB_NOT_FOUND")
	}

	checkQueryPlace := `SELECT id, name FROM places WHERE id = ?`

	errCheckPlace := db.Debug().Raw(checkQueryPlace, j.PlaceId).Scan(&places).Error

	if errCheckPlace != nil {
		helper.Logger("error", "In Server: "+errCheckPlace.Error())
		return nil, errors.New(errCheckPlace.Error())
	}

	isPlaceExist := len(places)

	if isPlaceExist == 0 {
		return nil, errors.New("PLACE_NOT_FOUND")
	}

	query := `UPDATE jobs SET title = ?, caption = ?, salary = ?, min_salary = ?, max_salary = ?, worker_count = ?, company_id = ?, place_id = ?, cat_id = ?, place_id = ?, is_draft = ?
	WHERE uid = ?`

	err := db.Debug().Exec(query, j.Title, j.Caption, j.Salary, j.MinSalary, j.MaxSalary, j.WorkerCount, j.CompanyId, j.PlaceId, j.CatId, j.PlaceId, j.IsDraft, j.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	for _, skill := range j.Skills {
		var JobCatId = uuid.NewV4().String()

		queryInsertJobSkillCategory := `INSERT INTO job_skill_categories (uid, name) VALUES (?, ?)`

		errInsertJobSkillCategory := db.Debug().Exec(queryInsertJobSkillCategory, JobCatId, skill).Error

		if errInsertJobSkillCategory != nil {
			helper.Logger("error", "In Server: "+errInsertJobSkillCategory.Error())
			return nil, errors.New(errInsertJobSkillCategory.Error())
		}

		queryInsertJobSkills := `INSERT INTO job_skills (job_id, cat_id)
		VALUES (?, ?)`

		errInsertJobSkills := db.Debug().Exec(queryInsertJobSkills, j.Id, JobCatId).Error

		if errInsertJobSkills != nil {
			helper.Logger("error", "In Server: "+errInsertJobSkills.Error())
			return nil, errors.New(errInsertJobSkills.Error())
		}
	}

	return map[string]any{}, nil
}

func JobSkillCategoryList() (map[string]any, error) {
	jobSkillCategoryList := []entities.JobSkillCategoryList{}

	queryJobSkillCategory := `SELECT uid AS id, name FROM job_skill_categories`

	errCheckCat := db.Debug().Raw(queryJobSkillCategory).Scan(&jobSkillCategoryList).Error

	if errCheckCat != nil {
		helper.Logger("error", "In Server: "+errCheckCat.Error())
		return nil, errors.New(errCheckCat.Error())
	}

	return map[string]any{
		"data": jobSkillCategoryList,
	}, nil
}

func JobSkillCategoryDelete(jscd *entities.JobSkillCategoryDelete) (map[string]any, error) {

	for _, skill := range jscd.Skills {
		queryDeleteJobSkillCategory := `DELETE FROM job_skill_categories WHERE uid = ?`

		errDeleteJobSkillCategory := db.Debug().Exec(queryDeleteJobSkillCategory, skill).Error

		if errDeleteJobSkillCategory != nil {
			helper.Logger("error", "In Server: "+errDeleteJobSkillCategory.Error())
			return nil, errors.New(errDeleteJobSkillCategory.Error())
		}

		queryDeleteJobSkill := `DELETE FROM job_skills WHERE job_id = ? AND cat_id = ?`

		errDeleteJobSkills := db.Debug().Exec(queryDeleteJobSkill, jscd.JobId, skill).Error

		if errDeleteJobSkills != nil {
			helper.Logger("error", "In Server: "+errDeleteJobSkills.Error())
			return nil, errors.New(errDeleteJobSkills.Error())
		}
	}

	return map[string]any{}, nil
}

func JobSkillCategoryStore(jscs *entities.JobSkillCategoryStore) (map[string]any, error) {
	query := `INSERT INTO job_skills (job_id, cat_id) 
	VALUES (?, ?)`

	err := db.Debug().Exec(query, jscs.JobId, jscs.CatId).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func JobFavourite(j *models.JobFavourite) (map[string]any, error) {
	jobFavourite := []entities.JobFavourite{}
	jobs := []entities.Job{}
	users := []entities.User{}

	checkQueryJob := `SELECT uid AS id FROM jobs WHERE uid = '` + j.JobId + `'`

	errCheckJob := db.Debug().Raw(checkQueryJob).Scan(&jobs).Error

	if errCheckJob != nil {
		helper.Logger("error", "In Server: "+errCheckJob.Error())
		return nil, errors.New(errCheckJob.Error())
	}

	isJobExist := len(jobs)

	if isJobExist == 0 {
		return nil, errors.New("job not found")
	}

	checkQueryUser := `SELECT uid AS id FROM users WHERE uid = '` + j.UserId + `'`

	errCheckUser := db.Debug().Raw(checkQueryUser).Scan(&users).Error

	if errCheckUser != nil {
		helper.Logger("error", "In Server: "+errCheckUser.Error())
		return nil, errors.New(errCheckUser.Error())
	}

	isUserExist := len(users)

	if isUserExist == 0 {
		return nil, errors.New("user not found")
	}

	checkQueryJobFavourite := `SELECT user_id, job_id FROM job_favourites WHERE user_id = '` + j.UserId + `' AND job_id = '` + j.JobId + `'`

	errCheckJobFavourite := db.Debug().Raw(checkQueryJobFavourite).Scan(&jobFavourite).Error

	if errCheckJobFavourite != nil {
		helper.Logger("error", "In Server: "+errCheckJobFavourite.Error())
		return nil, errors.New(errCheckJobFavourite.Error())
	}

	isJobFavouriteExist := len(jobFavourite)

	if isJobFavouriteExist == 0 {
		query := `INSERT INTO job_favourites (user_id, job_id) VALUES (?, ?)`

		err := db.Debug().Exec(query, j.UserId, j.JobId).Error

		if err != nil {
			helper.Logger("error", "In Server: "+err.Error())
			return nil, errors.New(err.Error())
		}
	} else {
		query := `DELETE FROM job_favourites WHERE user_id = ? AND job_id = ?`

		err := db.Debug().Exec(query, j.UserId, j.JobId).Error

		if err != nil {
			helper.Logger("error", "In Server: "+err.Error())
			return nil, errors.New(err.Error())
		}
	}

	return map[string]any{}, nil
}

func JobCategoryCount() (map[string]any, error) {
	categoryCounts := []entities.JobCategoryCount{}

	query := `SELECT jc.name, jc.logo AS icon, COUNT(j.cat_id) AS total
	FROM job_categories jc
	LEFT JOIN jobs j ON jc.uid = j.cat_id
	GROUP BY jc.name
	ORDER BY total DESC`

	err := db.Debug().Raw(query).Scan(&categoryCounts).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": categoryCounts,
	}, nil
}

func JobCategoryStore(j *models.JobCategoryStore) (map[string]any, error) {

	j.Id = uuid.NewV4().String()

	query := `INSERT INTO job_categories (uid, logo, name, type) VALUES (?, ?, ?, ?)`

	err := db.Debug().Exec(query, j.Id, j.Icon, j.Name, j.Type).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func CandidatePassesForm(dp *entities.DepartureForm) (map[string]any, error) {
	var dataUserFcm entities.InitFcm

	// Insert Departure
	queryDepartures := `INSERT INTO departures (content) VALUES (?)`

	resultDepartures, _ := db.DB().Exec(queryDepartures, dp.Content)

	lastID, errDepartures := resultDepartures.LastInsertId()
	if errDepartures != nil {
		helper.Logger("error", "In Server: "+errDepartures.Error())
		return nil, errors.New(errDepartures.Error())
	}

	// Insert Candidate Passes
	queryCandidatePasses := `INSERT INTO candidate_passes (departure_id, apply_job_id, user_candidate_id) VALUES (?, ?, ?)`

	errCandidatePasses := db.Debug().Exec(queryCandidatePasses, lastID, dp.ApplyJobId, dp.UserCandidateId).Error

	if errCandidatePasses != nil {
		helper.Logger("error", "In Server: "+errCandidatePasses.Error())
		return nil, errors.New(errCandidatePasses.Error())
	}

	// Insert Inbox
	queryInbox := `INSERT INTO inboxes (uid, field1, field2, user_id, type) VALUES (?, ?, ?, ?, ?)`

	errInbox := db.Debug().Exec(queryInbox, uuid.NewV4().String(), dp.Content, dp.ApplyJobId, dp.UserCandidateId, "departure").Error

	if errInbox != nil {
		helper.Logger("error", "In Server: "+errInbox.Error())
		return nil, errors.New(errInbox.Error())
	}

	// Update Apply Jobs
	queryUpdateApplyJobs := `UPDATE apply_jobs SET is_finish = ? WHERE uid = ?`
	errUpdateApplyJobs := db.Debug().Exec(queryUpdateApplyJobs, 1, dp.ApplyJobId).Error
	if errUpdateApplyJobs != nil {
		helper.Logger("error", "In Server: "+errUpdateApplyJobs.Error())
		return nil, errors.New(errUpdateApplyJobs.Error())
	}

	// Fcm
	queryUserFcm := `SELECT u.email, f.token, p.fullname FROM fcms f
	INNER JOIN profiles p ON p.user_id = f.user_id
	WHERE f.user_id = ?`

	rowUserFcm := db.Debug().Raw(queryUserFcm, dp.UserCandidateId).Row()

	errUserFcmRow := rowUserFcm.Scan(&dataUserFcm.Token, &dataUserFcm.Fullname)

	if errUserFcmRow != nil {
		if errors.Is(errUserFcmRow, sql.ErrNoRows) {
			helper.Logger("info", "No FCM data found for user")
		}

		helper.Logger("error", "In Server: "+errUserFcmRow.Error())
	}

	helper.SendEmail(dataUserFcm.Email, "TJL", "Jadwal Keberangkatan", dp.Content, "tjl-invitation")

	helper.SendFcm("Jadwal Keberangkatan", dataUserFcm.Fullname, dataUserFcm.Token, "notifications", "-")

	return map[string]any{
		"content":           dp.Content,
		"departure_id":      lastID,
		"apply_job_id":      dp.ApplyJobId,
		"user_candidate_id": dp.UserCandidateId,
	}, nil
}

func JobCategoryUpdate(j *models.JobCategoryUpdate) (map[string]any, error) {

	query := `UPDATE job_categories SET name = ?, logo = ?, type = ? WHERE uid = ?`

	err := db.Debug().Exec(query, j.Name, j.Icon, j.Type, j.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func JobCategoryDelete(j *models.JobCategoryDelete) (map[string]any, error) {

	query := `DELETE FROM job_categories WHERE uid = ?`

	err := db.Debug().Exec(query, j.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func JobCategory() (map[string]any, error) {
	categories := []entities.JobCategory{}

	query := `SELECT jc.uid AS id, jc.logo AS icon, jc.name, tj.name AS type
	FROM job_categories jc
	INNER JOIN type_jobs tj ON tj.id = jc.type 
	`

	err := db.Debug().Raw(query).Scan(&categories).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	isCategoryExist := len(categories)

	if isCategoryExist == 0 {
		return nil, errors.New("job cateogry not found")
	}

	return map[string]any{
		"data": categories,
	}, nil
}

func JobDelete(j *models.Job) (map[string]any, error) {
	errDeleteJob := db.Debug().Exec(`
	DELETE FROM jobs WHERE uid = ?`, j.Id).Error

	if errDeleteJob != nil {
		helper.Logger("error", "In Server: "+errDeleteJob.Error())
		return nil, errors.New(errDeleteJob.Error())
	}

	return map[string]any{}, nil
}

func TypeJobList() (map[string]any, error) {
	var typeJob []entities.TypeJob

	query := `SELECT id, name FROM type_jobs`

	err := db.Debug().Raw(query).Scan(&typeJob).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": typeJob,
	}, nil
}

func TypeJobStore(tjs *entities.TypeJobStore) (map[string]any, error) {
	query := `INSERT INTO type_jobs (name) VALUES (?)`

	err := db.Debug().Exec(query, tjs.Name).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func TypeJobUpdate(tju *entities.TypeJobUpdate) (map[string]any, error) {
	query := `UPDATE type_jobs SET name = ? WHERE id = ?`

	err := db.Debug().Exec(query, tju.Name, tju.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func TypeJobDelete(tjd *entities.TypeJobDelete) (map[string]any, error) {
	query := `DELETE FROM type_jobs WHERE id = ?`

	err := db.Debug().Exec(query, tjd.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}
