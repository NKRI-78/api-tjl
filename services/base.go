package services

import (
	"fmt"
	"os"
	helpers "superapps/helpers"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		helpers.Logger("error", "Error getting env")
	}

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbDriver := os.Getenv("DB_DRIVER")

	dbURI := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_0900_ai_ci&parseTime=True&loc=Asia%%2FJakarta",
		username, password, dbHost, dbPort, dbName,
	)

	conn, err := gorm.Open(dbDriver, dbURI)
	if err != nil {
		helpers.Logger("error", "In Server: "+err.Error())
	}

	os.Setenv("TZ", "Asia/Jakarta")

	db = conn
	db.Debug().AutoMigrate()
}

func GetDB() *gorm.DB {
	return db
}
