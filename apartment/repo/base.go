package repo

import (
	"fmt"
	"log"
	"sync"

	"github.com/pragmataW/apartment_management/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

type repo struct {
	db *gorm.DB
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	SslMode  string
}

func NewDb(conf DBConfig) (*gorm.DB) {
	once.Do(func() {
		var err error
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", conf.Host, conf.User, conf.Password, conf.DbName, conf.Port, conf.SslMode)
		log.Println(dsn)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil{
			log.Fatal(err)
		}
		err = db.AutoMigrate(&models.Apartment{})
		if err != nil{
			log.Fatal(err)
		}
		err = db.AutoMigrate(&models.Announcement{})
		if err != nil{
			log.Fatal(err)
		}
		err = db.AutoMigrate(&models.Merchant{})
		if err != nil{
			log.Fatal(err)
		}
	})
	return db
}

func NewRepo(db *gorm.DB) repo {
	return repo{
		db: db,
	}
}
