package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
)

func Summary(branchId string) (map[string]any, error) {

	var dataChartSummary entities.ChartSummaryResponse

	var genderQuery entities.GenderQuery
	var genderResponse = make([]entities.GenderResponse, 0)
	var countryQuery entities.CountryQuery
	var countryResponse = make([]entities.CountryResponse, 0)

	var applicantPerMonthdata entities.ApplicantPerMonthQuery

	var applicantPerMonthResponse = make([]entities.ApplicantPerMonthResponse, 0)

	// APPLICANTS PER MONTH
	queryApplicantsPerMonth := `
		SELECT DATE_FORMAT(aj.created_at, '%Y-%m') AS month, COUNT(*) AS total, b.name AS branch
		FROM apply_jobs aj
		INNER JOIN user_branches ub ON ub.user_id = aj.user_id
		INNER JOIN branchs b ON b.id = ub.branch_id
		WHERE b.id = ?
		GROUP BY month 
		ORDER BY month
	`

	rowsApplicantsPerMonth, err := db.Debug().Raw(queryApplicantsPerMonth, branchId).Rows()
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}
	defer rowsApplicantsPerMonth.Close()

	for rowsApplicantsPerMonth.Next() {
		errScanRows := db.ScanRows(rowsApplicantsPerMonth, &applicantPerMonthdata)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		applicantPerMonthResponse = append(applicantPerMonthResponse, entities.ApplicantPerMonthResponse{
			Month:  applicantPerMonthdata.Month,
			Branch: applicantPerMonthdata.Branch,
			Total:  applicantPerMonthdata.Total,
		})

	}

	// GENDER
	queryGender := `
		SELECT fb.gender, COUNT(*) AS total
		FROM form_biodatas fb
		INNER JOIN user_branches ub ON ub.user_id = fb.user_id
		INNER JOIN branchs b ON b.id = ub.branch_id 
		WHERE b.id = ?
		GROUP BY fb.gender
	`

	rowsGender, err := db.Debug().Raw(queryGender, branchId).Rows()
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}
	defer rowsGender.Close()

	for rowsGender.Next() {
		errScanRows := db.ScanRows(rowsGender, &genderQuery)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		genderResponse = append(genderResponse, entities.GenderResponse{
			Gender: genderQuery.Gender,
			Total:  genderQuery.Total,
		})
	}

	// COUNTRY
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
		errScanRows := db.ScanRows(rowsCountry, &countryQuery)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		countryResponse = append(countryResponse, entities.CountryResponse{
			Country: countryQuery.Country,
			Total:   countryQuery.Total,
		})
	}

	dataChartSummary.ApplicantsPerMonth = applicantPerMonthResponse
	dataChartSummary.Countries = countryResponse
	dataChartSummary.Genders = genderResponse

	return map[string]any{
		"data": dataChartSummary,
	}, nil
}

func AdminListUser() (map[string]any, error) {
	var adminListUserData []entities.AdminListUserResponse

	queryUsers := `SELECT p.user_id AS id, u.email, u.phone, p.avatar, p.fullname, ur.name AS role, u.created_at
		FROM users u 
		INNER JOIN profiles p ON p.user_id = u.uid
		INNER JOIN user_roles ur ON ur.id = u.role
	`

	rows, err := db.Debug().Raw(queryUsers).Rows()
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

		var branch entities.AdminListUserBranch
		if len(adminListUserBranchData) > 0 {
			branch = adminListUserBranchData[0]
		} else {
			branch = entities.AdminListUserBranch{
				Id:   0,
				Name: "-",
			}
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
