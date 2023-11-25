package utils

import (
	"fmt"
	"log"
	"os"

	model "basictrade/models"

	"github.com/google/uuid"
	// "github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func StartDB() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	host := os.Getenv("HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbport := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	if host == "" || user == "" || password == "" || dbport == "" || dbname == "" {
		log.Fatal("Incomplete database configuration. Please check your environment variables.")
	}	

	// "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	config := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, dbport, dbname)
	db, err = gorm.Open(mysql.Open(config), &gorm.Config{})
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}
	

	// AutoMigrate models
	db.Debug().AutoMigrate(
		&model.Admin{}, 
		&model.Product{}, 
		&model.Variant{},
	)

	db.Callback().Create().Before("gorm:before_create").Register("before_create", BeforeCreateUUID)

	fmt.Println("Connected to the database")
}

func GetDB() *gorm.DB {
	return db
}

// BeforeCreateUUID is a callback to set UUIDs before creating records
func BeforeCreateUUID(db *gorm.DB) {
    if _, ok := db.Statement.Schema.FieldsByName["uuid"]; ok {
        db.Statement.SetColumn("uuid", uuid.New().String())
    }
}
