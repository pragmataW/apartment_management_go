package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pragmataW/pragmaTraffic/mail/model"
)

func NewController(src IMailService) MailController {
	return MailController{
		Src: src,
	}
}

func (ctrl MailController) PostMail(c *fiber.Ctx) error {
	var req RequestBody

	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Message": "Bad Request",
		})
	}

	conf := model.MailConfig{
		FromName:  req.FromName,
		FromEmail: req.FromEmail,
		ToName:    req.ToName,
		ToEmail:   req.ToEmail,
		Subject:   req.Subject,
		Html:      req.Html,
		ApiKey:    ctrl.Src.GetApiKey(),
		SecretKey: ctrl.Src.GetSecretKey(),
	}

	err = ctrl.Src.SendMail(conf)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"Message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"Message": "Mail sent",
	})
}
