package helper

// import (
// 	"net/mail"
// )

import (
	// "math/rand"
	"bytes"
	crand "crypto/rand"
	"encoding/base32"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"regexp"
	"strings"

	// "time"
	entities "superapps/entities"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SendEmail(to, subject, body string) error {
	emailData := &entities.SendEmailRequest{
		To:      to,
		App:     "TJL",
		Subject: subject,
		Body:    body,
	}

	jsonData, err := json.Marshal(emailData)
	if err != nil {
		return err
	}

	resp, err := http.Post("https://api-email.inovatiftujuh8.com/api/v1/email", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to send email, status code: " + resp.Status)
	}

	return nil
}

func DecodeJwt(tokenP string) *jwt.Token {
	splitted := strings.Split(tokenP, " ")

	tokenPart := splitted[1]

	token, _ := jwt.Parse(tokenPart, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	return token
}

// NOT SECURE FOR USE
// func CodeOtp() string {

// 	rand.Seed(time.Now().UTC().UnixNano())

// 	b := make([]rune, 4)
// 	l := len(letterRunes)
// 	for i := range b {
// 		b[i] = letterRunes[rand.Intn(l)]
// 	}

// 	return string(b)
// }

func CodeOtpSecure() string {
	// Calculate the number of random bytes needed for the specified length
	numBytes := (4 * 5) / 8

	// Generate random bytes
	randomBytes := make([]byte, numBytes)
	_, err := crand.Read(randomBytes)
	if err != nil {
		return ""
	}

	// Encode the random bytes to base32
	otp := base32.StdEncoding.EncodeToString(randomBytes)

	// Truncate to the desired length
	otp = otp[:4]

	return string(otp)
}

func IsValidEmail(email string) bool {
	// Optional 1
	// _, err := mail.ParseAddress(email)
	// return err == nil

	emailRegex := regexp.MustCompile(`^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	return emailRegex.MatchString(email)
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
