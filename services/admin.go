package services

import (
	"database/sql"
	"errors"
	"fmt"
	entities "superapps/entities"
	helper "superapps/helpers"

	"github.com/jinzhu/gorm"
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
	statusSequence := []string{
		"IN_PROGRESS", "INTERVIEW", "ACCEPTED", "PROGRESS_MCU",
		"PROGRESS_LINK_SIAPKERJA", "PROGRESS_SKCK", "PROGRESS_VISA",
		"PROGRESS_TTD", "PROGRESS_OPP", "PROGRESS_DONE", "DECLINED",
	}

	var badgeCounts []struct {
		Status string
		Total  int
	}

	baseQuery := `
		SELECT js.name AS status, COUNT(*) AS total
		FROM users u
		INNER JOIN apply_jobs aj ON aj.user_id = u.uid
		INNER JOIN job_statuses js ON js.id = aj.status
		INNER JOIN user_branches ub ON ub.user_id = u.uid
		INNER JOIN branchs b ON b.id = ub.branch_id
	`

	var args []any
	if branchId != "" {
		baseQuery += " WHERE b.id = ?"
		args = append(args, branchId)
	}
	baseQuery += " GROUP BY js.name"

	err := db.Debug().Raw(baseQuery, args...).Scan(&badgeCounts).Error
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	countMap := make(map[string]int)
	for _, row := range badgeCounts {
		countMap[row.Status] = row.Total
	}

	var orderedResult []map[string]any
	for _, status := range statusSequence {
		orderedResult = append(orderedResult, map[string]any{
			"status": status,
			"total":  countMap[status],
		})
	}

	return map[string]any{
		"data": orderedResult,
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
		INNER JOIN branchs b ON b.id = ub.branch_id 
		WHERE u.role != '4'`

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

	if adminListUserData == nil {
		adminListUserData = []entities.AdminListUserResponse{}
	}

	return map[string]any{
		"data": adminListUserData,
	}, nil
}

func AdminListCandidate() (map[string]any, error) {
	return map[string]any{}, nil
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

func AdminListCandidateImportV2() ([]entities.ProfileImportResponse, error) {
	profiles := []entities.ProfileImport{}

	query := `SELECT p.user_id AS id, p.fullname, p.avatar, u.phone, u.email, u.enabled, 
		ur.name AS role,
		c.uid AS company_id,
		c.name AS company,
		jc.name AS position,
		j.salary,
		fb.id AS bio_id,
		fb.birthdate AS bio_birthdate,
		fb.gender AS bio_gender,
		fb.weight AS bio_weight,
		fb.height AS bio_height,
		fb.religion AS bio_religion,
		fb.status AS bio_status,
		fb.place AS bio_place,
		fp.detail_address AS bio_detail_address,
		fp.id AS bio_address_id,
		pro.id AS bio_province_id,
		pro.name AS bio_province,
		reg.id AS bio_city_id,
		reg.name AS bio_city,
		dis.id AS bio_district_id,
		dis.name AS bio_district, 
		vil.id AS bio_subdistrict_id,
		vil.name AS bio_subdistrict 
		FROM profiles p 
		INNER JOIN users u ON u.uid = p.user_id
		LEFT JOIN user_roles ur ON ur.id = u.role
		LEFT JOIN apply_jobs aj ON aj.user_id = p.user_id
		LEFT JOIN jobs j ON j.uid = aj.job_id
		LEFT JOIN job_categories jc ON jc.uid = j.cat_id
		LEFT JOIN companies c ON c.uid = j.company_id
		LEFT JOIN form_biodatas fb ON fb.user_id = p.user_id
		LEFT JOIN form_places fp ON fp.user_id = p.user_id
		LEFT JOIN provinces pro ON pro.id = fp.province_id
		LEFT JOIN regencies reg ON reg.id = fp.city_id
		LEFT JOIN districts dis ON dis.id = fp.district_id
		LEFT JOIN villages vil ON vil.id = fp.subdistrict_id
		WHERE u.via = ?
	`

	if err := db.Debug().Raw(query, "auto").Scan(&profiles).Error; err != nil {
		return nil, err
	}
	if len(profiles) == 0 {
		return nil, errors.New("no profiles found")
	}

	var result []entities.ProfileImportResponse

	for _, prof := range profiles {
		profileResp := entities.ProfileImportResponse{
			Id:        prof.Id,
			Fullname:  prof.Fullname,
			Avatar:    prof.Avatar,
			Email:     prof.Email,
			Role:      prof.Role,
			Phone:     prof.Phone,
			IsEnabled: prof.Enabled == 1,
			Job: entities.ProfileImportJobResponse{
				Id:       prof.CompanyId,
				Company:  prof.Company,
				Position: prof.Position,
				Salary:   prof.Salary,
			},
			Biodata: entities.Biodata{
				Personal: entities.ProfileFormBiodata{
					Id:        prof.BioId,
					Birthdate: prof.BioBirthdate,
					Gender:    prof.BioGender,
					Height:    prof.BioHeight,
					Weight:    prof.BioWeight,
					Religion:  prof.BioReligion,
					Place:     prof.BioPlace,
					Status:    prof.BioStatus,
				},
				Address: entities.ProfileFormPlace{
					Id:            prof.BioAddressId,
					DetailAddress: prof.BioDetailAddress,
					Province: entities.ProfileFormPlaceData{
						Id:   prof.BioProvinceId,
						Name: prof.BioProvince,
					},
					City: entities.ProfileCityPlaceData{
						Id:   prof.BioCityId,
						Name: prof.BioCity,
					},
					District: entities.ProfileDistrictPlaceData{
						Id:   prof.BioDistrictId,
						Name: prof.BioDistrict,
					},
					Subdistrict: entities.ProfileSubdistrictPlaceData{
						Id:   prof.BioSubdistrictId,
						Name: prof.BioSubdistrict,
					},
				},
			},
		}

		// Get Educations
		educations, err := getEducations(db, prof.Id)
		if err != nil {
			return nil, err
		}
		profileResp.Biodata.Educations = educations

		// Get Trainings
		trainings, err := getTrainings(db, prof.Id)
		if err != nil {
			return nil, err
		}
		profileResp.Biodata.Trainings = trainings

		// Get Experiences (Works)
		works, err := getWorks(db, prof.Id)
		if err != nil {
			return nil, err
		}
		profileResp.Biodata.Experiences = works

		// Get Languages
		languages, err := getLanguages(db, prof.Id)
		if err != nil {
			return nil, err
		}
		profileResp.Biodata.Languages = languages

		result = append(result, profileResp)
	}

	return result, nil
}

func logAndReturnError(message string, err error) error {
	helper.Logger("error", fmt.Sprintf("%s: %s", message, err.Error()))
	return errors.New(err.Error())
}

func getEducations(db *gorm.DB, userId string) ([]entities.ProfileFormEducation, error) {
	var dataEdu []entities.ProfileFormEducation

	queryEdu := `SELECT id, education_level, major, school_or_college, start_year, start_month, end_month, end_year, user_id 
	FROM form_educations WHERE user_id = ?`

	rows, err := db.Debug().Raw(queryEdu, userId).Rows()
	if err != nil {
		return nil, logAndReturnError("Education Query", err)
	}
	defer rows.Close()

	for rows.Next() {
		var edu entities.ProfileFormEducation
		if err := db.ScanRows(rows, &edu); err != nil {
			return nil, logAndReturnError("Education Scan", err)
		}

		var letters []entities.FormEducationLetter
		queryMedia := `SELECT id, path FROM form_exercise_medias WHERE exercise_id = ?`
		mediaRows, err := db.Debug().Raw(queryMedia, edu.Id).Rows()
		if err != nil {
			return nil, logAndReturnError("Education Media Query", err)
		}

		for mediaRows.Next() {
			var media entities.FormEducationLetter
			if err := db.ScanRows(mediaRows, &media); err != nil {
				return nil, logAndReturnError("Education Media Scan", err)
			}
			letters = append(letters, media)
		}
		mediaRows.Close()

		edu.Letters = letters
		dataEdu = append(dataEdu, edu)
	}

	return dataEdu, nil
}

func getTrainings(db *gorm.DB, userId string) ([]entities.ProfileFormExercise, error) {
	var dataTraining []entities.ProfileFormExercise

	query := `SELECT id, name, institution, start_month, start_year, end_month, end_year, user_id 
	FROM form_exercises WHERE user_id = ?`

	rows, err := db.Debug().Raw(query, userId).Rows()
	if err != nil {
		return nil, logAndReturnError("Training Query", err)
	}
	defer rows.Close()

	for rows.Next() {
		var training entities.ProfileFormExercise
		if err := db.ScanRows(rows, &training); err != nil {
			return nil, logAndReturnError("Training Scan", err)
		}

		var certs []entities.FormExerciseCertificate
		queryMedia := `SELECT id, path FROM form_exercise_medias WHERE exercise_id = ?`
		mediaRows, err := db.Debug().Raw(queryMedia, training.Id).Rows()
		if err != nil {
			return nil, logAndReturnError("Training Media Query", err)
		}

		for mediaRows.Next() {
			var cert entities.FormExerciseCertificate
			if err := db.ScanRows(mediaRows, &cert); err != nil {
				return nil, logAndReturnError("Training Media Scan", err)
			}
			certs = append(certs, cert)
		}
		mediaRows.Close()

		training.Certificates = certs
		dataTraining = append(dataTraining, training)
	}

	return dataTraining, nil
}

func getWorks(db *gorm.DB, userId string) ([]entities.ProfileFormWork, error) {
	var dataWork []entities.ProfileFormWork

	query := `SELECT id, work, position, institution, is_work, country, city, start_month, start_year, end_month, end_year, user_id 
	FROM form_works WHERE user_id = ?`

	rows, err := db.Debug().Raw(query, userId).Rows()
	if err != nil {
		return nil, logAndReturnError("Work Query", err)
	}
	defer rows.Close()

	for rows.Next() {
		var workRow struct {
			Id          int
			Work        string
			Position    string
			Institution string
			Country     string
			City        string
			IsWork      bool
			StartMonth  string
			StartYear   string
			EndMonth    string
			EndYear     string
		}
		if err := db.ScanRows(rows, &workRow); err != nil {
			return nil, logAndReturnError("Work Scan", err)
		}

		work := entities.ProfileFormWork{
			Id:          workRow.Id,
			Work:        workRow.Work,
			Position:    workRow.Position,
			Institution: workRow.Institution,
			IsWork:      workRow.IsWork,
			City:        workRow.City,
			Country:     workRow.Country,
			StartMonth:  workRow.StartMonth,
			StartYear:   workRow.StartYear,
			EndMonth:    workRow.EndMonth,
			EndYear:     workRow.EndYear,
		}
		dataWork = append(dataWork, work)
	}

	return dataWork, nil
}

func getLanguages(db *gorm.DB, userId string) ([]entities.ProfileFormLanguage, error) {
	var dataLanguage []entities.ProfileFormLanguage

	query := `SELECT id, language, level FROM form_languages WHERE user_id = ?`

	rows, err := db.Debug().Raw(query, userId).Rows()
	if err != nil {
		return nil, logAndReturnError("Language Query", err)
	}
	defer rows.Close()

	for rows.Next() {
		var lang entities.ProfileFormLanguage
		if err := db.ScanRows(rows, &lang); err != nil {
			return nil, logAndReturnError("Language Scan", err)
		}
		dataLanguage = append(dataLanguage, lang)
	}

	return dataLanguage, nil
}

func AdminListCandidateImport() (map[string]any, error) {

	profiles := []entities.Profile{}
	education := entities.ProfileFormEducation{}
	exercise := entities.ProfileFormExercise{}
	educationMedia := entities.FormExerciseCertificate{}
	exerciseMedia := entities.FormExerciseCertificate{}
	work := entities.ProfileFormWorkQuery{}
	language := entities.ProfileFormLanguage{}

	query := `SELECT p.user_id AS id, p.fullname, p.avatar, u.phone, u.email, u.enabled, 
	fb.id AS bio_id,
	fb.birthdate AS bio_birthdate,
	fb.gender AS bio_gender,
	fb.weight AS bio_weight,
	fb.height AS bio_height,
	fb.religion AS bio_religion,
	fb.status AS bio_status,
	fb.place AS bio_place,
	fp.detail_address AS bio_detail_address,
	fp.id AS bio_address_id,
	pro.id AS bio_province_id,
	pro.name AS bio_province,
	reg.id AS bio_city_id,
	reg.name AS bio_city,
	dis.id AS bio_district_id,
	dis.name AS bio_district, 
	vil.id AS bio_subdistrict_id,
	vil.name AS bio_subdistrict 
	FROM profiles p 
	INNER JOIN users u ON u.uid = p.user_id
	LEFT JOIN form_biodatas fb ON fb.user_id = p.user_id
	LEFT JOIN form_places fp ON fp.user_id = p.user_id
	LEFT JOIN provinces pro ON pro.id = fp.province_id
	LEFT JOIN regencies reg ON reg.id = fp.city_id
	LEFT JOIN districts dis ON dis.id = fp.district_id
	LEFT JOIN villages vil ON vil.id = fp.subdistrict_id
	WHERE u.via = ?`

	err := db.Debug().Raw(query, "auto").Scan(&profiles).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	isProfileExist := len(profiles)

	if isProfileExist == 0 {
		return nil, errors.New("profile not found")
	}

	profile := entities.ProfileResponse{}

	var enabled bool

	if profiles[0].Enabled == 1 {
		enabled = true
	} else {
		enabled = false
	}

	// Education

	var dataEdu = make([]entities.ProfileFormEducation, 0)

	queryEdu := `SELECT id, education_level, major, school_or_college, start_year, start_month, end_month, end_year, user_id 
	FROM form_educations WHERE user_id  = '` + profiles[0].Id + `'`

	rows, errEdu := db.Debug().Raw(queryEdu).Scan(&education).Rows()

	if errEdu != nil {
		helper.Logger("error", "In Server: "+errEdu.Error())
		return nil, errors.New(errEdu.Error())
	}

	defer rows.Close()

	for rows.Next() {
		errScanRows := db.ScanRows(rows, &education)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		queryFormEducationMedia := `SELECT id, path
		FROM form_exercise_medias WHERE exercise_id  = ?`

		rowsFormEducationMedia, errFormExerciseMedia := db.Debug().Raw(queryFormEducationMedia, education.Id).Scan(&educationMedia).Rows()

		if errFormExerciseMedia != nil {
			helper.Logger("error", "In Server: "+errFormExerciseMedia.Error())
			return nil, errors.New(errFormExerciseMedia.Error())
		}

		var dataFormEducationLetter = make([]entities.FormEducationLetter, 0)

		for rowsFormEducationMedia.Next() {
			errScanFormEducationMedia := db.ScanRows(rowsFormEducationMedia, &exerciseMedia)

			if errScanFormEducationMedia != nil {
				helper.Logger("error", "In Server: "+errScanFormEducationMedia.Error())
				return nil, errors.New(errScanFormEducationMedia.Error())
			}

			dataFormEducationLetter = append(dataFormEducationLetter, entities.FormEducationLetter{
				Id:   educationMedia.Id,
				Path: educationMedia.Path,
			})
		}

		dataEdu = append(dataEdu, entities.ProfileFormEducation{
			Id:              education.Id,
			EducationLevel:  education.EducationLevel,
			Major:           education.Major,
			SchoolOrCollege: education.SchoolOrCollege,
			StartYear:       education.StartYear,
			EndYear:         education.EndYear,
			StartMonth:      education.StartMonth,
			EndMonth:        education.EndMonth,
			Letters:         dataFormEducationLetter,
		})
	}

	// Trainning

	var dataTraining = make([]entities.ProfileFormExercise, 0)

	queryTraining := `SELECT id, name, institution, start_month, start_year, end_month, end_year, user_id 
	FROM form_exercises WHERE user_id  = ?`

	rows, errTrainning := db.Debug().Raw(queryTraining, profiles[0].Id).Scan(&exercise).Rows()

	if errTrainning != nil {
		helper.Logger("error", "In Server: "+errTrainning.Error())
		return nil, errors.New(errTrainning.Error())
	}

	defer rows.Close()

	for rows.Next() {
		errScanRows := db.ScanRows(rows, &exercise)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		queryFormExerciseMedia := `SELECT id, path
		FROM form_exercise_medias WHERE exercise_id  = ?`

		rowsFormExerciseMedia, errFormExerciseMedia := db.Debug().Raw(queryFormExerciseMedia, exercise.Id).Scan(&exerciseMedia).Rows()

		if errFormExerciseMedia != nil {
			helper.Logger("error", "In Server: "+errFormExerciseMedia.Error())
			return nil, errors.New(errFormExerciseMedia.Error())
		}

		defer rowsFormExerciseMedia.Close()

		var dataFormExerciseCertificate = make([]entities.FormExerciseCertificate, 0)

		for rowsFormExerciseMedia.Next() {
			errScanFormExerciseMedia := db.ScanRows(rowsFormExerciseMedia, &exerciseMedia)

			if errScanFormExerciseMedia != nil {
				helper.Logger("error", "In Server: "+errScanFormExerciseMedia.Error())
				return nil, errors.New(errScanFormExerciseMedia.Error())
			}

			dataFormExerciseCertificate = append(dataFormExerciseCertificate, entities.FormExerciseCertificate{
				Id:   exerciseMedia.Id,
				Path: exerciseMedia.Path,
			})
		}

		dataTraining = append(dataTraining, entities.ProfileFormExercise{
			Id:           exercise.Id,
			Name:         exercise.Name,
			Institution:  exercise.Institution,
			StartYear:    exercise.StartYear,
			StartMonth:   exercise.StartMonth,
			EndYear:      exercise.EndYear,
			EndMonth:     exercise.EndMonth,
			Certificates: dataFormExerciseCertificate,
		})
	}

	// Work

	var dataWork = make([]entities.ProfileFormWork, 0)

	queryWork := `SELECT id, work, position, institution, is_work, country, city, start_month, start_year, end_month, end_year, user_id 
	FROM form_works WHERE user_id  = ?`

	rows, errWork := db.Debug().Raw(queryWork, profiles[0].Id).Scan(&work).Rows()

	if errWork != nil {
		helper.Logger("error", "In Server: "+errWork.Error())
		return nil, errors.New(errWork.Error())
	}

	defer rows.Close()

	for rows.Next() {
		errScanRows := db.ScanRows(rows, &work)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		var isWork bool

		if work.IsWork == 1 {
			isWork = true
		} else {
			isWork = false
		}

		dataWork = append(dataWork, entities.ProfileFormWork{
			Id:          work.Id,
			Position:    work.Position,
			Institution: work.Institution,
			Work:        work.Work,
			IsWork:      isWork,
			City:        work.City,
			Country:     work.Country,
			StartMonth:  work.StartMonth,
			EndMonth:    work.EndMonth,
			StartYear:   work.StartYear,
			EndYear:     work.EndYear,
		})
	}

	// Language

	var dataLanguage = make([]entities.ProfileFormLanguage, 0)

	queryLanguage := `SELECT id, language, level
	FROM form_languages WHERE user_id  = '` + profiles[0].Id + `'`

	rows, errLanguage := db.Debug().Raw(queryLanguage).Scan(&language).Rows()

	if errLanguage != nil {
		helper.Logger("error", "In Server: "+errLanguage.Error())
		return nil, errors.New(errLanguage.Error())
	}

	defer rows.Close()

	for rows.Next() {
		errScanRows := db.ScanRows(rows, &language)

		if errScanRows != nil {
			helper.Logger("error", "In Server: "+errScanRows.Error())
			return nil, errors.New(errScanRows.Error())
		}

		dataLanguage = append(dataLanguage, language)
	}

	profile.Id = profiles[0].Id
	profile.Avatar = profiles[0].Avatar
	profile.Phone = profiles[0].Phone
	profile.Email = profiles[0].Email
	profile.Fullname = profiles[0].Fullname
	profile.Role = profiles[0].Role
	profile.IsEnabled = enabled
	profile.Job = entities.ProfileJobResponse{
		Id:   profiles[0].JobId,
		Name: profiles[0].JobName,
	}
	profile.Biodata = entities.Biodata{
		Personal: entities.ProfileFormBiodata{
			Id:        profiles[0].BioId,
			Birthdate: profiles[0].BioBirthdate,
			Gender:    profiles[0].BioGender,
			Height:    profiles[0].BioHeight,
			Weight:    profiles[0].BioWeight,
			Religion:  profiles[0].BioReligion,
			Place:     profiles[0].BioPlace,
			Status:    profiles[0].BioStatus,
		},
		Address: entities.ProfileFormPlace{
			Id:            profiles[0].BioAddressId,
			DetailAddress: profiles[0].BioDetailAddress,
			Province: entities.ProfileFormPlaceData{
				Id:   profiles[0].BioProvinceId,
				Name: profiles[0].BioProvince,
			},
			City: entities.ProfileCityPlaceData{
				Id:   profiles[0].BioCityId,
				Name: profiles[0].BioCity,
			},
			District: entities.ProfileDistrictPlaceData{
				Id:   profiles[0].BioDistrictId,
				Name: profiles[0].BioDistrict,
			},
			Subdistrict: entities.ProfileSubdistrictPlaceData{
				Id:   profiles[0].BioSubdistrictId,
				Name: profiles[0].BioSubdistrict,
			},
		},
		Educations:  dataEdu,
		Trainings:   dataTraining,
		Experiences: dataWork,
		Languages:   dataLanguage,
	}

	return map[string]any{
		"data": profile,
	}, nil
}

func ImportUserCreate(ius *entities.ImportUserStore) (map[string]any, error) {
	// INSERT USER
	queryInsertUser := `INSERT INTO users (uid, email, phone, enabled, via) 
		VALUES (?, ?, ?, ?, ?)`
	err := db.Debug().Exec(queryInsertUser, ius.UserId, ius.Email, ius.Phone, 1, "auto").Error
	if err != nil {
		helper.Logger("error", "In Server (InsertUser): "+err.Error())
		return nil, err
	}

	// INSERT PROFILE
	queryInsertProfile := `INSERT INTO profiles (fullname, user_id) VALUES (?, ?)`
	err = db.Debug().Exec(queryInsertProfile, ius.Fullname, ius.UserId).Error
	if err != nil {
		helper.Logger("error", "In Server (InsertProfile): "+err.Error())
		return nil, err
	}

	// INSERT FORM BIODATA
	queryInsertBiodata := `INSERT INTO form_biodatas (birthdate, gender, weight, height, status, religion, place, user_id) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	err = db.Debug().Exec(queryInsertBiodata, ius.Birthdate, ius.Gender, ius.Weight, ius.Height, ius.MaritalStatus, ius.Religion, ius.Place, ius.UserId).Error
	if err != nil {
		helper.Logger("error", "In Server (InsertBiodata): "+err.Error())
		return nil, err
	}

	// INSERT FORM EDUCATION
	// queryInsertEducation := `INSERT INTO form_educations (education_level, major, school_or_college, start_month, start_year, end_month, end_year, user_id)
	// 	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	// err = db.Debug().Exec(queryInsertEducation, ius.EducationLevel, ius.Major, ius.SchoolOrCollege, ius.StartMonthEducation, ius.StartYearEducation, ius.EndMonthEducation, ius.EndYearEducation, ius.UserId).Error
	// if err != nil {
	// 	helper.Logger("error", "In Server (InsertEducation): "+err.Error())
	// 	return nil, err
	// }

	// INSERT FORM EXERCISE
	// queryInsertExercise := `INSERT INTO form_exercises (name, institution, start_month, start_year, end_month, end_year, user_id)
	// 	VALUES (?, ?, ?, ?, ?, ?, ?)`
	// err = db.Debug().Exec(queryInsertExercise, ius.NameInstitution, ius.Institution, ius.StartMonthExercise, ius.StartYearExercise, ius.EndMonthExercise, ius.EndYearExercise, ius.UserId).Error
	// if err != nil {
	// 	helper.Logger("error", "In Server (InsertExercise): "+err.Error())
	// 	return nil, err
	// }

	// INSERT FORM WORK
	// queryInsertWork := `INSERT INTO form_works (position, institution, work, country, city, start_month, start_year, end_month, end_year, is_work, user_id)
	// 	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	// err = db.Debug().Exec(queryInsertWork, ius.PositionWork, ius.InstitutionWork, ius.Work, ius.CountryWork, ius.CityWork, ius.StartMonthWork, ius.StartYearWork, ius.EndMonthWork, ius.EndYearWork, 0, ius.UserId).Error
	// if err != nil {
	// 	helper.Logger("error", "In Server (InsertWork): "+err.Error())
	// 	return nil, err
	// }

	// INSERT FORM LANGUAGE
	// queryInsertLanguage := `INSERT INTO form_languages (level, language, user_id) VALUES (?, ?, ?)`
	// err = db.Debug().Exec(queryInsertLanguage, ius.Level, ius.Language, ius.UserId).Error
	// if err != nil {
	// 	helper.Logger("error", "In Server (InsertLanguage): "+err.Error())
	// 	return nil, err
	// }

	return map[string]any{"message": "Import successful"}, nil
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
