package services

import (
	"errors"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/pragmataW/apartment_management/dto"
	"github.com/pragmataW/apartment_management/models"
	mocks "github.com/pragmataW/apartment_management/service_mocks"
	"github.com/stretchr/testify/assert"
)

func TestLoginAdmin(t *testing.T) {
	configManagerMock := new(mocks.IConfigManager)
	src := NewService(WithConfigManager(configManagerMock))

	configManagerMock.On("GetAdminPassword").Return("123")
	configManagerMock.On("GetJwtKey").Return("123")

	actual, err := src.LoginAdmin("123")
	assert.NoError(t, err)
	assert.NotEqual(t, "", actual)
}

func TestLoginAdminButPasswordDoesNotMatch(t *testing.T) {
	configManagerMock := new(mocks.IConfigManager)
	src := NewService(WithConfigManager(configManagerMock))

	configManagerMock.On("GetAdminPassword").Return("123")
	configManagerMock.On("GetJwtKey").Return("123")

	actual, err := src.LoginAdmin("1234")
	assert.Error(t, err)
	assert.IsType(t, dto.UserDoesNotExists{}, err)
	assert.Equal(t, "", actual)
}

func TestLoginUser(t *testing.T) {
	mockRepo := new(mocks.IRepo)
	mockEncryptor := new(mocks.IEncrypt)
	mockConfigManager := new(mocks.IConfigManager)

	service := NewService(WithRepo(mockRepo), WithConfigManager(mockConfigManager), WithEncryptor(mockEncryptor))

	email := "test@example.com"
	flatNo := 101
	password := "securepassword"
	encryptedPassword := "encryptedpassword"

	// Mock repo behavior
	mockRepo.On("GetPasswordAndFlatNoByEmail", email).Return(encryptedPassword, flatNo, nil)
	// Mock encryptor behavior
	mockEncryptor.On("Decrypt", encryptedPassword).Return(password, nil)
	// Mock config manager behavior
	mockConfigManager.On("GetJwtKey").Return("secretkey")

	// Call the service method
	token, err := service.LoginUser(flatNo, email, password)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify that the expected methods were called
	mockRepo.AssertCalled(t, "GetPasswordAndFlatNoByEmail", email)
	mockEncryptor.AssertCalled(t, "Decrypt", encryptedPassword)
	mockConfigManager.AssertCalled(t, "GetJwtKey")
}

func TestLoginUserInvalidCredentials(t *testing.T) {
	mockRepo := new(mocks.IRepo)
	mockEncryptor := new(mocks.IEncrypt)
	mockConfigManager := new(mocks.IConfigManager)

	service := NewService(WithRepo(mockRepo), WithConfigManager(mockConfigManager), WithEncryptor(mockEncryptor))
	email := "test@example.com"
	flatNo := 101
	password := "securepassword"
	encryptedPassword := "encryptedpassword"
	wrongPassword := "wrongpassword"

	// Mock repo behavior
	mockRepo.On("GetPasswordAndFlatNoByEmail", email).Return(encryptedPassword, flatNo, nil)
	// Mock encryptor behavior
	mockEncryptor.On("Decrypt", encryptedPassword).Return(wrongPassword, nil)

	// Call the service method
	token, err := service.LoginUser(flatNo, email, password)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "", token)
	assert.IsType(t, dto.UserDoesNotExists{}, err)

	// Verify that the expected methods were called
	mockRepo.AssertCalled(t, "GetPasswordAndFlatNoByEmail", email)
	mockEncryptor.AssertCalled(t, "Decrypt", encryptedPassword)
}

func TestLoginUserButUserDoesNotExists(t *testing.T) {
	configManagerMock := new(mocks.IConfigManager)
	EncryptorMock := new(mocks.IEncrypt)
	repoMock := new(mocks.IRepo)

	src := NewService(
		WithConfigManager(configManagerMock),
		WithEncryptor(EncryptorMock),
		WithRepo(repoMock),
	)

	configManagerMock.On("GetJwtKey").Return("123")
	EncryptorMock.On("Encrypt", "123").Return("123", nil)
	repoMock.On("CheckCredentials", "deneme@mail.com", "123", 1).Return(false, nil)

	actual, err := src.LoginUser(1, "deneme@mail.com", "123")

	assert.Error(t, err)
	assert.IsType(t, dto.UserDoesNotExists{}, err)
	assert.Equal(t, "", actual)
}

func TestCreateFlat(t *testing.T) {
	repoMock := new(mocks.IRepo)
	src := NewService(WithRepo(repoMock))

	repoMock.On("CreateFlat", 1).Return(nil)

	err := src.CreateFlat(1)
	assert.NoError(t, err)
}

func TestCreateFlatButError(t *testing.T) {
	repoMock := new(mocks.IRepo)
	src := NewService(WithRepo(repoMock))

	repoMock.On("CreateFlat", 1).Return(errors.New("random error"))

	err := src.CreateFlat(1)
	assert.Error(t, err)
}

func TestUpdateFlatOwner(t *testing.T) {
	encryptorMock := new(mocks.IEncrypt)
	repoMock := new(mocks.IRepo)

	src := NewService(
		WithEncryptor(encryptorMock),
		WithRepo(repoMock),
	)

	apartment := Apartment{
		FlatNo:       1,
		OwnerName:    "Yusuf",
		OwnerSurname: "Çiftçi",
		Mail:         "deneme@mail.com",
		Password:     "123",
		DuesCount:    0,
	}

	encryptorMock.On("Encrypt", apartment.Password).Return("123", nil)
	repoMock.On("UpdateFlatOwner", apartment.ToApartmentModel()).Return(nil)

	err := src.UpdateFlatOwner(apartment)
	assert.NoError(t, err)
}

func TestDeleteFlat(t *testing.T) {
	repoMock := new(mocks.IRepo)
	src := NewService(WithRepo(repoMock))

	repoMock.On("DeleteFlat", 1).Return(nil)
	err := src.DeleteFlat(1)
	assert.NoError(t, err)
}

func TestDeleteFlatButWithError(t *testing.T) {
	repoMock := new(mocks.IRepo)
	src := NewService(WithRepo(repoMock))

	repoMock.On("DeleteFlat", 1).Return(errors.New("random error"))
	err := src.DeleteFlat(1)
	assert.Error(t, err)
}

func TestGetAllInfoAboutFlat(t *testing.T) {
	repoMock := new(mocks.IRepo)
	EncryptorMock := new(mocks.IEncrypt)

	service := NewService(
		WithRepo(repoMock),
		WithEncryptor(EncryptorMock),
	)

	flatNo := 1
	repoReturn := models.Apartment{
		FlatNo:       flatNo,
		OwnerName:    "deneme",
		OwnerSurname: "deneme",
		Mail:         "deneme@mail.com",
		Password:     "123",
		DuesCount:    1,
	}

	repoMock.On("GetAllInfoAboutFlat", 1).Return(repoReturn, nil)
	EncryptorMock.On("Decrypt", repoReturn.Password).Return("123", nil)

	actual, err := service.GetAllInfoAboutFlat(flatNo)
	assert.NoError(t, err)
	assert.Equal(t, repoReturn.FlatNo, actual.FlatNo)
	assert.Equal(t, repoReturn.OwnerName, actual.OwnerName)
	assert.Equal(t, repoReturn.OwnerSurname, actual.OwnerSurname)
	assert.Equal(t, repoReturn.Mail, actual.Mail)
	assert.Equal(t, repoReturn.Password, actual.Password)
	assert.Equal(t, repoReturn.DuesCount, actual.DuesCount)
}

func TestGetAllInfoAboutAllFlat(t *testing.T) {
	repoMock := new(mocks.IRepo)
	service := NewService(
		WithRepo(repoMock),
	)

	repoReturn := []models.Apartment{
		{
			FlatNo:       1,
			OwnerName:    "deneme1",
			OwnerSurname: "deneme1",
			Mail:         "deneme1@mail.com",
			Password:     "123",
			DuesCount:    1,
		},
		{
			FlatNo:       2,
			OwnerName:    "deneme2",
			OwnerSurname: "deneme2",
			Mail:         "deneme2@mail.com",
			Password:     "456",
			DuesCount:    2,
		},
	}

	repoMock.On("GetAllInfoAboutAllFlats").Return(repoReturn, nil)

	actual, err := service.GetAllInfoAboutAllFlat()
	assert.NoError(t, err)

	assert.Equal(t, len(repoReturn), len(actual))
	for i, repoApartment := range repoReturn {
		assert.Equal(t, repoApartment.FlatNo, actual[i].FlatNo)
		assert.Equal(t, repoApartment.OwnerName, actual[i].OwnerName)
		assert.Equal(t, repoApartment.OwnerSurname, actual[i].OwnerSurname)
		assert.Equal(t, repoApartment.Mail, actual[i].Mail)
		assert.Equal(t, repoApartment.Password, actual[i].Password)
		assert.Equal(t, repoApartment.DuesCount, actual[i].DuesCount)
	}

	repoMock.AssertExpectations(t)
}

func TestAddDues(t *testing.T) {
	repoMock := new(mocks.IRepo)
	service := NewService(
		WithRepo(repoMock),
	)

	flatNo := 1

	repoMock.On("AddDues", flatNo).Return(nil)

	err := service.AddDues(flatNo)
	assert.NoError(t, err)
	repoMock.AssertExpectations(t)
}

func TestAddDues_Error(t *testing.T) {
	repoMock := new(mocks.IRepo)
	service := NewService(
		WithRepo(repoMock),
	)

	flatNo := 1
	expectedError := errors.New("some error")
	repoMock.On("AddDues", flatNo).Return(expectedError)

	err := service.AddDues(flatNo)
	assert.Equal(t, expectedError, err)
	repoMock.AssertExpectations(t)
}

func TestDeleteDues(t *testing.T) {
	repoMock := new(mocks.IRepo)
	service := NewService(
		WithRepo(repoMock),
	)

	flatNo := 1

	repoMock.On("DeleteDues", flatNo).Return(nil)

	err := service.DeleteDues(flatNo)
	assert.NoError(t, err)

	repoMock.AssertExpectations(t)
}

func TestDeleteDues_Error(t *testing.T) {
	repoMock := new(mocks.IRepo)
	service := NewService(
		WithRepo(repoMock),
	)

	flatNo := 1
	expectedError := errors.New("some error")

	repoMock.On("DeleteDues", flatNo).Return(expectedError)

	err := service.DeleteDues(flatNo)
	assert.Equal(t, expectedError, err)

	repoMock.AssertExpectations(t)
}

func TestChangeDuesPrice(t *testing.T) {
	service := NewService()

	newPrice := 100.0

	service.ChangeDuesPrice(newPrice)

	dto.Mutx.Lock()
	assert.Equal(t, newPrice, dto.DuesPrice)
	dto.Mutx.Unlock()
}

func TestChangePayDay(t *testing.T) {
	service := NewService()

	validPayDay := 15

	err := service.ChangePayDay(validPayDay)
	assert.NoError(t, err)

	dto.Mutx.Lock()
	assert.Equal(t, validPayDay, dto.PayDay)
	dto.Mutx.Unlock()
}

func TestChangePayDay_InvalidLow(t *testing.T) {
	service := NewService()

	invalidPayDay := 0

	err := service.ChangePayDay(invalidPayDay)
	assert.Error(t, err)
	assert.IsType(t, dto.PayDayRangeError{}, err)
}

func TestChangePayDay_InvalidHigh(t *testing.T) {
	service := NewService()

	invalidPayDay := 29

	err := service.ChangePayDay(invalidPayDay)
	assert.Error(t, err)
	assert.IsType(t, dto.PayDayRangeError{}, err)
}

func TestAddAnnouncement(t *testing.T) {
	repoMock := new(mocks.IRepo)
	service := NewService(WithRepo(repoMock))

	announcement := Announcement{AnnouncementID: 1, Title: "Test Title", Content: "Test Content"}
	announcementModel := announcement.ToAnnouncementModel()

	repoMock.On("AddAnnouncement", announcementModel).Return(nil)

	err := service.AddAnnouncement(announcement)
	assert.NoError(t, err)
	repoMock.AssertExpectations(t)
}

func TestAddAnnouncement_Error(t *testing.T) {
	repoMock := new(mocks.IRepo)
	service := NewService(WithRepo(repoMock))

	announcement := Announcement{AnnouncementID: 1, Title: "Test Title", Content: "Test Content"}
	announcementModel := announcement.ToAnnouncementModel()

	expectedError := errors.New("some error")
	repoMock.On("AddAnnouncement", announcementModel).Return(expectedError)

	err := service.AddAnnouncement(announcement)
	assert.Equal(t, expectedError, err)
	repoMock.AssertExpectations(t)
}

func TestGetAllAnnouncements(t *testing.T) {
	repoMock := new(mocks.IRepo)
	service := NewService(WithRepo(repoMock))

	modelAnnouncements := []models.Announcement{
		{AnnouncementID: 1, Title: "Test Title 1", Content: "Test Content 1"},
		{AnnouncementID: 2, Title: "Test Title 2", Content: "Test Content 2"},
	}
	repoMock.On("GetAllAnnouncements").Return(modelAnnouncements, nil)

	actual, err := service.GetAllAnnouncements()
	assert.NoError(t, err)
	assert.Len(t, actual, 2)

	for i, announcement := range actual {
		assert.Equal(t, modelAnnouncements[i].AnnouncementID, announcement.AnnouncementID)
		assert.Equal(t, modelAnnouncements[i].Title, announcement.Title)
		assert.Equal(t, modelAnnouncements[i].Content, announcement.Content)
	}
	repoMock.AssertExpectations(t)
}

func TestGetAllAnnouncements_Error(t *testing.T) {
	repoMock := new(mocks.IRepo)
	service := NewService(WithRepo(repoMock))

	expectedError := errors.New("some error")
	repoMock.On("GetAllAnnouncements").Return(nil, expectedError)

	actual, err := service.GetAllAnnouncements()
	assert.Equal(t, expectedError, err)
	assert.Empty(t, actual)
	repoMock.AssertExpectations(t)
}

func TestSendMail_Success(t *testing.T) {
	configManagerMock := new(mocks.IConfigManager)
	configManagerMock.On("GetFromMail").Return("ciftciyusuf700@gmail.com")
	configManagerMock.On("GetMailServer").Return("http://localhost:2001/sendMail")

	client := resty.New()
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	src := NewService(WithConfigManager(configManagerMock), WithRestyClient(client))

	httpmock.RegisterResponder("POST", "http://localhost:2001/sendMail",
		httpmock.NewStringResponder(200, `OK`))

	err := src.SendMail("Test Subject", "<h1>Test Body</h1>", "ciftciyusuf700@gmail.com")
	assert.NoError(t, err)
}

func TestSendMail_Error(t *testing.T) {
	configManagerMock := new(mocks.IConfigManager)
	configManagerMock.On("GetFromMail").Return("ciftciyusuf700@gmail.com")
	configManagerMock.On("GetMailServer").Return("http://localhost:2001/sendMail")

	restyClient := resty.New()
	src := NewService(WithConfigManager(configManagerMock), WithRestyClient(restyClient))

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:2001/sendMail",
		httpmock.NewStringResponder(500, `Internal Server Error`))

	restyClient.SetTransport(httpmock.DefaultTransport)

	err := src.SendMail("Test Subject", "<h1>Test Body</h1>", "test@mail.com")

	assert.Error(t, err)
	assert.IsType(t, dto.SendMailError{Message: "send mail error"}, err)
}
