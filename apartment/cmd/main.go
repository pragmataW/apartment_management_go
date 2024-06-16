package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pragmataW/apartment_management/controller"
	configmanager "github.com/pragmataW/apartment_management/pkg/config_manager"
	"github.com/pragmataW/apartment_management/pkg/encrypt"
	"github.com/pragmataW/apartment_management/repo"
	"github.com/pragmataW/apartment_management/services"
)

var (
	host     string
	port     int
	user     string
	password string
	dbName   string
	sslMode  string
	chiper   string
	jwtKey   string
)

func main() {
	repoCfg := repo.DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DbName:   dbName,
		SslMode:  sslMode,
	}

	repo := repo.NewRepo(
		repo.NewDb(repoCfg),
	)
	encrypt := encrypt.NewEncryptor(chiper)
	cfgManager := configmanager.NewConfigManager()
	service := services.NewService(
		services.WithConfigManager(cfgManager),
		services.WithRepo(repo),
		services.WithEncryptor(encrypt),
	)
	ctrl := controller.NewController(
		controller.WithConfigManager(cfgManager),
		controller.WithService(service),
	)

	go func() {
		err := service.IncreaseDuesAutomatically()
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
	}()

	app := fiber.New()
	ctrl.RegisterRoutes(app, jwtKey)
	log.Fatal(app.Listen(":2009"))
}

func init() {
	host = os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	port, _ = strconv.Atoi(portStr)
	user = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbName = os.Getenv("POSTGRES_DB")
	sslMode = os.Getenv("POSTGRES_SSL")
	chiper = os.Getenv("CHIPER")
	jwtKey = os.Getenv("JWT_KEY")
}
