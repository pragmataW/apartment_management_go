package configmanager

import (
	"log"
	"os"
	"strconv"
)

type configManager struct {
	adminPassword string
	jwtKey        string
	fromMail      string
	mailServer    string
	merchantID    int
	merchantKey   string
	merchantSalt  string
	okUrl         string
	failUrl       string
}

func NewConfigManager() configManager {
	merchantId, err := strconv.Atoi(os.Getenv("MERCHANT_ID"))
	if err != nil {
		log.Fatal(err)
	}

	return configManager{
		adminPassword: os.Getenv("ADMIN_PASS"),
		jwtKey:        os.Getenv("JWT_KEY"),
		fromMail:      os.Getenv("FROM_MAIL"),
		mailServer:    os.Getenv("MAIL_SERVER"),
		merchantID:    merchantId,
		merchantKey:   os.Getenv("MERCHANT_KEY"),
		merchantSalt:  os.Getenv("MERCHANT_SALT"),
		okUrl:         os.Getenv("PAYMENT_OK_URL"),
		failUrl:       os.Getenv("PAYMENT_FAIL_URL"),
	}
}

func (c configManager) GetAdminPassword() string {
	return c.adminPassword
}

func (c configManager) GetJwtKey() string {
	return c.jwtKey
}

func (c configManager) GetFromMail() string {
	return c.fromMail
}

func (c configManager) GetMailServer() string {
	return c.mailServer
}

func (c configManager) GetMerchantID() int {
	return c.merchantID
}

func (c configManager) GetMerchantKey() string {
	return c.merchantKey
}

func (c configManager) GetMerchantSalt() string {
	return c.merchantSalt
}

func (c configManager) GetOkUrl() string {
	return c.okUrl
}

func (c configManager) GetFailUrl() string {
	return c.failUrl
}

