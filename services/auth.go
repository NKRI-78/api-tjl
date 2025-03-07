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

func VerifyOtp(u *models.User) (map[string]interface{}, error) {
	var user entities.UserOtp

	// Gunakan parameterized query untuk mencegah SQL Injection
	err := db.Debug().Raw(`
		SELECT uid, enabled, created_at 
		FROM users 
		WHERE (email = ? OR phone = ?) AND otp = ?`, u.Val, u.Val, u.Otp).
		First(&user).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
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

	// Update status akun dengan parameterized query
	errUpdate := db.Debug().Exec(`
		UPDATE users SET enabled = 1, email_active_date = NOW() 
		WHERE uid = ?`, user.Uid).Error

	if errUpdate != nil {
		helper.Logger("error", "In Server: "+errUpdate.Error())
		return nil, errUpdate
	}

	// Buat token setelah akun diaktifkan
	token, err := middleware.CreateToken(user.Uid)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	return map[string]any{"token": token["token"]}, nil
}

func ResendOtp(u *models.User) (map[string]interface{}, error) {

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
		errUpdateResendOtp := db.Debug().Exec(`UPDATE users SET otp = '` + otp + `', otp_date = NOW() WHERE email = '` + u.Val + `'`).Error

		if errUpdateResendOtp != nil {
			helper.Logger("error", "In Server: "+errUpdateResendOtp.Error())
			return nil, errors.New(errUpdateResendOtp.Error())
		}

		errEmail := helper.SendEmail(u.Val, "TJL", otp)
		if errEmail != nil {
			helper.Logger("error", "Failed to send email: "+errEmail.Error())
		}
	}

	return map[string]interface{}{
		"otp": otp,
	}, nil
}

func Login(u *models.User) (map[string]interface{}, error) {

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

		errEmail := helper.SendEmail(u.Val, "TJL", otp)
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

	token, err := middleware.CreateToken(user.Id)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	access := token["token"]

	return map[string]any{"token": access}, nil
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

	errEmail := helper.SendEmail(user.Email, "TJL", otp)
	if errEmail != nil {
		helper.Logger("error", "Failed to send email: "+errEmail.Error())
	}

	token, err := middleware.CreateToken(user.Id)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}

	access := token["token"]

	return map[string]any{"token": access}, nil
}



