package services

import (
	helper "superapps/helpers"
	"fmt"
	"os"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {

	env := godotenv.Load()

	if env != nil {
		helper.Logger("error", "Error getting env")
	}

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbDriver := os.Getenv("DB_DRIVER")

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, dbHost, dbPort, dbName)

	conn, err := gorm.Open(dbDriver, dbURI)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
	}

	db = conn
	db.Debug().AutoMigrate() 
}

func GetDB() *gorm.DB {
	return db
}
