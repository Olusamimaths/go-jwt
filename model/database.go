package model

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var DB *gorm.DB
var err error

var secretKey = envVariable("SECRET")

func envVariable(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func SetDBClient() {
	var (
		host     = envVariable("DB_HOST")
		port     = envVariable("DB_PORT")
		user     = envVariable("DB_USER")
		dbname   = envVariable("DB_NAME")
		password = envVariable("DB_PASSWORD")
	)

	dns := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		user,
		dbname,
		password,
	)

	DB, err = gorm.Open("postgres", dns)
	DB.AutoMigrate(User{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connection to DB successful")
}

type User struct {
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

func (u *User) GeneratePasswordHash() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	u.Password = string(bytes)
	return err
}

func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
