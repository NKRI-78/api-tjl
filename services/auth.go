package services

import (
	"errors"
	entities "superapps/entities"
	helper "superapps/helpers"
	middleware "superapps/middlewares"
	models "superapps/models"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func DeleteUser(d *entities.DeleteUser) (map[string]any, error) {

	query := `DELETE FROM users WHERE uid = ?`

	err := db.Debug().Exec(query, d.Id).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{}, nil
}

func RoleList() (map[string]any, error) {
	var userRole []entities.UserRoles

	query := `SELECT id, name FROM user_roles`

	err := db.Debug().Raw(query).Scan(&userRole).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	return map[string]any{
		"data": userRole,
	}, nil
}

func RegisterUserBranch(rub *entities.RegisterUserBranch) (map[string]any, error) {
	Id := uuid.NewV4().String()

	hashedPassword, err := helper.Hash(rub.Password)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	// INSERT USER
	queryInsertUser := `INSERT INTO users (uid, email, phone, password, role, enabled) VALUES (?, ?, ?, ?, ?, ?)`

	errInsertUser := db.Debug().Exec(queryInsertUser, Id, rub.Email, rub.Phone, hashedPassword, rub.RoleId, 1).Error

	if errInsertUser != nil {
		helper.Logger("error", "In Server: "+errInsertUser.Error())
		return nil, errors.New(errInsertUser.Error())
	}

	// INSERT PROFILE
	queryInsertProfile := `INSERT INTO profiles (fullname, user_id) VALUES (?, ?)`

	errInsertUserProfile := db.Debug().Exec(queryInsertProfile, rub.Fullname, Id).Error

	if errInsertUserProfile != nil {
		helper.Logger("error", "In Server: "+errInsertUserProfile.Error())
		return nil, errors.New(errInsertUserProfile.Error())
	}

	// INSERT USER BRANCH
	queryInsertUserBranch := `INSERT INTO user_branches (branch_id, user_id) VALUES (?, ?)`

	errInsertUserBranch := db.Debug().Exec(queryInsertUserBranch, rub.BranchId, Id).Error

	if errInsertUserBranch != nil {
		helper.Logger("error", "In Server: "+errInsertUserBranch.Error())
		return nil, errors.New(errInsertUserBranch.Error())
	}

	helper.SendEmail(rub.Email, "TJL", "Registrasi Berhasil", rub.Password, "tjl-create-user-branch")

	return map[string]any{}, nil
}

func UpdateUserBranch(uub *entities.UpdateUserBranch) (map[string]any, error) {
	// users := []entities.CheckAccount{}

	hashedPassword, err := helper.Hash(uub.Password)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	// CHECK EMAIL only if provided
	// if uub.Email != "" {
	// errCheckEmail := db.Debug().Raw(`SELECT email FROM users WHERE email = ? AND enabled = 1`, uub.Email).Scan(&users).Error
	// if errCheckEmail != nil {
	// 	helper.Logger("error", "In Server: "+errCheckEmail.Error())
	// 	return nil, errors.New(errCheckEmail.Error())
	// }

	// if len(users) == 1 {
	// 	return nil, errors.New("E-mail already exists")
	// }
	// }

	// UPDATE USER (with or without email depending on whether it's provided)
	var errUpdateUser error
	// if uub.Email != "" {
	// 	errUpdateUser = db.Debug().Exec(`
	// 		UPDATE users SET password = ?, phone = ?, role = ?, updated_at = NOW()
	// 		WHERE uid = ?`, hashedPassword, uub.Phone, uub.RoleId, uub.Id).Error
	// } else {
	errUpdateUser = db.Debug().Exec(`
			UPDATE users SET password = ?, phone = ?, role = ?, updated_at = NOW() 
			WHERE uid = ?`, hashedPassword, uub.Phone, uub.RoleId, uub.Id).Error
	// }

	if errUpdateUser != nil {
		helper.Logger("error", "In Server: "+errUpdateUser.Error())
		return nil, errUpdateUser
	}

	// UPDATE PROFILE
	errUpdateProfile := db.Debug().Exec(`
		UPDATE profiles SET updated_at = NOW(), fullname = ?
		WHERE user_id = ?`, uub.Fullname, uub.Id).Error

	if errUpdateProfile != nil {
		helper.Logger("error", "In Server: "+errUpdateProfile.Error())
		return nil, errUpdateProfile
	}

	helper.SendEmail(uub.Email, "TJL", "Update Akun Berhasil", uub.Password, "tjl-update-user-branch")

	return map[string]any{}, nil
}

func ForgotPassword(fp *entities.ForgotPassword) (map[string]any, error) {
	users := []entities.ForgotPassword{}

	err := db.Debug().Raw(`SELECT email FROM users WHERE email = ?`, fp.Email).Scan(&users).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	isUserExist := len(users)

	if isUserExist == 0 {
		return nil, errors.New("USER_NOT_FOUND")
	}

	hashedPassword, errHashedPassword := helper.Hash(fp.NewPassword)

	if errHashedPassword != nil {
		helper.Logger("error", "In Server: "+errHashedPassword.Error())
		return nil, errHashedPassword
	}

	errUpdate := db.Debug().Exec(`UPDATE users SET password = ? WHERE email = ?`, hashedPassword, fp.Email).Error

	if errUpdate != nil {
		helper.Logger("error", "In Server: "+errUpdate.Error())
		return nil, errUpdate
	}

	return map[string]any{}, nil
}

func UpdateEmail(ue *models.UpdateEmail) (map[string]any, error) {
	users := []entities.CheckAccount{}

	err := db.Debug().Raw(`SELECT email FROM users WHERE email = ? AND enabled = 1`, ue.NewEmail).Scan(&users).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	isUserExist := len(users)

	if isUserExist == 1 {
		return nil, errors.New("EMAIL_ALREADY_EXIST_AND_ACTIVE")
	}

	otp := helper.CodeOtpSecure()

	errUpdate := db.Debug().Exec(`
		UPDATE users SET email = ?, otp = ?, created_at = NOW()
		WHERE email = ? AND enabled = 0`, ue.NewEmail, otp, ue.OldEmail).Error

	if errUpdate != nil {
		helper.Logger("error", "In Server: "+errUpdate.Error())
		return nil, errUpdate
	}

	errEmail := helper.SendEmail(ue.NewEmail, "TJL", "Verification Account", otp, "-")
	if errEmail != nil {
		helper.Logger("error", "Failed to send email: "+errEmail.Error())
	}

	return map[string]any{}, nil
}

func VerifyOtp(u *models.User) (map[string]any, error) {
	var user entities.UserOtp

	err := db.Debug().Raw(`
		SELECT uid, enabled, created_at 
		FROM users 
		WHERE (email = ? OR phone = ?) AND otp = ?`, u.Val, u.Val, u.Otp).
		First(&user).Error

	if err != nil {
		helper.Logger("error", "In Server: USER_OR_OTP_IS_INVALID")
		return nil, errors.New("USER_OR_OTP_IS_INVALID")
	}

	if user.Enabled == 1 {
		helper.Logger("error", "In Server: Account is already active")
		return nil, errors.New("ACCOUNT_IS_ALREADY_ACTIVE")
	}

	// Cek expired OTP (lebih efisien dengan time.Since)
	if time.Since(user.CreatedAt.UTC()) >= time.Minute {
		helper.Logger("error", "In Server: Otp is expired")
		return nil, errors.New("OTP_IS_EXPIRED")
	}

	errUpdate := db.Debug().Exec(`
		UPDATE users SET enabled = 1, email_active_date = NOW() 
		WHERE uid = ?`, user.Uid).Error

	if errUpdate != nil {
		helper.Logger("error", "In Server: "+errUpdate.Error())
		return nil, errUpdate
	}

	token, err := middleware.CreateToken("-", user.Uid)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	return map[string]any{"token": token["token"]}, nil
}

func ResendOtp(u *models.User) (map[string]any, error) {

	users := []entities.UserOtp{}
	query := `SELECT enabled, otp_date FROM users
	WHERE (email = '` + u.Val + `' OR phone = '` + u.Val + `')`

	err := db.Debug().Raw(query).Scan(&users).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	isUserExist := len(users)

	if isUserExist == 0 {
		return nil, errors.New("USER_NOT_FOUND")
	}

	emailActive := users[0].Enabled
	createdAt := users[0].CreatedAt

	if emailActive == 1 {
		helper.Logger("error", "In Server: Account is already active")
		return nil, errors.New("ACCOUNT_IS_ALREADY_ACTIVE")
	}

	currentTime := time.Now()
	elapsed := currentTime.Sub(createdAt)

	otp := helper.CodeOtpSecure()

	if elapsed >= 1*time.Minute {
		errUpdateResendOtp := db.Debug().Exec(`UPDATE users SET otp = '` + otp + `', created_at = NOW(), otp_date = NOW() WHERE email = '` + u.Val + `'`).Error

		if errUpdateResendOtp != nil {
			helper.Logger("error", "In Server: "+errUpdateResendOtp.Error())
			return nil, errors.New(errUpdateResendOtp.Error())
		}

		errEmail := helper.SendEmail(u.Val, "TJL", "Verification Account", otp, "-")
		if errEmail != nil {
			helper.Logger("error", "Failed to send email: "+errEmail.Error())
		}
	}

	return map[string]any{
		"otp": otp,
	}, nil
}

func Login(u *models.User) (map[string]any, error) {

	user := entities.User{}

	users := []entities.UserLogin{}
	query := `SELECT uid, enabled, password FROM users WHERE email = '` + u.Val + `' OR phone = '` + u.Val + `'`

	err := db.Debug().Raw(query).Scan(&users).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	isUserExist := len(users)

	if isUserExist == 0 {
		return nil, errors.New("USER_NOT_FOUND")
	}

	otp := helper.CodeOtpSecure()

	emailActive := users[0].Enabled
	user.Id = users[0].Uid

	if emailActive == 0 {
		err := db.Debug().Exec(`UPDATE users SET otp = '` + otp + `', otp_date = NOW()
		WHERE email = '` + u.Val + `' OR phone = '` + u.Val + `'`).Error

		if err != nil {
			helper.Logger("error", "In Server: "+err.Error())
			return nil, errors.New(err.Error())
		}

		errEmail := helper.SendEmail(u.Val, "TJL", "Verification Account", otp, "-")
		if errEmail != nil {
			helper.Logger("error", "Failed to send email: "+errEmail.Error())
		}

		helper.Logger("error", "In Server: Please activate your account")
		return nil, errors.New("PLEASE_ACTIVATE_YOUR_ACCOUNT")
	}

	passHashed := users[0].Password

	err = helper.VerifyPassword(passHashed, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New("CREDENTIALS_IS_INCORRECT")
	}

	token, err := middleware.CreateToken("-", user.Id)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	access := token["token"]

	return map[string]any{"token": access}, nil
}

func LoginAdmin(u *models.UserAdmin) (entities.AdminResponse, error) {

	users := []entities.UserAdmin{}

	query := `SELECT u.uid AS user_id, b.id AS branch_id, u.enabled, u.password, p.fullname, p.avatar, ur.name AS role
	FROM users u
	INNER JOIN profiles p ON p.user_id = u.uid
	LEFT JOIN user_branches ub ON ub.user_id = u.uid
	LEFT JOIN branchs b ON b.id = ub.branch_id
	INNER JOIN user_roles ur ON ur.id = u.role
	WHERE u.email = '` + u.Val + `' OR u.phone = '` + u.Val + `' 
	LIMIT 1`

	err := db.Debug().Raw(query).Scan(&users).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return entities.AdminResponse{}, errors.New(err.Error())
	}

	isUserExist := len(users)

	if isUserExist == 0 {
		return entities.AdminResponse{}, errors.New("USER_NOT_FOUND")
	}

	passHashed := users[0].Password

	err = helper.VerifyPassword(passHashed, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		helper.Logger("error", "In Server: "+err.Error())
		return entities.AdminResponse{}, errors.New("CREDENTIALS_IS_INCORRECT")
	}

	token, err := middleware.CreateToken(users[0].BranchId, users[0].UserId)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return entities.AdminResponse{}, err
	}

	access := token["token"]

	return entities.AdminResponse{
		ID:       users[0].UserId,
		Avatar:   users[0].Avatar,
		Fullname: users[0].Fullname,
		Role:     users[0].Role,
		Token:    access,
	}, nil
}

func Register(u *models.User) (map[string]any, error) {

	hashedPassword, err := helper.Hash(u.Password)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	user := entities.User{}

	user.Id = uuid.NewV4().String()

	user.JobId = u.JobId
	user.BranchId = u.BranchId
	user.Avatar = u.Avatar
	user.Fullname = u.Fullname
	user.Email = u.Email
	user.Phone = u.Phone
	user.Password = string(hashedPassword)

	otp := helper.CodeOtpSecure()

	users := []entities.CheckAccount{}
	jobs := []entities.CheckJobs{}

	errCheckAccount := db.Debug().Raw(`SELECT email FROM users WHERE email = '` + u.Email + `'`).Scan(&users).Error

	if errCheckAccount != nil {
		helper.Logger("error", "In Server: "+errCheckAccount.Error())
		return nil, errors.New(errCheckAccount.Error())
	}

	errCheckJobs := db.Debug().Raw(`SELECT uid FROM job_categories WHERE uid = '` + u.JobId + `'`).Scan(&jobs).Error

	if errCheckJobs != nil {
		helper.Logger("error", "In Server: "+errCheckJobs.Error())
		return nil, errors.New(errCheckJobs.Error())
	}

	isUserExist := len(users)

	if isUserExist == 1 {
		helper.Logger("error", "In Server: User already exist")
		return nil, errors.New("USER_ALREADY_EXIST")
	}

	isJobExist := len(jobs)

	if isJobExist == 0 {
		helper.Logger("error", "In Server: Job not found")
		return nil, errors.New("JOB_NOT_FOUND")
	}

	errInsertUser := db.Debug().Exec(`INSERT INTO users (uid, email, phone, password, otp) 
	VALUES ('` + user.Id + `', '` + user.Email + `', '` + user.Phone + `', '` + user.Password + `', '` + otp + `')`).Error

	if errInsertUser != nil {
		helper.Logger("error", "In Server: "+errInsertUser.Error())
		return nil, errors.New(errInsertUser.Error())
	}

	errInsertProfile := db.Debug().Exec(`INSERT INTO profiles (user_id, fullname, avatar) VALUES ('` + user.Id + `', '` + user.Fullname + `', '` + user.Avatar + `')`).Error

	if errInsertProfile != nil {
		helper.Logger("error", "In Server: "+errInsertProfile.Error())
		return nil, errors.New(errInsertProfile.Error())
	}

	errInsertUserJobPick := db.Debug().Exec(`INSERT INTO user_pick_jobs (user_id, job_id) VALUES ('` + user.Id + `', '` + user.JobId + `')`).Error

	if errInsertUserJobPick != nil {
		helper.Logger("error", "In Server: "+errInsertUserJobPick.Error())
		return nil, errors.New(errInsertUserJobPick.Error())
	}

	errInsertUserBranch := db.Debug().Exec(`INSERT INTO user_branches (user_id, branch_id) VALUES ('` + user.Id + `', '` + user.BranchId + `')`).Error

	if errInsertUserBranch != nil {
		helper.Logger("error", "In Server: "+errInsertUserBranch.Error())
		return nil, errors.New(errInsertUserBranch.Error())
	}

	errEmail := helper.SendEmail(user.Email, "TJL", "Verification Account", otp, "-")
	if errEmail != nil {
		helper.Logger("error", "Failed to send email: "+errEmail.Error())
	}

	token, err := middleware.CreateToken("-", user.Id)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}

	access := token["token"]

	return map[string]any{"token": access}, nil
}
