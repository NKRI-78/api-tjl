package services

import (
	"database/sql"
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
)

func AdminCandidatePassesBadges(branchId string) (map[string]any, error) {
	var badgeCount int

	baseQuery := `
	SELECT COUNT(*) AS badge_count
	FROM (
		SELECT aj.uid
		FROM apply_jobs aj
		LEFT JOIN candidate_passes cp ON cp.apply_job_id = aj.uid
		INNER JOIN user_branches ub ON ub.user_id = aj.user_id
		INNER JOIN branchs b ON b.id = ub.branch_id
		WHERE aj.status = ?
		AND cp.apply_job_id IS NULL
		AND EXISTS (
			SELECT 1 FROM apply_job_documents ajd
			WHERE ajd.apply_job_id = aj.uid
		)
	`

	var args []any
	args = append(args, "3")

	// Conditionally add branch filter
	if branchId != "" {
		baseQuery += " AND b.id = ?"
		args = append(args, branchId)
	}

	// Close the subquery
	baseQuery += ") AS filtered_jobs"

	err := db.Raw(baseQuery, args...).Row().Scan(&badgeCount)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	return map[string]any{
		"data": badgeCount,
	}, nil
}

func AdminApplyJobBadges(branchId string) (map[string]any, error) {
	var dataAdminApplyJobBadges entities.AdminApplyJobBadges

	baseQuery := `
		SELECT COUNT(*) AS total 
		FROM users u 
		INNER JOIN apply_jobs aj ON aj.user_id = u.uid
		INNER JOIN user_branches ub ON ub.user_id = u.uid
		INNER JOIN branchs b ON b.id = ub.branch_id
		WHERE aj.status = 1
	`

	var args []any

	// Conditionally add branch filter
	if branchId != "" {
		baseQuery += " AND b.id = ?"
		args = append(args, branchId)
	}

	row := db.Debug().Raw(baseQuery, args...).Row()
	err := row.Scan(&dataAdminApplyJobBadges.Total)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	return map[string]any{
		"data": dataAdminApplyJobBadges.Total,
	}, nil
}

func Summary(branchId string) (map[string]any, error) {
	var dataChartSummary entities.ChartSummaryResponse

	var genderQuery entities.GenderQuery
	var genderResponse = make([]entities.GenderResponse, 0)

	var countryQuery entities.CountryQuery
	var countryResponse = make([]entities.CountryResponse, 0)

	var applicantPerMonthData entities.ApplicantPerMonthQuery
	var applicantPerBranchData entities.ApplicantPerMonthQuery
	var applicantPerMonthResponse = make([]entities.ApplicantPerMonthResponse, 0)
	var applicantPerBranchResponse = make([]entities.ApplicantPerMonthResponse, 0)

	// === APPLICANTS PER MONTH ===
	queryApplicantsPerMonth := `
		SELECT DATE_FORMAT(aj.created_at, '%Y-%m') AS month, COUNT(*) AS total, b.name AS branch
		FROM apply_jobs aj
		INNER JOIN user_branches ub ON ub.user_id = aj.user_id
		INNER JOIN branchs b ON b.id = ub.branch_id
	`
	var args []any
	if branchId != "" {
		queryApplicantsPerMonth += " WHERE b.id = ?"
		args = append(args, branchId)
	}
	queryApplicantsPerMonth += " GROUP BY month ORDER BY month"

	rowsApplicantsPerMonth, err := db.Debug().Raw(queryApplicantsPerMonth, args...).Rows()
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}
	defer rowsApplicantsPerMonth.Close()

	for rowsApplicantsPerMonth.Next() {
		if err := db.ScanRows(rowsApplicantsPerMonth, &applicantPerMonthData); err != nil {
			helper.Logger("error", "In Server: "+err.Error())
			return nil, err
		}
		applicantPerMonthResponse = append(applicantPerMonthResponse, entities.ApplicantPerMonthResponse{
			Month:  applicantPerMonthData.Month,
			Branch: applicantPerMonthData.Branch,
			Total:  applicantPerMonthData.Total,
		})
	}

	// === APPLICANTS PER BRANCH ===
	queryApplicantsPerBranch := `
		SELECT DATE_FORMAT(aj.created_at, '%Y-%m') AS month, COUNT(*) AS total, b.name AS branch
		FROM apply_jobs aj
		INNER JOIN user_branches ub ON ub.user_id = aj.user_id
		INNER JOIN branchs b ON b.id = ub.branch_id
		GROUP BY b.id
	`

	rowsApplicantsPerBranch, err := db.Debug().Raw(queryApplicantsPerBranch).Rows()
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}
	defer rowsApplicantsPerBranch.Close()

	for rowsApplicantsPerBranch.Next() {
		if err := db.ScanRows(rowsApplicantsPerBranch, &applicantPerBranchData); err != nil {
			helper.Logger("error", "In Server: "+err.Error())
			return nil, err
		}
		applicantPerBranchResponse = append(applicantPerBranchResponse, entities.ApplicantPerMonthResponse{
			Month:  applicantPerBranchData.Month,
			Branch: applicantPerBranchData.Branch,
			Total:  applicantPerBranchData.Total,
		})
	}

	// === GENDER ===
	queryGender := `
		SELECT fb.gender, COUNT(*) AS total
		FROM form_biodatas fb
		INNER JOIN user_branches ub ON ub.user_id = fb.user_id
		INNER JOIN apply_jobs aj ON aj.user_id  = ub.user_id
		INNER JOIN branchs b ON b.id = ub.branch_id
	`
	args = []any{}
	if branchId != "" {
		queryGender += " WHERE b.id = ?"
		args = append(args, branchId)
	}
	queryGender += " GROUP BY fb.gender"

	rowsGender, err := db.Debug().Raw(queryGender, args...).Rows()
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}
	defer rowsGender.Close()

	for rowsGender.Next() {
		if err := db.ScanRows(rowsGender, &genderQuery); err != nil {
			helper.Logger("error", "In Server: "+err.Error())
			return nil, err
		}
		genderResponse = append(genderResponse, entities.GenderResponse{
			Gender: genderQuery.Gender,
			Total:  genderQuery.Total,
		})
	}

	// === COUNTRY === (no filtering by branchId)
	queryCountry := `
		SELECT p.name AS country, COUNT(*) AS total
		FROM jobs j
		INNER JOIN places p ON p.id = j.place_id 
		GROUP BY p.id
	`

	rowsCountry, err := db.Debug().Raw(queryCountry).Rows()
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}
	defer rowsCountry.Close()

	for rowsCountry.Next() {
		if err := db.ScanRows(rowsCountry, &countryQuery); err != nil {
			helper.Logger("error", "In Server: "+err.Error())
			return nil, err
		}
		countryResponse = append(countryResponse, entities.CountryResponse{
			Country: countryQuery.Country,
			Total:   countryQuery.Total,
		})
	}

	dataChartSummary.ApplicantsPerMonth = applicantPerMonthResponse
	dataChartSummary.ApplicantsPerBranch = applicantPerBranchResponse
	dataChartSummary.Countries = countryResponse
	dataChartSummary.Genders = genderResponse

	return map[string]any{
		"data": dataChartSummary,
	}, nil
}

func AdminListUser(Type, branchId string) (map[string]any, error) {
	var adminListUserData []entities.AdminListUserResponse

	queryUsers := `SELECT p.user_id AS id, u.email, u.phone, p.avatar, p.fullname, ur.name AS role, u.created_at
		FROM users u 
		INNER JOIN profiles p ON p.user_id = u.uid
		INNER JOIN user_roles ur ON ur.id = u.role
		INNER JOIN user_branches ub ON ub.user_id = u.uid
		INNER JOIN branchs b ON b.id = ub.branch_id`

	var args []any
	whereClause := ""

	if branchId != "" {
		whereClause += " WHERE b.id = ?"
		args = append(args, branchId)
	}

	if Type != "" && Type != "manual" {
		if whereClause == "" {
			whereClause += " WHERE u.via = 'auto'"
		} else {
			whereClause += " AND u.via = 'auto'"
		}
	}

	queryUsers += whereClause

	rows, err := db.Debug().Raw(queryUsers, args...).Rows()
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var adminListUserQuery entities.AdminListUser

		if err := rows.Scan(
			&adminListUserQuery.Id,
			&adminListUserQuery.Email,
			&adminListUserQuery.Phone,
			&adminListUserQuery.Avatar,
			&adminListUserQuery.Fullname,
			&adminListUserQuery.Role,
			&adminListUserQuery.CreatedAt,
		); err != nil {
			helper.Logger("error", "In Server: "+err.Error())
			return nil, errors.New(err.Error())
		}

		var adminListUserBranchData []entities.AdminListUserBranch
		queryBranches := `SELECT b.id, b.name
			FROM user_branches ub 
			LEFT JOIN branchs b ON b.id = ub.branch_id 
			WHERE ub.user_id = ?`
		if errUserBranch := db.Debug().Raw(queryBranches, adminListUserQuery.Id).Scan(&adminListUserBranchData).Error; errUserBranch != nil {
			helper.Logger("error", "In Server: "+errUserBranch.Error())
			return nil, errors.New(errUserBranch.Error())
		}

		branch := entities.AdminListUserBranch{
			Id:   0,
			Name: "-",
		}
		if len(adminListUserBranchData) > 0 {
			branch = adminListUserBranchData[0]
		}

		adminListUserData = append(adminListUserData, entities.AdminListUserResponse{
			Id:        adminListUserQuery.Id,
			Avatar:    adminListUserQuery.Avatar.String,
			Fullname:  adminListUserQuery.Fullname,
			Email:     adminListUserQuery.Email,
			Phone:     adminListUserQuery.Phone,
			Role:      adminListUserQuery.Role,
			Branch:    branch,
			CreatedAt: adminListUserQuery.CreatedAt,
		})
	}

	return map[string]any{
		"data": adminListUserData,
	}, nil
}

func ViewPdfDeparture(userId, applyJobId string) (map[string]any, error) {
	query := `SELECT field1 AS content FROM inboxes 
	WHERE user_id = ? AND field2 = ?`

	var content string

	row := db.Debug().Raw(query, userId, applyJobId).Row()
	err := row.Scan(&content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return map[string]any{
				"data": nil,
			}, nil
		}
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	updateQuery := `UPDATE inboxes SET is_read = true WHERE field2 = ? AND user_id = ?`
	if err := db.Debug().Exec(updateQuery, applyJobId, userId).Error; err != nil {
		helper.Logger("error", "Failed to update is_read: "+err.Error())
	}

	return map[string]any{
		"data": content,
	}, nil
}

func ViewPdfApplyJobOffline(applyJobId string) (map[string]any, error) {
	query := `SELECT content FROM apply_job_offlines 
	WHERE apply_job_id = ?`

	var content string

	row := db.Debug().Raw(query, applyJobId).Row()
	err := row.Scan(&content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return map[string]any{
				"data": nil,
			}, nil
		}
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	return map[string]any{
		"data": content,
	}, nil
}

func ImportUserCreate(ius *entities.ImportUserStore) (map[string]any, error) {
	// INSERT USER
	queryInsertUser := `INSERT INTO users (uid, email, phone, password, enabled, via) 
	VALUES (?, ?, ?, ?, 1)`

	errInsertUser := db.Debug().Exec(queryInsertUser, ius.UserId, ius.Email, ius.Phone, ius.Password, "auto").Error

	if errInsertUser != nil {
		helper.Logger("error", "In Server: "+errInsertUser.Error())
		return nil, errors.New(errInsertUser.Error())
	}

	// INSERT PROFILE
	queryInsertProfile := `INSERT INTO profiles(fullname, user_id) 
	VALUES (?, ?)`

	errInsertProfile := db.Debug().Exec(queryInsertProfile, ius.Fullname, ius.UserId).Error

	if errInsertProfile != nil {
		helper.Logger("error", "In Server: "+errInsertProfile.Error())
		return nil, errors.New(errInsertProfile.Error())
	}

	// INSERT FORM BIODATA
	queryInsertBiodata := `INSERT INTO form_biodatas (birthdate, gender, weight, height, status, religion, place, user_id) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	errInsertFormBiodata := db.Debug().Exec(queryInsertBiodata, ius.Birthdate, ius.Gender, ius.Weight, ius.Height, ius.MaritalStatus, ius.Religion, ius.Place, ius.UserId).Error

	if errInsertFormBiodata != nil {
		helper.Logger("error", "In Server: "+errInsertFormBiodata.Error())
		return nil, errors.New(errInsertFormBiodata.Error())
	}

	// INSERT FORM EDUCATION
	queryInsertFormEducation := `INSERT INTO form_educations (education_level, major, school_or_college, start_month, start_year, end_month, end_year, user_id) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	errInsertForumEducation := db.Debug().Exec(queryInsertFormEducation, ius.EducationLevel, ius.Major, ius.SchoolOrCollege, ius.StartMonthEducation, ius.StartYearEducation, ius.EndMonthEducation, ius.EndYearEducation, ius.UserId).Error

	if errInsertForumEducation != nil {
		return nil, errors.New(errInsertForumEducation.Error())
	}

	// INSERT FORM EXERCISE
	queryInsertFormExercise := `INSERT INTO form_exercises (name, institution, start_month, start_year, end_month, end_year, user_id) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	errInsertFormExercise := db.Debug().Exec(queryInsertFormExercise, ius.NameInstitution, ius.Institution, ius.StartMonthExercise, ius.StartYearExercise, ius.EndMonthExercise, ius.EndYearExercise, ius.UserId).Error

	if errInsertFormExercise != nil {
		return nil, errors.New(errInsertFormExercise.Error())
	}

	// INSERT FORM WORK
	queryInsertFormWork := `INSERT INTO form_works (position, institution, work, country, city, start_month, start_year, end_month, end_year, is_work, user_id)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	errInsertFormWork := db.Debug().Exec(queryInsertFormWork, ius.PositionWork, ius.InstitutionWork, ius.Work, ius.CountryWork, ius.CityWork, ius.StartMonthWork, ius.StartYearWork, ius.EndMonthWork, ius.EndYearWork, 0, ius.UserId).Error

	if errInsertFormWork != nil {
		return nil, errors.New(errInsertFormWork.Error())
	}

	// INSERT FORM LANGUAGE
	queryInsertFormLanguage := `INSERT INTO form_languages (level, language, user_id) 
	VALUES (?, ?, ?)`

	errInsertFormLanguage := db.Debug().Exec(queryInsertFormLanguage, ius.Level, ius.Language, ius.UserId).Error

	if errInsertFormLanguage != nil {
		return nil, errors.New(errInsertFormLanguage.Error())
	}

	return map[string]any{}, nil

}

// func ViewOfflinePdfDeparture(userId, applyJobId string) (map[string]any, error) {
// 	query := `SELECT d.content FROM departures d
// 	INNER JOIN candidate_passes cp
// 	ON cp.departure_id = d.id
// 	WHERE cp.user_candidate_id = ? AND cp.apply_job_id = ?`

// 	var content string

// 	row := db.Debug().Raw(query, userId, applyJobId).Row()
// 	err := row.Scan(&content)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return map[string]any{
// 				"data": nil,
// 			}, nil
// 		}
// 		helper.Logger("error", "In Server: "+err.Error())
// 		return nil, err
// 	}

// 	updateQuery := `UPDATE inboxes SET is_read = true WHERE field2 = ? AND user_id = ?`
// 	if err := db.Debug().Exec(updateQuery, applyJobId, userId).Error; err != nil {
// 		helper.Logger("error", "Failed to update is_read: "+err.Error())
// 	}

// 	return map[string]any{
// 		"data": content,
// 	}, nil
// }
