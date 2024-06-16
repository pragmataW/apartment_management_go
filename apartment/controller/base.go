package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/pragmataW/apartment_management/dto"
	"github.com/pragmataW/apartment_management/services"
)

var (
	validate = validator.New()
)

type IService interface {
	LoginAdmin(password string) (string, error)
	LoginUser(flatNo int, mail string, password string) (string, error)
	CreateFlat(flatNo int) error
	UpdateFlatOwner(apartment services.Apartment) error
	DeleteFlat(flatNo int) error
	GetAllInfoAboutFlat(flatNo int) (services.Apartment, error)
	GetAllInfoAboutAllFlat() ([]services.Apartment, error)
	AddDues(flatNo int) error
	DeleteDues(flatNo int) error
	ChangeDuesPrice(price float64)
	ChangePayDay(payDay int) error
	AddAnnouncement(announcement services.Announcement) error
	GetAllAnnouncements() ([]services.Announcement, error)
	SendMail(subject string, body string, mail string) error
	IncreaseDuesAutomatically() error
	GetPaymentToken(payment dto.PaymentSendReq) (string, error)
	PaymentCallback(merchantOID string) error
}

type IConfigManager interface {
	GetMerchantID() int
	GetMerchantKey() string
	GetMerchantSalt() string
	GetFailUrl() string
	GetOkUrl() string
}

type controller struct {
	Service       IService
	ConfigManager IConfigManager
}

type controllerOption func(*controller)

func NewController(options ...controllerOption) *controller {
	controller := &controller{}
	for _, option := range options {
		option(controller)
	}
	return controller
}

func WithService(service IService) controllerOption {
	return func(c *controller) {
		c.Service = service
	}
}

func WithConfigManager(configManager IConfigManager) controllerOption {
	return func(c *controller) {
		c.ConfigManager = configManager
	}
}
