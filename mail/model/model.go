package model

type MailConfig struct {
	FromName  string
	FromEmail string
	ToName    string
	ToEmail   string
	Subject   string
	Html      string
	ApiKey    string
	SecretKey string
}
