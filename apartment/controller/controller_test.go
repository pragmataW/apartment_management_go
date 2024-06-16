package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	mocks "github.com/pragmataW/apartment_management/controller_mocks"
	"github.com/pragmataW/apartment_management/dto"
	"github.com/pragmataW/apartment_management/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginAdmin(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	token := "token"
	mockService.On("LoginAdmin", "adminPassword").Return(token, nil)

	app := fiber.New()
	app.Post("/login/admin", controller.LoginAdmin)

	req := httptest.NewRequest("POST", "/login/admin", strings.NewReader(`{"password":"adminPassword"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	expectedBody := `{"message":"status ok"}`
	assert.Equal(t, expectedBody, string(body))

	cookies := resp.Cookies()
	var found bool
	for _, cookie := range cookies {
		if cookie.Name == "Authentication" {
			found = true
			assert.Equal(t, token, cookie.Value)
		}
	}

	if !found {
		t.Error("Authentication cookie not found")
	}
}

func TestLoginAdminButWrongPassword(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	token := ""
	password := "123"
	mockService.On("LoginAdmin", password).Return(token, dto.PasswordMatchError{
		Message: "password does not match",
	})

	app := fiber.New()
	app.Post("/login/admin", controller.LoginAdmin)

	req := httptest.NewRequest("POST", "/login/admin", strings.NewReader(
		`
		{
			"password":"123"
		}
		`,
	))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, string(respBody), `{"message":"password does not match"}`)
}

func TestLoginUser(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	token := "token"
	flatNo := 1
	mail := "deneme@gmail.com"
	password := "123"
	mockService.On("LoginUser", flatNo, mail, password).Return(token, nil)

	app := fiber.New()
	app.Post("/login/user", controller.LoginUser)

	req := httptest.NewRequest("POST", "/login/user", strings.NewReader(`
		{
			"flat_no":1,
			"password": "123",
			"mail": "deneme@gmail.com"
		}
	`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	expectedBody := `{"message":"status ok"}`
	assert.Equal(t, expectedBody, string(body))

	cookies := resp.Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "Authentication" {
			found = true
			assert.Equal(t, token, cookie.Value)
			break
		}
	}
	assert.Equal(t, true, found)
}

func TestLoginUserButWithWrongPassword(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	token := ""
	flatNo := 1
	mail := "deneme@gmail.com"
	password := "123"
	mockService.On("LoginUser", flatNo, mail, password).Return(token, dto.PasswordMatchError{Message: "password does not match"})

	app := fiber.New()
	app.Post("/login/user", controller.LoginUser)

	req := httptest.NewRequest("POST", "/login/user", strings.NewReader(`
		{
			"flat_no":1,
			"password": "123",
			"mail": "deneme@gmail.com"
		}
	`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	expectedBody := `{"message":"invalid credentials"}`
	assert.Equal(t, expectedBody, string(body))

	cookies := resp.Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "Authentication" {
			found = true
			assert.Equal(t, token, cookie.Value)
			break
		}
	}
	assert.Equal(t, false, found)
}

func TestLogout(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:  "Authentication",
			Value: "token",
			Path:  "/",
		})
		return c.Next()
	})

	app.Post("/logout", controller.Logout)

	logoutReq := httptest.NewRequest("POST", "/logout", nil)
	logoutReq.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(logoutReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	cookies := resp.Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "Authentication" {
			if time.Until(cookie.Expires) > 0 {
				found = true
			}
		}
	}
	assert.False(t, found)
}

func TestCreateFlat(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	flatNo := 1
	mockService.On("CreateFlat", flatNo).Return(nil)

	app := fiber.New()
	app.Post("/createFlat/:flatNo", controller.CreateFlat)

	req := httptest.NewRequest("POST", "/createFlat/1", strings.NewReader(""))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"message":"status ok"}`, string(respBody))

	mockService.AssertExpectations(t)
}

func TestCreateFlatButFlatAlreadyExists(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	flatNo := 1
	mockService.On("CreateFlat", flatNo).Return(dto.FlatAlreadyExists{Message: "flat already exists"})

	app := fiber.New()
	app.Post("/createFlat/:flatNo", controller.CreateFlat)

	req := httptest.NewRequest("POST", "/createFlat/1", strings.NewReader(""))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"message":"flat already exists"}`, string(respBody))

	mockService.AssertExpectations(t)
}

func TestUpdateFlatOwner(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	// Mocked data
	apartment := services.Apartment{
		FlatNo:       1,
		OwnerName:    "John",
		OwnerSurname: "Doe",
		Mail:         "john.doe@example.com",
		Password:     "securepassword",
		DuesCount:    5,
	}
	mockService.On("UpdateFlatOwner", apartment).Return(nil)

	app := fiber.New()
	app.Put("/updateFlatOwner", controller.UpdateFlatOwner)

	reqBody := `{
		"flat_no": 1,
		"owner_name": "John",
		"owner_surname": "Doe",
		"mail": "john.doe@example.com",
		"password": "securepassword",
		"dues_count": 5
	}`
	req := httptest.NewRequest("PUT", "/updateFlatOwner", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	expectedBody := `{"message":"status ok"}`
	assert.Equal(t, expectedBody, string(body))

	mockService.AssertExpectations(t)
}

func TestUpdateFlatOwnerBadRequest(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	app := fiber.New()
	app.Put("/updateFlatOwner", controller.UpdateFlatOwner)

	reqBody := `{`
	req := httptest.NewRequest("PUT", "/updateFlatOwner", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "bad request")
}

func TestUpdateFlatOwnerValidationFailed(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	app := fiber.New()
	app.Put("/updateFlatOwner", controller.UpdateFlatOwner)

	reqBody := `{
		"flat_no": 1,
		"owner_name": "",
		"owner_surname": "Doe",
		"mail": "john.doe@example.com",
		"password": "securepassword",
		"dues_count": 5
	}`
	req := httptest.NewRequest("PUT", "/updateFlatOwner", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	expectedMessage := "Key: 'ApartmentRequest.OwnerName' Error:Field validation for 'OwnerName' failed on the 'required' tag"
	assert.Contains(t, string(body), expectedMessage)
}

func TestUpdateFlatOwnerServiceError(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	apartment := services.Apartment{
		FlatNo:       1,
		OwnerName:    "John",
		OwnerSurname: "Doe",
		Mail:         "john.doe@example.com",
		Password:     "securepassword",
		DuesCount:    5,
	}
	mockService.On("UpdateFlatOwner", apartment).Return(errors.New("service error"))

	app := fiber.New()
	app.Put("/updateFlatOwner", controller.UpdateFlatOwner)

	reqBody := `{
		"flat_no": 1,
		"owner_name": "John",
		"owner_surname": "Doe",
		"mail": "john.doe@example.com",
		"password": "securepassword",
		"dues_count": 5
	}`
	req := httptest.NewRequest("PUT", "/updateFlatOwner", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "service error")

	mockService.AssertExpectations(t)
}

func TestDeleteFlat(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	flatNo := 1
	mockService.On("DeleteFlat", flatNo).Return(nil)

	app := fiber.New()
	app.Delete("/deleteFlat/:flatNo", controller.DeleteFlat)

	req := httptest.NewRequest("DELETE", "/deleteFlat/1", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	expectedBody := `{"message":"status ok"}`
	assert.Equal(t, expectedBody, string(body))

	mockService.AssertExpectations(t)
}

func TestDeleteFlatInvalidParam(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	app := fiber.New()
	app.Delete("/deleteFlat/:flatNo", controller.DeleteFlat)

	req := httptest.NewRequest("DELETE", "/deleteFlat/invalid", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	expectedMessage := "invalid query parameter: flatNo"
	assert.Contains(t, string(body), expectedMessage)
}

func TestDeleteFlatNotFound(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	flatNo := 1
	mockService.On("DeleteFlat", flatNo).Return(dto.ThereIsNoFlat{Message: "there is no flat with the given number"})

	app := fiber.New()
	app.Delete("/deleteFlat/:flatNo", controller.DeleteFlat)

	req := httptest.NewRequest("DELETE", "/deleteFlat/1", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	expectedMessage := "there is no flat with the given number"
	assert.Contains(t, string(body), expectedMessage)

	mockService.AssertExpectations(t)
}

func TestDeleteFlatInternalError(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	flatNo := 1
	mockService.On("DeleteFlat", flatNo).Return(errors.New("internal server error"))

	app := fiber.New()
	app.Delete("/deleteFlat/:flatNo", controller.DeleteFlat)

	req := httptest.NewRequest("DELETE", "/deleteFlat/1", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	expectedMessage := "internal server error"
	assert.Contains(t, string(body), expectedMessage)

	mockService.AssertExpectations(t)
}

func TestDeleteFlatNegativeParam(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	app := fiber.New()
	app.Delete("/deleteFlat/:flatNo", controller.DeleteFlat)

	req := httptest.NewRequest("DELETE", "/deleteFlat/-1", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	expectedMessage := "invalid query parameter: flatNo"
	assert.Contains(t, string(body), expectedMessage)
}

func TestGetAllInfoAboutFlatSuccess(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	expectedApartment := services.Apartment{
		FlatNo:       1,
		OwnerName:    "John",
		OwnerSurname: "Doe",
		Mail:         "john.doe@example.com",
		Password:     "password123",
		DuesCount:    3,
	}

	mockService.On("GetAllInfoAboutFlat", 1).Return(expectedApartment, nil)

	app := fiber.New()
	app.Get("/flat/:flatNo", controller.GetAllInfoAboutFlat)

	req := httptest.NewRequest("GET", "/flat/1", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var respBody dto.ApartmentResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	assert.NoError(t, err)

	expectedResponse := dto.ApartmentResponse{
		FlatNo:       expectedApartment.FlatNo,
		OwnerName:    expectedApartment.OwnerName,
		OwnerSurname: expectedApartment.OwnerSurname,
		Mail:         expectedApartment.Mail,
		DuesCount:    expectedApartment.DuesCount,
		Password:     expectedApartment.Password,
	}
	assert.Equal(t, expectedResponse, respBody)

	mockService.AssertExpectations(t)
}

func TestGetAllInfoAboutFlatInvalidParam(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	mockService.On("GetAllInfoAboutFlat", 0).Return(services.Apartment{}, errors.New("invalid flat number"))

	app := fiber.New()
	app.Get("/flat/:flatNo", controller.GetAllInfoAboutFlat)

	req := httptest.NewRequest("GET", "/flat/0", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var respBody fiber.Map
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	assert.NoError(t, err)

	expectedMessage := "invalid flat number"
	assert.Equal(t, expectedMessage, respBody["message"])

	mockService.AssertExpectations(t)
}

func TestGetAllInfoAboutAllFlatSuccess(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	expectedApartments := []services.Apartment{
		{
			FlatNo:       1,
			OwnerName:    "John",
			OwnerSurname: "Doe",
			Mail:         "john.doe@example.com",
			Password:     "password123",
			DuesCount:    3,
		},
		{
			FlatNo:       2,
			OwnerName:    "Jane",
			OwnerSurname: "Smith",
			Mail:         "jane.smith@example.com",
			Password:     "abc123",
			DuesCount:    2,
		},
	}

	mockService.On("GetAllInfoAboutAllFlat").Return(expectedApartments, nil)

	app := fiber.New()
	app.Get("/flats", controller.GetAllInfoAboutAllFlat)

	req := httptest.NewRequest("GET", "/flats", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var respBody []dto.ApartmentResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	assert.NoError(t, err)

	expectedResponses := []dto.ApartmentResponse{
		{
			FlatNo:       expectedApartments[0].FlatNo,
			OwnerName:    expectedApartments[0].OwnerName,
			OwnerSurname: expectedApartments[0].OwnerSurname,
			Mail:         expectedApartments[0].Mail,
			DuesCount:    expectedApartments[0].DuesCount,
			Password:     expectedApartments[0].Password,
		},
		{
			FlatNo:       expectedApartments[1].FlatNo,
			OwnerName:    expectedApartments[1].OwnerName,
			OwnerSurname: expectedApartments[1].OwnerSurname,
			Mail:         expectedApartments[1].Mail,
			DuesCount:    expectedApartments[1].DuesCount,
			Password:     expectedApartments[1].Password,
		},
	}
	assert.ElementsMatch(t, expectedResponses, respBody)

	mockService.AssertExpectations(t)
}

func TestAddDuesSuccess(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	flatNo := 1

	mockService.On("AddDues", flatNo).Return(nil)

	app := fiber.New()
	app.Post("/addDues/:flatNo", controller.AddDues)

	req := httptest.NewRequest("POST", "/addDues/1", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"message":"status ok"}`, string(body))

	mockService.AssertExpectations(t)
}

func TestDeleteDuesSuccess(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	flatNo := 1

	mockService.On("DeleteDues", flatNo).Return(nil)

	app := fiber.New()
	app.Delete("/deleteDues/:flatNo", controller.DeleteDues)

	req := httptest.NewRequest("DELETE", "/deleteDues/1", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"message":"status ok"}`, string(body))

	mockService.AssertExpectations(t)
}

func TestChangeDuesPriceSuccess(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	reqBody := dto.ChangeDuesPriceReq{
		Price: 500,
	}

	mockService.On("ChangeDuesPrice", reqBody.Price).Return(nil)

	app := fiber.New()
	app.Put("/changeDuesPrice", controller.ChangeDuesPrice)

	reqBodyBytes, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest("PUT", "/changeDuesPrice", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"message":"status ok"}`, string(body))

	mockService.AssertExpectations(t)
}

func TestChangePayDaySuccess(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	reqBody := dto.ChangePayDayReq{
		PayDay: 15,
	}

	mockService.On("ChangePayDay", reqBody.PayDay).Return(nil)

	app := fiber.New()
	app.Put("/changePayDay", controller.ChangePayDay)

	reqBodyBytes, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest("PUT", "/changePayDay", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"message":"status ok"}`, string(body))

	mockService.AssertExpectations(t)
}

func TestChangePayDayBadRequest(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	mockService.On("ChangePayDay", mock.Anything).Return(errors.New("negative pay day is not allowed"))

	app := fiber.New()
	app.Put("/changePayDay", controller.ChangePayDay)

	reqBody := dto.ChangePayDayReq{
		PayDay: -1,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest("PUT", "/changePayDay", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "negative pay day is not allowed")

	mockService.AssertExpectations(t)
}

func TestGetDuesPrice(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	app := fiber.New()
	app.Get("/getDuesPrice", controller.GetDuesPrice)

	expectedResp := dto.DuesPriceResponse{
		DuesPrice: dto.DuesPrice,
	}

	req := httptest.NewRequest("GET", "/getDuesPrice", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var actualResp dto.DuesPriceResponse
	err = json.NewDecoder(resp.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Equal(t, expectedResp.DuesPrice, actualResp.DuesPrice)
}

func TestGetPayDay(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	app := fiber.New()
	app.Get("/getPayDay", controller.GetPayDay)

	expectedResp := dto.PayDayResponse{
		PayDay: dto.PayDay,
	}

	req := httptest.NewRequest("GET", "/getPayDay", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var actualResp dto.PayDayResponse
	err = json.NewDecoder(resp.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Equal(t, expectedResp.PayDay, actualResp.PayDay)
}

func TestAddAnnouncement(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	mockRequest := dto.AnnouncementRequest{
		AnnouncementID: 1,
		Title:          "Test Announcement",
		Content:        "This is a test announcement.",
	}

	mockBody, _ := json.Marshal(mockRequest)
	app := fiber.New()
	app.Post("/addAnnouncement", controller.AddAnnouncement)

	req := httptest.NewRequest("POST", "/addAnnouncement", bytes.NewReader(mockBody))
	req.Header.Set("Content-Type", "application/json")

	mockService.On("AddAnnouncement", mock.Anything).Return(nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), `"message":"status ok"`)

	mockService.AssertExpectations(t)
}

func TestGetAllAnnouncements(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	mockAnnouncements := []services.Announcement{
		{AnnouncementID: 1, Title: "Announcement 1", Content: "Content 1"},
		{AnnouncementID: 2, Title: "Announcement 2", Content: "Content 2"},
	}

	mockService.On("GetAllAnnouncements").Return(mockAnnouncements, nil)

	app := fiber.New()
	app.Get("/getAllAnnouncements", controller.GetAllAnnouncements)

	req := httptest.NewRequest("GET", "/getAllAnnouncements", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var respAnnouncements []dto.AnnouncementResponse
	err = json.NewDecoder(resp.Body).Decode(&respAnnouncements)
	assert.NoError(t, err)
	assert.Equal(t, len(mockAnnouncements), len(respAnnouncements))

	mockService.AssertExpectations(t)
}

func TestSendMail(t *testing.T) {
	mockService := new(mocks.IService)
	controller := NewController(WithService(mockService))

	mockRequest := dto.MailRequest{
		Subject: "Test Subject",
		Body:    "This is a test mail body.",
		ToMail:  "test@example.com",
	}

	mockBody, _ := json.Marshal(mockRequest)
	app := fiber.New()
	app.Post("/sendMail", controller.SendMail)

	req := httptest.NewRequest("POST", "/sendMail", bytes.NewReader(mockBody))
	req.Header.Set("Content-Type", "application/json")

	mockService.On("SendMail", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), `"message":"status ok"`)

	mockService.AssertExpectations(t)
}

func TestGetPaymentToken(t *testing.T) {
	// Create mock instances
	mockConfigManager := new(mocks.IConfigManager)
	mockService := new(mocks.IService)

	// Configure mock config manager
	mockConfigManager.On("GetMerchantID").Return(123)                 // Example return value
	mockConfigManager.On("GetMerchantKey").Return("mocked_key")       // Example return value
	mockConfigManager.On("GetMerchantSalt").Return("mocked_salt")     // Example return value
	mockConfigManager.On("GetOkUrl").Return("http://mock.ok/url")     // Example return value
	mockConfigManager.On("GetFailUrl").Return("http://mock.fail/url") // Example return value

	// Create the controller with mock service and config manager
	controller := NewController(
		WithService(mockService),
		WithConfigManager(mockConfigManager),
	)

	// Mock request body
	mockUserBasket := [][]interface{}{
		{"item1", 1},
		{"item2", 2},
	}
	mockRequest := dto.PaymentGetReq{
		UserName:    "Test User",
		UserAddress: "Test Address",
		UserPhone:   "123456789",
		UserBasket:  mockUserBasket,
		DebugOn:     "1",
		TestMode:    "1",
	}
	mockBody, _ := json.Marshal(mockRequest)

	// Mock expectations for service method
	mockService.On("GetPaymentToken", mock.Anything).Return("mocked_token", nil)

	// Create Fiber app instance for testing
	app := fiber.New()
	app.Post("/getPaymentToken", controller.GetPaymentToken)

	// Create HTTP request
	req := httptest.NewRequest("POST", "/getPaymentToken", bytes.NewReader(mockBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert response status code
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Assert response body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), `"token":"mocked_token"`)

	// Verify mock expectations
	mockConfigManager.AssertExpectations(t)
	mockService.AssertExpectations(t)
}
