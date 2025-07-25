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
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	// "time"
	entities "superapps/entities"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func DefaultIfEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// Ex : Kamis, 20 Maret 2025 13:24 WIB
func FormatDate(t time.Time) string {
	days := map[string]string{
		"Sunday":    "Minggu",
		"Monday":    "Senin",
		"Tuesday":   "Selasa",
		"Wednesday": "Rabu",
		"Thursday":  "Kamis",
		"Friday":    "Jumat",
		"Saturday":  "Sabtu",
	}

	months := map[string]string{
		"January":   "Januari",
		"February":  "Februari",
		"March":     "Maret",
		"April":     "April",
		"May":       "Mei",
		"June":      "Juni",
		"July":      "Juli",
		"August":    "Agustus",
		"September": "September",
		"October":   "Oktober",
		"November":  "November",
		"December":  "Desember",
	}

	day := days[t.Weekday().String()]
	month := months[t.Month().String()]
	return fmt.Sprintf("%s, %02d %s %d %02d:%02d WIB", day, t.Day(), month, t.Year(), t.Hour(), t.Minute())
}



func SendEmail(to, app, subject, data, Type string) error {

	body := data

	emailData := &entities.SendEmailRequest{
		To:      to,
		App:     app,
		Subject: subject,
		Body:    body,
		Type:    Type,
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
	} else {
		Logger("info", "Send Email Success")
	}

	return nil
}

func TimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "Beberapa detik yang lalu"
	case duration < time.Hour:
		return fmt.Sprintf("%d menit yang lalu", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%d jam yang lalu", int(duration.Hours()))
	case duration < 7*24*time.Hour:
		return fmt.Sprintf("%d hari yang lalu", int(duration.Hours()/24))
	case duration < 30*24*time.Hour:
		return fmt.Sprintf("%d minggu yang lalu", int(duration.Hours()/(24*7)))
	case duration < 365*24*time.Hour:
		return fmt.Sprintf("%d bulan yang lalu", int(duration.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%d tahun yanh lalu", int(duration.Hours()/(24*365)))
	}
}

func SendFcm(title, message, token, BroadcastType, JobAppId string) error {

	FcmData := &entities.SendFcmRequest{
		Token:         token,
		Title:         title,
		Body:          message,
		JobAppId:      JobAppId,
		BroadcastType: BroadcastType,
	}

	jsonData, err := json.Marshal(FcmData)
	if err != nil {
		return err
	}

	resp, err := http.Post("https://api-fcm-office.inovatiftujuh8.com/api/v1/firebase/fcm", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to send fcm, status code: " + resp.Status)
	}

	return nil
}

func DecodeJwt(val string) *jwt.Token {
	splitted := strings.Split(val, " ")

	tokenPart := splitted[1]

	token, _ := jwt.Parse(tokenPart, func(token *jwt.Token) (any, error) {
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

func FormatIDR(amount float64) string {
	// Convert amount to a string with no decimal places
	amountStr := fmt.Sprintf("%.0f", amount)
	n := len(amountStr)

	if n <= 3 {
		return "Rp." + amountStr
	}

	var result []string
	for i, c := range amountStr {
		if (n-i)%3 == 0 && i != 0 {
			result = append(result, ".")
		}
		result = append(result, string(c))
	}

	return "Rp." + strings.Join(result, "")
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
