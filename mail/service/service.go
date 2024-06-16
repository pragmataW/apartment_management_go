package service

import (
	"fmt"
	"log"

	"github.com/mailjet/mailjet-apiv3-go"
	"github.com/pragmataW/pragmaTraffic/mail/model"
)

func NewService(apiKey string, secretKey string) MailService {
	return MailService{
		ApiKey: apiKey,
		SecretKey: secretKey,
	}
}

func (s MailService) GetApiKey() string{
	return s.ApiKey
}

func (s MailService) GetSecretKey() string{
	return s.SecretKey
}

func (s MailService) SendMail(cnf model.MailConfig) error {
	client := mailjet.NewMailjetClient(cnf.ApiKey, cnf.SecretKey)
	messageInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: cnf.FromEmail,
				Name:  cnf.FromName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: cnf.ToEmail,
					Name:  cnf.ToName,
				},
			},
			Subject:  cnf.Subject,
			HTMLPart: cnf.Html,
		},
	}
	message := mailjet.MessagesV31{Info: messageInfo}
	res, err := client.SendMailV31(&message)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("Response: %+v\n", res)
	return nil
}