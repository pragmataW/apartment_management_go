package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/pragmataW/pragmaTraffic/mail/controller"
	"github.com/pragmataW/pragmaTraffic/mail/service"
)

var (
	apiKey    string
	secretKey string
)

func main() {
	src := service.NewService(apiKey, secretKey)
	ctrl := controller.NewController(src)

	app := fiber.New()
	app.Post("/sendMail", ctrl.PostMail)

	log.Fatal(app.Listen(":2001"))
}

func init() {
	apiKey = os.Getenv("MAILJET_API_KEY")
	secretKey = os.Getenv("MAILJET_SECRET_KEY")
}
