package services

import (
	"github.com/go-resty/resty/v2"
	"github.com/pragmataW/apartment_management/models"
)

type IRepo interface {
	CreateFlat(flatNo int) error
	UpdateFlatOwner(apartment models.Apartment) error
	DeleteFlat(flatNo int) error
	GetAllInfoAboutFlat(flatNo int) (models.Apartment, error)
	GetAllInfoAboutAllFlats() ([]models.Apartment, error)
	GetDuesCount(flatNo int) (int, error)
	AddDues(flatNo int) error
	AddDuesForAll() error
	DeleteDues(flatNo int) error
	DeleteDuesByEmail(email string) error
	GetPasswordAndFlatNoByEmail(email string) (string, int, error)
	GetAllAnnouncements() ([]models.Announcement, error)
	AddAnnouncement(announcement models.Announcement) error
	AddMerchant(uuid string, email string) error
	GetEmailFromMerchant(merchantOID string) (string, error)
}

type IEncrypt interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(encryptedText string) (string, error)
}

type IConfigManager interface {
	GetAdminPassword() string
	GetJwtKey() string
	GetFromMail() string
	GetMailServer() string
}

type service struct {
	Repo          IRepo
	Encryptor     IEncrypt
	ConfigManager IConfigManager
	RestyClient   *resty.Client
}

type serviceOption func(*service)

func NewService(options ...serviceOption) *service {
	serviceObject := &service{}
	for _, option := range options {
		option(serviceObject)
	}
	return serviceObject
}

func WithRepo(repo IRepo) serviceOption {
	return func(s *service) {
		s.Repo = repo
	}
}

func WithEncryptor(encryptor IEncrypt) serviceOption {
	return func(s *service) {
		s.Encryptor = encryptor
	}
}

func WithConfigManager(configManager IConfigManager) serviceOption {
	return func(s *service) {
		s.ConfigManager = configManager
	}
}

func WithRestyClient(client *resty.Client) serviceOption {
	return func(s *service) {
		s.RestyClient = client
	}
}
