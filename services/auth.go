package services

import (
	"errors"
	"fmt"
	entities "superapps/entities"
	helper "superapps/helpers"
	middleware "superapps/middlewares"
	models "superapps/models"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func VerifyOtp(u *models.User) (map[string]interface{}, error) {

	users := []entities.UserOtp{}
	query := `SELECT uid, email_active, otp_date FROM users 
	WHERE (email = '` + u.Val + `' OR phone = '` + u.Val + `') AND otp = '` + u.Otp + `'`

	fmt.Println((query))

	err := db.Debug().Raw(query).Scan(&users).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	isUserExist := len(users)

	if isUserExist == 0 {
		return nil, errors.New("user not found")
	}

	uid := users[0].Uid
	emailActive := users[0].EmailActive
	otpDate := users[0].OtpDate

	if emailActive == 1 {
		helper.Logger("error", "In Server: Account is already active")
		return nil, errors.New("account is already active")
	}

	currentTime := time.Now()
	elapsed := currentTime.Sub(otpDate)

	if elapsed >= 1*time.Minute {
		helper.Logger("error", "In Server: Otp is expired")
		return nil, errors.New("otp is expired")
	}

	errUpdateEmailActive := db.Debug().Exec(`UPDATE users SET email_active = 1, email_active_date = NOW()
		WHERE email = '` + u.Val + `'
	`).Error

	if errUpdateEmailActive != nil {
		helper.Logger("error", "In Server: "+errUpdateEmailActive.Error())
		return nil, errors.New(errUpdateEmailActive.Error())
	}

	token, err := middleware.CreateToken(uid)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	access := token["token"]

	return map[string]interface{}{"token": access}, nil
}

func ResendOtp(u *models.User) (map[string]interface{}, error) {

	users := []entities.UserOtp{}
	query := `SELECT email_active, otp_date FROM users
	WHERE (email = '` + u.Val + `' OR phone = '` + u.Val + `')`

	err := db.Debug().Raw(query).Scan(&users).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	isUserExist := len(users)

	if isUserExist == 0 {
		return nil, errors.New("user not found")
	}

	emailActive := users[0].EmailActive
	otpDate := users[0].OtpDate

	if emailActive == 1 {
		helper.Logger("error", "In Server: Account is already active")
		return nil, errors.New("account is already active")
	}

	currentTime := time.Now()
	elapsed := currentTime.Sub(otpDate)

	otp := helper.CodeOtpSecure()

	if elapsed >= 1*time.Minute {
		errUpdateResendOtp := db.Debug().Exec(`UPDATE users SET otp = '` + otp + `', otp_date = NOW() WHERE email = '` + u.Val + `'`).Error

		if errUpdateResendOtp != nil {
			helper.Logger("error", "In Server: "+errUpdateResendOtp.Error())
			return nil, errors.New(errUpdateResendOtp.Error())
		}
	}

	return map[string]interface{}{
		"otp": otp,
	}, nil
}

func Login(u *models.User) (map[string]interface{}, error) {

	user := entities.User{}

	users := []entities.UserLogin{}
	query := `SELECT uid, email_active, password FROM users WHERE email = '` + u.Val + `' OR phone = '` + u.Val + `'`

	err := db.Debug().Raw(query).Scan(&users).Error

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New(err.Error())
	}

	isUserExist := len(users)

	if isUserExist == 0 {
		return nil, errors.New("user not found")
	}

	emailActive := users[0].EmailActive
	user.Id = users[0].Uid

	if emailActive == 0 {
		err := db.Debug().Exec(`UPDATE users SET otp = '` + helper.CodeOtpSecure() + `', otp_date = NOW() 
		WHERE email = '` + u.Val + `' OR phone = '` + u.Val + `'`).Error

		if err != nil {
			helper.Logger("error", "In Server: "+err.Error())
			return nil, errors.New(err.Error())
		}

		helper.Logger("error", "In Server: Please activate your account")
		return nil, errors.New("please activate your account")
	}

	passHashed := users[0].Password

	err = helper.VerifyPassword(passHashed, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, errors.New("credentials is incorrect")
	}

	token, err := middleware.CreateToken(user.Id)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	access := token["token"]

	return map[string]interface{}{"token": access}, nil
}

func Register(u *models.User) (map[string]interface{}, error) {

	hashedPassword, err := helper.Hash(u.Password)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		return nil, err
	}

	user := entities.User{}

	user.Id = uuid.NewV4().String()

	user.Avatar = u.Avatar
	user.Fullname = u.Fullname
	user.Email = u.Email
	user.Phone = u.Phone
	user.Password = string(hashedPassword)

	otp := helper.CodeOtpSecure()

	users := []entities.CheckAccount{}
	errCheckAccount := db.Debug().Raw(`SELECT email FROM users WHERE email = '` + u.Email + `'`).Scan(&users).Error

	if errCheckAccount != nil {
		helper.Logger("error", "In Server: "+errCheckAccount.Error())
		return nil, errors.New(errCheckAccount.Error())
	}

	isUserExist := len(users)

	if isUserExist == 1 {
		helper.Logger("error", "In Server: User already exist")
		return nil, errors.New("user already exist")
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

	emailBody := `<div style="font-family: Helvetica,Arial,sans-serif;min-width:1000px;overflow:auto;line-height:2">
	<div style="margin:50px auto; width:70%; padding: 20px 0;">
	<p style="font-size:1.1em;">Hi,</p>
	<p>Use the following OTP to complete your Sign Up procedures. OTP is valid for 2 minutes</p>
	<h2 style="background: #00466a; margin: 0 auto; width: max-content; padding: 0 10px; color: #fff; border-radius: 4px;">` + otp + `</h2>
	<p style="font-size:0.9em;">Regards, <br/>TJL</p>
	<hr style="border:none;border-top:1px solid #eee" />
	</div></div>`
	
    errEmail := helper.SendEmail(user.Email, "TJL", emailBody)
    if errEmail != nil {
        helper.Logger("error", "Failed to send email: "+errEmail.Error())
    }

	token, err := middleware.CreateToken(user.Id)
	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}

	access := token["token"]

	return map[string]interface{}{"token": access}, nil
}