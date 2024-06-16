package controller

import "github.com/pragmataW/pragmaTraffic/mail/model"

type IMailService interface {
	SendMail(model.MailConfig) error
	GetApiKey() string
	GetSecretKey() string
}

type MailController struct {
	Src IMailService
}

type RequestBody struct {
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
	ToName    string `json:"to_name"`
	ToEmail   string `json:"to_email"`
	Subject   string `json:"subject"`
	Html      string `json:"html"`
}
